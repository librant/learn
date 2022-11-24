#### Flannel

Flannel是一种基于overlay网络的跨主机容器网络解决方案，也就是将TCP数据包封装在另一种网络包里面进行路由转发和通信。

#### 实现原理

1、网络模式
- UDP （基本弃用）
- VxLAN
- Host-GW

- VxLAN 模式
![img_2.png](img_2.png)

不同节点之间的通信：

1、pod中的数据，根据pod的路由信息，发送到网桥 cni0
2、cni0 根据节点路由表，将数据发送到隧道设备flannel.1
3、flannel.1 查看数据包的目的ip，从flanneld获取对端隧道设备的必要信息，封装数据包
4、flannel.1 将数据包发送到对端设备。对端节点的网卡接收到数据包，发现数据包为overlay数据包，解开外层封装，并发送内层封装到flannel.1 设备
5、flannel.1 设备查看数据包，根据路由表匹配，将数据发送给cni0设备
6、cni0匹配路由表，发送数据到网桥

![img_3.png](img_3.png)

- Host-GW 模式
![img_4.png](img_4.png)

host-gw采用纯静态路由的方式，要求所有宿主机都在一个局域网内，跨局域网无法进行路由。


