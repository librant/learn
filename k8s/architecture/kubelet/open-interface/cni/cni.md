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

