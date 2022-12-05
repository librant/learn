
#### client-go 源码结构
- discovery
  - 提供 DiscoveryClient 发现客户端
  - api-versions 
  - api-resources
- dynamic
  - 提供 DynamicClient 动态客户端
- informers 
  - 内置资源的 Informer 实现
- kubernetes 
  - 提供 ClientSet 客户端
- listers 
  - 内置资源提供 lister 功能
- plugin 
  - 云服务提供的插件
- rest
  - 提供 RestClient 客户端
- scale 
  - 提供 ScaleClient 客户端
- tools 
  - 提供常用工具
- transport 
  - 提供安全的 TCP 链接，支持 HTTP Stream
- util 
  - 提供常用方法，WorkQueue 工作队列

---
1) RESTClient 客户端  
最基础的客户端  
```go
config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
restClient, err := rest.RESTClientFor(config)
```

2）ClientSet 客户端

3）DynamicClient 客户端

4）DiscoveryClient 客户端

---
Informer 机制：  
保证消息的实时性，可靠性、顺序性

1) Informer 机制架构设计

- Reflector：list/watch
  - watch 指定的 kubernetes 资源变化（Added/Updated/Deleted）,
    将其资源对象放在本地缓存 DeltaFIFO 中
- DeltaFIFO：对象缓存队列
  - 队列的基本操作方法：Add/Update/Delete/List/Pop/Close 
  - Delta 是一个资源对象存储，保存资源操作的类型：Added/Updated/Deleted/Sync
- Indexer：用来存储资源对象并自带索引功能的本地存储
  - Reflector 从 DeltaFIFO 中消费出来的资源对象存储至 Indexer

2) 资源 Informer   
每一个 kubernetes 资源都实现了 informer 机制
- informer()
- lister()

3) Shared Informer 共享机制   
Shared Informer 可以使同一类资源 Informer 共享一个 Reflector,
通过 map 数据结构实现共享的 Informer 机制   


4) Reflector   
用于监控指定资源的 kubernetes 资源的变化（Added/Updated/Deleted）,
并将资源存放在本地缓存 DeltaFIFO 中   

1、获取资源列表数据
- ListAndWatch() 获取该对象下的所有数据 
  - List() 在程序第一次运行时，获取该资源下所有对象数据并存储至 DeltaFIFO 中
  - GetResourceVersion() 获取资源版本号
    - client-go 执行 Watch 操作时可以根据 ResourceVersion 
        来确定当前资源对象是否发生变化
  - ExtractList() 用于将资源数据转换成资源对象列表
    - 将 runtime.Object 对象转换成 []runtime.Object 对象
  - syncWith：用于将资源对象列表中的资源对象和资源版本号存储至 DeltaFIFO 中，
    并会替换已存在的资源版本号
  - setLastSyncResourceVersion() 设置最新的资源版本号

2、监控资源对象
- Watch() 通过 HTTP 协议与 kube-apiserver 建立长链接，接收资源变更事件
  - 分块传输编码（Chunked Transfer Encoding）


5) DeltaFIFO   
生产者是 Reflector 调用 Add() 方法，消费者是 Controller 调用的 Pop() 方法
- FIFO 一个先进先出的队列，拥有队列操作的基本方法
- Delta 一个资源对象的存储，可以保存资源对象的操作类型（Added/Updated/Deleted/Sync）

1、生产者方法
- queueActionLocked() 函数
  - 通过 f.KeyOf() 函数计算出资源对象的 Key
  - 如果操作类型是 Sync, 则标识该数据来源于 Indexer, 如果 Indexer 中的资源对象已经被删除，则直接返回
  - 将 actionType 和资源对象构成 Delta，添加到 items 中，并通过 dedupDeltas() 进行去重操作
  - 更新构造后的 Delta 并通过 cond.Broadcast() 通知所有消费者解除阻塞

2、消费者方法
- Pop() 函数
  - 从 DeltaFIFO 头部取出最早进入队列中的资源对象数据
  - Pop() 中需要传入 process() 函数，用于接收并处理对象的回调方法
  - 如果 process() 回调处理出错，将该对象重新存入队列

3、Resync 机制   
- 将 Indexer 本地存储中的资源对象同步至 DeltaFIFO 中，并将这些对象设置为 Sync 的操作类型


5）Indexer   
- 存储资源对象并自带索引功能的本地存储，Reflector 从 DeltaFIFO 中将消费出来的资源对象存储至 Indexer
- Indexer 的存储结构
  - Indices：存储缓存器
    - map[string]Index
  - Index：存储缓存数据， K/V
    - map[string]sets.String
  - Indexers：存储器索引
    - map[string]IndexFunc：key 为索引器名称，value 为索引器实现函数
  - IndexFunc：索引器函数
  
1、ThreadSafeMap 并发安全存储
- 增删改查操作方法（Add/Update/Delete/List/Get/Replace/Resync）

2、Indexer 索引器
- updateIndices 
- deleteFromIndices 

3、Indexer 核心实现
- index.ByIndex() 通过执行索引器函数得到索引结果

---
WorkQueue 工作队列
- 有序：按照添加的顺序处理元素
- 去重：相同的元素在同一时间不会被重复处理
- 并发性：多生产者和消费者
- 标记机制：标记一个元素是否被处理，允许重新入队
- 通知机制：shutdown 方法通过信号量通知不再接收新的元素
- 延迟：支持延迟队列
- 限速：支持限速队列
- Metric：支持 metric 监控指标

支持三种队列：
- Interface：FIFO 队列接口，先进先出队列，支持去重机制
- DelayingInterface：延迟队列接口
- RateLimitingInterface：限速队列接口

1) FIFO 队列   
支持最基本的队列方法
- Add()：给队列添加元素 item 
- Len()：返回当前队列的长度
- Get()：获取队列头部的一个元素
- Done()：标记队列中该元素已经被处理
- ShutDown()：关闭队列
- ShuttingDown()：查询队列是否正在关闭

FIFO 队列数据结构：
- queue []t：实际存储元素的地方
- dirty set：保证去重，和并发时，元素只被处理一次
- processing set：用于标记机制

2) 延迟队列   
在原有功能上增加了 AddAfter() 方法，延迟一段时间后再将元素插入 FIFO 队列中   
- waitingLoop() 消费元素数据
- waitingForAddCh

3) 限速队列   
在原有功能上增加了 AddRateLimited(), Forget(), NumRequests() 方法
- When()：获取指定元素应该等待的时间
- Forget()：释放指定元素，清空该元素的排队数
- NumRequests()：获取指定元素的排队数

限速算法：
- 令牌桶算法（BucketRateLimiter）
- 排队指数算法
- 计数器算法
- 混合算法

---
EventBroadcaster 事件管理器   
用于展示集群内发生的情况
- Event 资源对象
  - Normal
  - Warning
- 事件管理机制
  - EventRecorder 事件（Event）生产者
  - EventBroadcaster 事件（Event）消费者
  - broadcasterWatcher 观察者（Watcher）管理

1) EventRecorder   
```shell
k8s.io/client-go/tools/record/event.go

type EventRecorder interface {} 
```
- Event：对刚发生的事件进行记录
- Eventf：通过 fmt.Printf() 格式化输出事件的格式
- PastEventf：允许自定义事件发生的时间，以记录已经发生过的消息
- AnnotationEventf：附加注释（Annotation）字段的 Eventf

2) EventBroadcaster   
消费 EventRecorder 记录的事件并将其分发给目前所有链接的 broadcasterWatcher   
- 阻塞分发机制
  - WaitIfChannelFull
- 非阻塞分发机制
  - DropIfChannelFull

3) broadcasterWatcher   
kubernetes 系统组件自定义处理事件的方式
- StartLogging：将事件写入日志中
- StartRecordingToSink：将事件上报至 API Server

StartEventWatcher():   
- 内部运行一个 goroutine，用于不断监控 EventBroadcaster 来发现事件并调用相关函数对事件进行处理


---
### 【好文推荐】

[client-go源码分析](https://qiankunli.github.io/2020/07/20/client_go.html)