
#### syncLoop 源码分析

```go
// syncLoop is the main loop for processing changes.
func (kl *Kubelet) syncLoop(updates <-chan kubetypes.PodUpdate, handler SyncHandler) {}
```

```go
plegCh := kl.pleg.Watch(): 从通道中读取事件

--> kl.syncLoopIteration(updates, handler, syncTicker.C, housekeepingTicker.C, plegCh)
```

---
```go
// syncLoopIteration reads from various channels and dispatches pods to the given handler.
func (kl *Kubelet) syncLoopIteration(configCh <-chan kubetypes.PodUpdate, handler SyncHandler,
	syncCh <-chan time.Time, housekeepingCh <-chan time.Time, plegCh <-chan *pleg.PodLifecycleEvent) bool {}
```


```go
// SyncHandler is an interface implemented by Kubelet, for testability
type SyncHandler interface {
    HandlePodAdditions(pods []*v1.Pod)
    HandlePodUpdates(pods []*v1.Pod)
    HandlePodRemoves(pods []*v1.Pod)
    HandlePodReconcile(pods []*v1.Pod)
    HandlePodSyncs(pods []*v1.Pod)
    HandlePodCleanups() error
}

//   - configCh: dispatch the pods for the config change to the appropriate
//     handler callback for the event type

u, open := <-configCh: pod 配置变更事件通道，根据事件类型，调用对应的 handler 的处理函数
```
- kubetypes.ADD:
  handler.HandlePodAdditions(u.Pods)