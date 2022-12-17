#### Flannel

[Flannel](https://github.com/flannel-io/flannel) 源码

Flannel是一种基于 overlay 网络的跨主机容器网络解决方案，也就是将 TCP 数据包封装在另一种网络包里面进行路由转发和通信。

#### 实现原理

1、网络模式 backend
- UDP （基本弃用）
  - cni0 --> bridge 网桥设备  
  - flannel0 --> tun 虚拟设备  
  - flanneld --> agent (DaemonSet) 网络包的封包和解包   
```shell
// 将三层的 IP 报文封装在一个 UDP 的报文中, 三层报文的两个 IP 不在同一个节点上
节点1 的 PodA --> cni0 --> flannel0 (tunl设备) --> flanneld:8285 --> 
eth0 (节点1 的公网 IP 所在网卡) --> internet --> eth1 (节点2 的公网 IP 所在网卡)
--> flanneld:8285 --> flannel0 (tunl 设备) --> cni0 --> 节点2 的 PodB
```

- VxLAN
  - VTEP 的 MAC 地址是通过 APIServer 处 watch Node 发现 
  - flannel.1 (flannel.[VNI])
  - flanneld：将 VTEP 设备相关信息上报到 etcd 中, 当网络中有新节点加入集群中，并向 etcd 注册时：
    - 在节点中创建一条该节点所属网络的路由表 --》 Pod 中的流量路由到 flannel.1 接口 （route -n）
    - 在节点中添加一条该节点 IP 及 VTEP 设备的静态 ARP 缓存 (arp -n) 
```shell
// 将三层的 IP 报文封装在一个 UDP 的报文中, 三层报文的两个 IP 不在同一个节点上
节点1 的 PodA --> cni0 --> flannel.1 (VTEP 设备)  --> eth0 (节点1 的公网 IP 所在网卡) 
--> internet --> eth1 (节点2 的公网 IP 所在网卡) --> flannel.1 (VTEP 设备) --> cni0 --> 节点2 的 PodB
```  

- Host-GW
  - 通过容器路由表 IP 到 cni0
  - 到达 cni0 的 IP 包匹配到 HostA 中的路由规则 --> Host B
  - IP 包通过物理网络到达 Host B 中的 eth1
  - 到达 eth1 中的 IP 包匹配到 Host B 中的路由表项， IP 包转发到 cni0
  - cni0 将 IP 包转发给链接在 cni0 上的目的容器
```shell
// 增加路由来将报文转发到对应的节点上
节点1 的 PodA --> cni0 --> eth0 (节点1 的公网 IP 所在网卡) 
--> internet --> eth1 (节点2 的公网 IP 所在网卡)  --> cni0 --> 节点2 的 PodB
```

- VxLAN 模式
![img_2.png](img_2.png)

1）创建 VXLAN 设备，不再监听任何 l2 miss/l3 miss 事件消息  
2）为远端的子网创建路由  
3）为远端的主机创建静态 ARP 表项  
4）创建 FDB 表项， 包括 VTEP MAC 和远端 Flannel 对外 IP  

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

---
Flannel 使用 etcd 存放网络配置、已分配 的 subnet、host 的 IP 等信息；  

2、数据转发流程
![img_6.png](img_5.png)

- 1）容器直接使用目标容器的 IP 进行访问，默认通过容器内部的 eth0 发送出去；
- 2）报文通过 veth pair 被发送到 vethXXX；
- 3）vethXXX 是直接连接到虚拟交换机 docker0 的，报文通过虚拟 bridge docker0 发送出去；
- 4）查找路由表，外部容器 ip 的报文都会转发到 flannel0 虚拟网卡，这是一个 P2P 的虚拟网卡，然后报文就被转发到监听在另一端的 flanneld；
- 5）flanneld 通过 etcd 维护了各个节点之间的路由表，把原来的报文 UDP 封装一层，通过配置的 iface 发送出去；
- 6）报文通过主机之间的网络找到目标主机；
- 7）报文继续往上，到传输层，交给监听在 8285 端口的 flanneld 程序处理；
- 8）数据被解包，然后发送给 flannel0 虚拟网卡；
- 9）查找路由表，发现对应容器的报文要交给 docker0；
- 10）docker0 找到连到自己的容器，把报文发送过去；

3、Flannel 启动  
1）先于 Docker 启动  
2）从 etcd 中获取 network（大网）的配置信息
3）划分 subnet（子网），并在 etcd 中进行注册  
4）将子网信息记录到 flannel 维护的 /run/flannel/subnet.env 文件中  
```shell
Flannel_NETWORK=10.0.0.0/16
Flannel_SUBNET=10.0.34.1/24
Flannel_MTU=1472
Flannel_IPMASQ=false
```
5）将 subnet.env 转写成一个 Docker 的环境变量文件 /run/flannel/docker
```shell
DOCKER_OPT_BIP="-bip=10.0.34.1/24"
DOCKER_OPT_IPMASQ="--ip-masq=true"
DOCKER_OPT_MTU="--mtu=1472"
DOCKER_NETWORK_OPTIONS="--bip=10.0.34.1/24 --ip-masq=true --mtu=1472"
```

---

1、跨 node 之间的 pod 进行通信：  
- pod1根据自己的路由表将包发给CNI0
- CNI0根据自己的路由表将包发给flannel.1
- flannel.1根据flanneld提供的IP、MAC信息进行内层封装
- linux内核进行外层封装
- 将包发给node2
- node2接受到包后，内核识别到这是一个vxlan，拆开后将包发给flannel.1设备
- flannel.1根据路由表将包发往CNI0
- CNI0根据ARP协议获得目标MAC地址，将包发给pod2

![img_6.png](img_6.png)

2、flannel 解决的问题：   
- 只要配置 master 节点, 自动配置集群各节点的子网，网关
- 自动创建 cni0 网桥，用于单节点容器间的互连，自动设置 cni0 网卡 ip，并作为节点内容器的网关
- 根据实际使用的类型udp/vxlan/host-gw 在跨节点上自动配置相应的网络路由以及封装规则
  - UDP 模式：创建 flannel0（tun）设备，flanneld 进配置进行 udp 外层（公网ip）分装
  - VxLan 模式：创建 flannel.1(VTEP) 设备，设置相应的fdb转发规则，使用内核 vxlan 模块进行外层（公网ip）封装
  - host-gw 模式：自动配置网段路由进行路由转发
- 自动为外部网络的访问，创建 NAT 规则，用于容器内部访问外部网络
- 通过创建 veth-pair 设备对，一端放在容器内部，另外一端放在 cni0 网桥上，保障容器内部可以直接与 cni0 网桥通信
- flannel 的 udp/vxlan 属于 overlay 的网络技术，安全方面有保障，另外也提供了一个性能较高的 host-gw 方案

---

### Flannel

- 为每个 node 分配 subnet，容器将自动从该子网中获取 IP 地址
- 当有 node 加入到网络中时，为每个 node 增加路由配置

1) Host-gw 方式：
![img_7.png](img_7.png)

---
参考文档：
- [CNI网络插件之flannel](https://blog.csdn.net/m0_48594855/article/details/107145561)
- https://jimmysong.io/kubernetes-handbook/concepts/flannel.html




