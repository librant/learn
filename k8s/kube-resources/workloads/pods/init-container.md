
#### Init 容器

Pod 能够具有多个容器，应用运行在容器里面，但是它也可能有一个或多个先于应用容器启动的 Init 容器；

- Init 容器总是运行到成功为止
- 每个 Init 容器都必须在下一个 Init 容器启动之前成功完成
- PodSpec: initContainers 字段定义
- Init 容器支持应用容器的全部字段和特性，包括资源限制、数据卷和安全设置
- Init 容器不支持 Readiness Probe

1) 使用 Init 容器
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
  labels:
    app: myapp
spec:
  containers:
    - name: myapp-container
      image: busybox
      command: ['sh', '-c', 'echo The app is running! && sleep 3600']
  initContainers:
    - name: init-myservice
      image: busybox
      command: ['sh', '-c', 'until nslookup myservice; do echo waiting for myservice; sleep 2; done;']
    - name: init-mydb
      image: busybox
      command: ['sh', '-c', 'until nslookup mydb; do echo waiting for mydb; sleep 2; done;']

```

2) 具体行为   
- 在所有的 Init 容器没有成功之前，Pod 将不会变成 Ready 状态
- Init 容器的端口将不会在 Service 中进行聚集
- 正在初始化中的 Pod 处于 Pending 状态，但应该会将 Initializing 状态设置为 true
- 如果 Pod 重启，所有 Init 容器必须重新执行
- 对 Init 容器 spec 的修改被限制在容器 image 字段，修改其他字段都不会生效
- 更改 Init 容器的 image 字段，等价于重启该 Pod

资源：
- 在所有 Init 容器上定义的，任何特殊资源请求或限制的最大值，是 有效初始请求/限制
- Pod 对资源的有效请求/限制要高于：
  - 所有应用容器对某个资源的请求/限制之和
  - 对某个资源的有效初始请求/限制
- 基于有效请求/限制完成调度，这意味着 Init 容器能够为初始化预留资源，这些资源在 Pod 生命周期过程中并没有被使用。
- Pod 的 有效 QoS 层，是 Init 容器和应用容器相同的 QoS 层

Pod 重启原因：
- 用户更新 PodSpec 导致 Init 容器镜像发生改变。应用容器镜像的变更只会重启应用容器。
- Pod 基础设施容器被重启。这不多见，但某些具有 root 权限可访问 Node 的人可能会这样做。
- 当 restartPolicy 设置为 Always，Pod 中所有容器会终止，强制重启，由于垃圾收集导致 Init 容器完整的记录丢失。