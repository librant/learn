
#### kube-proxy
出站流量的负载均衡器，监控 kubernetes API Service 并持续将服务的 IP (ClusterIP) 映射到运行状况良好的 Pod， 
落实到主机上就是 iptables/ipvs 等路由规则；

1、Service
- kubernetes 使用 Labels 将多个相关的 Pod 组合成一个逻辑单元
- Service 具有稳定的 IP 地址和端口，并在一组匹配的后端 Pod 之间提供负载均衡

2、Service 的三个 port
- port: 标识 service 暴露的服务端口，也是客户端访问的端口
- targetPort: 应用程序实际监听 Pod 内流量的端口
- nodePort: 集群外部访问 service 入口的一种方式

kube-proxy 通过节点上的 iptables 规则管理 port 和 targetPort 端口的重映射过程；

3、Service 的类型
- ClusterIP: 方便集群内 pod 到 pod 的调用
- Load Balancer
- NodePort: 乞丐版的 Load Balancer，为 service 在集群的每个节点上分配一个真实的端口（会同时分配 ClusterIP）
- Headless: 无头服务，不分配 ClusterIP

4、Ingress
授权入站连接到达集群内服务的规则集合，支持自定义 Service 的访问策略；
internet --> ingress --> service
- Ingress 可以基于客户端请求的 URL 做流量分发，转发给不同的 service 后端
- Ingress controller: 用户自己实现
  - List/Watch Service/Endpoints/Ingress 对象，并根据信息刷新外部 LB 的规则

外部请求 --> Load Balancer --> Ingress Controller --> list/watch --> Service/Endpoints/Ingress

5、kube-proxy 转发规则  
--proxy-mode: 参数进行配置
- userspace
- iptables
  - NodePort 类型的 service 创建的 iptables 规则：
    - KUBE-NODEPORTS  
    - KUBE-SVC-XXXXXXXXXXXX
    - KUBE-SEP-XXXXXXXXXXXX (后端有几个节点，就有几条链)
    - KUBE-SERVICES
  - SNAT
    - KUBE-MARK-MASQ 0x4000/0x4000
- ipvs
  - LVS 的负载均衡模块，基于 netfilter

---
参考文章：  
[记一次Docker/Kubernetes上无法解释的连接超时原因探寻之旅](https://blog.csdn.net/M2l0ZgSsVc7r69eFdTj/article/details/81380446)

