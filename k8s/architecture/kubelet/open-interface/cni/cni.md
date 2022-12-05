#### CNI（Container Network Interface）

```go
type CNI interface {
    AddNetworkList (net *NetworkConfigList, rt *RuntimeConf) (types.Result, error)
    DelNetworkList (net *NetworkConfigList, rt *RuntimeConf) error
    AddNetwork (net *NetworkConfig, rt *RuntimeConf) (types.Result, error)
    DelNetwork (net *NetworkConfig, rt *RuntimeConf) error
}
```

CNI 常用插件：
- loopback
- Bridge
- PTP
- Macvlan
- IPvlan
- third-party

---
- 二层负载均衡：基于 MAC 地址的二层负载均衡。
- 三层负载均衡：基于 IP 地址的负载均衡。
- 四层负载均衡：基于 IP+端口 的负载均衡。
- 七层负载均衡：基于 URL 等应用层信息的负载均衡。

---
kubelet 要使用 CNI 网络驱动需要配置启动参数： 
- --network-plugin=cni
- --cni-conf-dir (默认为：/etc/cni/net.d)
- --cni-bin-dir (默认为：/opt/cni/bin)

---
1、CNI 网络插件的工作流程  
1）kubelet 通过调用 CRI 接口（RunSandbox()） 创建 pause 容器，生成对应的 network namespace  
2）调用网络驱动（driver） --> CNI 的方式 --> 具体的 CNI 插件  
3）CNI 插件给 pause 容器配置正确的网络，Pod 中的其他容器共用 pause 的网络栈

2、CNI 工作原理  
- cni 配置文件目录
- cni 二进制文件

3、cni 解决的问题  
- 容器 IP 地址重复问题
- 容器 IP 地址路由问题

---
CNI 网络插件的开发方式   
- CNI 插件的详细工作流程
1）kubelet 的 grpc-client 调用 CRI grpc-server (dockerd/containerd)，创建一个 pod
2) grpc-server 按照一定的流程去 pull image, 创建 Sandbox(pause), 创建 netns，启动容器，将容器加入 Sandbox()
3) grpc-server 读取主机上 cni 配置（/etc/cni/net.d），获取 cni 的 name
4) 在 （/opt/kubernetes/cni/bin） 下访问 name 的二进制文件，
   grpc-server 传入 containerID, netns，eth-name，pod-name 等参数信息

- CNI 插件开发框架
1) cmdAdd
2) cmdDel

- CNI 设计考量   
1）容器运行时必须在调用任何插件之前为容器创建一个新的网络命名空间   
2）容器运行时必须按顺序为每个网络执行相应的插件，将容器添加到每个网络中   
3）在完成容器生命周期后，运行时必须以相反的顺序执行插件（相对于执行添加容器的顺序）以将容器与网络断开连接   
4）容器运行时不能为同一容器调用并行操作，但可以为不同的容器调用并行操作   
5）容器必须由 ContainerID 唯一标识。存储状态的插件应该使用（网络名称，容器 ID）的主键来完成   
6）运行时不能调用同一个网络名称或容器 ID 执行两次 ADD（没有相应的 DEL）

---
参考项目：  
- https://github.com/containernetworking/cni
- https://github.com/y805939188/simple-k8s-cni