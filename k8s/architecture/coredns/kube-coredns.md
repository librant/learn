
1、DNS 服务的功能  
- 负责解析 kubernetes 集群内的 Pod 和 Service 域名，集群内容器使用；

2、kube-coredns 部署
- 部署好后对暴露一个 service，集群内的容器通过访问该服务的 ClusterIP + 53 端口获得域名解析服务
- Service 中的 ClusterIP 一般情况下是固定的
- 容器内的进程把 DNS Server 写入 /etc/resolv.conf 文件， 由 kubelet 刷新

3、域名解析基本原来
- FQDN（完全限定域名）引用，或者通过 service 本身的 name 引用；
- DNS 记录
  - A 记录：用于将域或者子域指向某个 IP 地址的 DNS 记录的基本类型
    - your-svc.your-namespace.svc.cluster.local
  - SRV 记录：通过描述某些服务协议和地址促进服务发现的
  - CNAME 记录：用于将域或者子域指向另一个主机名
    - 用于联合服务的跨集群服务发现

4、DNS 使用
- kubernetes 域名解析策略: 域名解析策略对应 Pod 配置中的 dnsPolicy
  - None：允许 Pod 忽略 kubernetes 环境中的 DNS 配置，应使用 dnsConfigPod 规范中的字段提供所有 DNS 设置
  - ClusterFirstWithHostNet：对于 HostNetWork 运行的 Pod，明确设置该值
  - ClusterFirst：任何与配置集群域不匹配的 DNS 查询将转发到从宿主机上继承的上游域名服务器（默认配置）
  - Default：Pod 从宿主机上继承名称解析配置

5、调试 DNS
```shell
nslookup: cant resolve kubernetes.default
```
- 检查容器中的 resolve.conf 文件
- 检查 kube-coredns 插件是否已经启用（coredns pod 是否运行正常， service 对应的 endpoints 是否存在）

6、网络故障定位指南
- IP 转发  
tcpdump 发送大量重复的 SYN 数据包，但没有收到 ACK
  - 检查 net.ipv4.ip_forward 是否开启

- 桥接  
kubernetes 通过 bridge-netfilter 配置使 iptables 规则应用在 Linux 网桥上
  - 检查 net.bridge.bridge-nf-call-iptables 是否开启

- Pod CIDR 冲突  
当 Pod 子网和主机网络出现冲突时，Pod 和 Pod 之间通信就会因为路由中断
  - 检查网络配置，VLAN 或 VPC 之间不会有重叠
  - 有冲突，修改 kubelet 的 pod-cidr 参数指定的地址范围

- hairpin
Pod 无法通过 Service IP 访问自己，可能就是 hairpin 设置问题
  - kubelet 提供 --hairpin-mode 的标志 /sys/devices/virtual/net/
    - hairpin-veth
    - promiscuous-bridge

- 故障排查工具
  - tcpdump
```shell
tcpdump -i any host 172.x.x.x
```
  - nslookup

---

1、CoreDNS 三个特点：
- 插件化（Plugin）：基于 Caddy 服务框架，CoreDNS 实现一个插件链的结构
- 配置简单化：引入表达了更强的 DSL
- 一体化的解决方案：单独的可执行文件
