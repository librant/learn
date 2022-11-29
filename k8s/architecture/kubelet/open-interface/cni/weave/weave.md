
Weave: 支持数据加密的网络插件

1、Weave 的工作模式  
自己集成高可用的数据存储功能
- UDP 
- fastpath (VXLAN 和 OVS)

2、Weave 实现原理  
- 数据平面：通过封包的实现了 L2 overlay

3、Weave 网络模型
1) 指定容器 IP
```shell
weave run 192.168.0.2/24 -itd docker.io/centos /bin/bash
```
2) 容器互联
```shell
weave connect 103.10.86.239
```
3) 查看容器路由信息
```shell
weave ps
```

4、Weave 其他特性  
- 应用隔离
- 安全性
- 服务发现
- 无需额外的集群存储
- 性能
- 组播支持

5、WeaveScope 网络监控  
对 Docker 运行网络进行监控


