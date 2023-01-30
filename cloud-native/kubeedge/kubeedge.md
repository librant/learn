### KubeEdge 学习笔记


1、编程模式学习

github.com/kubeedge/kubeedge/cloud/cmd/cloudcore/app/server.go

```go
    // 使用注册的模式
    registerModules(config)
```

```go
// registerModules register all the modules started in cloudcore
func registerModules(c *v1alpha1.CloudCoreConfig) {
	cloudhub.Register(c.Modules.CloudHub)
	edgecontroller.Register(c.Modules.EdgeController)
	devicecontroller.Register(c.Modules.DeviceController)
	nodeupgradejobcontroller.Register(c.Modules.NodeUpgradeJobController)
	synccontroller.Register(c.Modules.SyncController)
	cloudstream.Register(c.Modules.CloudStream, c.CommonConfig)
	router.Register(c.Modules.Router)
	dynamiccontroller.Register(c.Modules.DynamicController)
}
```

每个模块实现自己的注册函数；

--》 cloudhub.Register(c.Modules.CloudHub)
