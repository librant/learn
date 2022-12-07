
#### Etcd 组件介绍
保存集群所有的网络配置和对象的状态信息   

需要用 etcd 协同存储配置：
- 网络插件（flannel/calico 等）
- kubernetes 集群本身（各种状态和元数据信息）

---
1) ETCD 原理   
raft 一致性算法，是一款分布式的 K/V 

---
参考文档：
- [Etcd 架构与实现解析](https://jolestar.com/etcd-architecture/)