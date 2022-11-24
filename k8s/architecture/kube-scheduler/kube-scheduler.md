
#### kube-scheduler 组件的启动流程

1、内置调度算法的注册

2、Cobra 命令行解析

3、实例化 Scheduler 对象

4、运行 EventBroadcaster 事件管理器

5、运行 HTTP/HTTPS 服务

6、运行 Informer 同步资源

7、领导者选主实例化

8、运行 sched.Run() 调度器

---

- 预选调度算法集 defaultPredicates
- 优选调度算法集 defaultPriorities

```shell
factory.RegisterAlgorithmProvider() --> algorithmProviderMap
```

亲和性和反亲和性：
```go
type Affinity struct {
	NodeAffinity    *NodeAffinity
	PodAffinity     *PodAffinity
	PodAntiAffinity *PodAntiAffinity
} 
```
---

1、kube-scheduler: watch api-server: Added/Updated/Deleted

- podInformer
- nodeInformer
SchedulingCache(Indexer)

--> SchedulingQueue  (low/mid/high)

```shell
scheduleOne():
sched.config.NextPod(): 从优先级队列中获取一个优先级最高的待调度 Pod 资源对象（阻塞模式）
sched.schedule(pod): 调度函数执行预选调度算法和优选调度算法
sched.preempt(): 抢占低优先级 Pod 资源
sched.bind(): 将合适的节点与 Pod 资源对象绑定在一起
```



