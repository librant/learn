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
1）kubelet 通过调用 CRI 接口（CreateSandbox()） 创建 pause 容器，生成对应的 network namespace  
2）调用网络驱动（driver） --> CNI 的方式 --> 具体的 CNI 插件  
3）CNI 插件给 pause 容器配置正确的网络，Pod 中的其他容器共用 pause 的网络栈

2、CNI 工作原理  
- cni 配置文件目录
- cni 二进制文件

3、cni 解决的问题  
- 容器 IP 地址重复问题
- 容器 IP 地址路由问题

---
参考项目：  
- https://github.com/containernetworking/cni