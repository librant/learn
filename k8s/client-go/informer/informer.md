
#### client-go 中的 informer 源码分析
client-go informer:
![img.png](img.png)

kubernetes提供了client-go以方便使用go语言进行二次快发

1、informerFactory   
- 用来管理需要多少个对象的 informer 实例

```go
// 创建一个 informer factory
kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)

// factory 已经为所有 k8s 的内置资源对象提供了创建对应 informer 实例的方法，调用具体 informer 实例的 Lister 或 Informer 方法
// 就完成了将 informer 注册到 factory 的过程
deploymentLister := kubeInformerFactory.Apps().V1().Deployments().Lister()

// 启动注册到 factory 的所有 informer
kubeInformerFactory.Start(stopCh)
```

2、SharedInformerFactory 结构
- 统一管理控制器中需要的各资源对象的 informer 实例

```go
type sharedInformerFactory struct {
   client           kubernetes.Interface    // clientset
   namespace        string                  // 关注的 namespace，可以通过 WithNamespace Option 配置
   tweakListOptions internalinterfaces.TweakListOptionsFunc
   lock             sync.Mutex
   defaultResync    time.Duration                               //前面传过来的时间，如30s
   customResync     map[reflect.Type]time.Duration              //自定义resync时间
   informers        map[reflect.Type]cache.SharedIndexInformer  //针对每种类型资源存储一个informer，informer的类型是ShareIndexInformer
   startedInformers map[reflect.Type]bool                       //每个informer是否都启动了
}
```

- client: clientset 客户端
- namespace: factory 关注的命名空间（默认： All Namespace）
- defaultResync: 用于初始化持有的 shareIndexInformer 的 resyncCheckPeriod 和 defaultEventHandlerResyncPeriod 字段，
  用于定时的将 local store（Indexer） 同步到 deltaFIFO
- customResync: 支持针对每一个 informer 来配置 resync 时间,
  通过 WithCustomResyncConfig 这个 Option 配置，否则就用指定的 defaultResync
- informers: factory 管理的 informer 集合
- startedInformers: 记录已经启动的 informer 集合

1) 新建一个 sharedInformerFactory
```go
func NewSharedInformerFactoryWithOptions(client kubernetes.Interface, defaultResync time.Duration, options ...SharedInformerOption) SharedInformerFactory {
   factory := &sharedInformerFactory{
      client:           client,          //clientset，对原生资源来说，这里可以直接使用kube clientset
      namespace:        v1.NamespaceAll, //可以看到默认是监听所有ns下的指定资源
      defaultResync:    defaultResync,   //30s
      //以下初始化map结构
      informers:        make(map[reflect.Type]cache.SharedIndexInformer),
      startedInformers: make(map[reflect.Type]bool),
      customResync:     make(map[reflect.Type]time.Duration),
   }
   return factory
}
```
2) 启动 factory 下的所有 informer
```go
func (f *sharedInformerFactory) Start(stopCh <-chan struct{}) {
   f.lock.Lock()
   defer f.lock.Unlock()

   for informerType, informer := range f.informers {
      if !f.startedInformers[informerType] {
         //直接起gorouting调用informer的Run方法，并且标记对应的informer已经启动
         go informer.Run(stopCh)
         f.startedInformers[informerType] = true
      }
   }
}
```

3) 等待 informer 的 cache 被同步
- sharedInformerFactory 的 WaitForCacheSync() 将会不断调用 factory 持有的所有 informer 的 HasSynced() 方法，直到返回 true
- informer 的 HasSynced() 方法调用的自己持有的 controller 的 HasSynced() 方法
（informer 结构持有 controller 对象，下文会分析 informer 的结构）
- informer 中的 controller 的 HasSynced() 方法则调用的是 controller 持有的 deltaFIFO 对象的 HasSynced() 方法

```go
func (f *sharedInformerFactory) WaitForCacheSync(stopCh <-chan struct{}) map[reflect.Type]bool {}
```

4) factory 为自己添加 informer
```go
func (f *sharedInformerFactory) InformerFor(obj runtime.Object, newFunc internalinterfaces.NewInformerFunc) cache.SharedIndexInformer {}
```
根据对象的类型，返回已经实现的 informer

5) shareIndexInformer 对应的 newFunc 的实现
- client-go 中已经为所有内置对象都提供了 NewInformerFunc()

- pod  
  factory.Core().V1().Pods() 为 factory 添加一个 pod 对应的 shareIndexInformer 的实现
- deployment  
  factory.Apps().V1().Deployments() 为 factory 添加一个 pod 对应的 shareIndexInformer 的实现

3、shareIndexInformer 结构

```go
type sharedIndexInformer struct {
   indexer    Indexer       //informer中的底层缓存cache
   controller Controller    //持有reflector和deltaFIFO对象，reflector对象将会listWatch对象添加到deltaFIFO，同时更新indexer cahce，更新成功则通过sharedProcessor触发用户配置的Eventhandler

   processor             *sharedProcessor //持有一系列的listener，每个listener对应用户的EventHandler
   cacheMutationDetector MutationDetector //可以先忽略，这个对象可以用来监测local cache是否被外部直接修改

   // This block is tracked to handle late initialization of the controller
   listerWatcher ListerWatcher //deployment的listWatch方法
   objectType    runtime.Object

   // resyncCheckPeriod is how often we want the reflector's resync timer to fire so it can call
   // shouldResync to check if any of our listeners need a resync.
   resyncCheckPeriod time.Duration
   // defaultEventHandlerResyncPeriod is the default resync period for any handlers added via
   // AddEventHandler (i.e. they don't specify one and just want to use the shared informer's default
   // value).
   defaultEventHandlerResyncPeriod time.Duration
   // clock allows for testability
   clock clock.Clock

   started, stopped bool
   startedLock      sync.Mutex

   // blockDeltas gives a way to stop all event distribution so that a late event handler
   // can safely join the shared informer.
   blockDeltas sync.Mutex
}
```

- indexer: 底层缓存，其实就是一个 map 记录对象
- controller: informer 内部 一个 controller
- reflector: 根据用户定义的 ListWatch() 方法获取对象并更新增量队列 DeltaFIFO
- processor: 知道如何处理 DeltaFIFO 队列中的对象，实现是 sharedProcessor{}
- listerWatcher: 知道如何 list 对象和 watch 对象的方法
- objectType: pod{}/deployment{}
- resyncCheckPeriod: 给自己的 controller 的 reflector 每隔多少s<尝试>调用 listener 的 shouldResync() 方法
- defaultEventHandlerResyncPeriod: 通过 AddEventHandler() 方法给 informer 配置回调时如果没有配置的默认值，
  这个值用在 processor 的 listener 中判断是否需要进行 resync，最小1s

1) sharedIndexInformer 的 Run() 方法
```go
// k8s.io/client-go/tools/cache/shared_informer.go
func (s *sharedIndexInformer) Run(stopCh <-chan struct{}) {}
```

2) 为 shareIndexInformer 创建 controller
- 通过执行 reflector.Run() 方法启动 reflector，开启对指定对象的 listAndWatch() 过程，获取的对象将添加到 reflector 的 deltaFIFO 中
- 通过不断执行 processLoop() 方法，从 DeltaFIFO pop 出对象，再调用 reflector 的 Process（就是shareIndexInformer的HandleDeltas方法）处理

3) controller 的 processLoop() 方法

```go
func (c *controller) processLoop() {
   for {
      obj, err := c.config.Queue.Pop(PopProcessFunc(c.config.Process))
      if err != nil {
         if err == ErrFIFOClosed {
            return
         }
         if c.config.RetryOnError {
            // This is the safe way to re-enqueue.
            c.config.Queue.AddIfNotPresent(obj)
         }
      }
   }
}
```

4) deltaFIFO pop 出来的对象处理逻辑
![img_1.png](img_1.png)
sharedIndexInformer 的 HandleDeltas 处理从 deltaFIFO pod 出来的增量时, 先尝试更新到本地缓存 cache,
更新成功的话会调用 processor.distribute() 方法向 processor 中的 listener 添加 notification，
listener 启动之后会不断获取 notification 回调用户的 EventHandler 方法

- Sync: reflector list 到对象时 Replace 到 deltaFIFO 时 daltaType 为 Sync 或者 resync 把 local strore 中的对象加回到 deltaFIFO
- Added、Updated: reflector watch 到对象时根据 watch event type 是 Add 还是 Modify 对应 deltaType 为 Added 或者 Updated
- Deleted: reflector watch 到对象的 watch event type 是 Delete 或者 re-list Replace 到 deltaFIFO 时 local store 多出的对象以 Delete 的方式加入 deltaFIFO

```go
// k8s.io/client-go/tools/cache/shared_informer.go
func (s *sharedIndexInformer) HandleDeltas(obj interface{}) error {}
```

---
**参考文档：**   
1) https://jimmysong.io/kubernetes-handbook/develop/client-go-informer-sourcecode-analyse.html