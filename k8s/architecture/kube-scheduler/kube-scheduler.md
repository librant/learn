
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


