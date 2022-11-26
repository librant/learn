#### IPvlan 

从一个主机接口虚拟出多个虚拟网络接口；
区别于 Macvlan 的是：所有的虚拟接口都有相同的 MAC 地址；

1、工作模式
- L2 模式
- L3 模式

---
```shell
ip netns add net1
ip netns add net2

ip link add ipv1 link eth0 type ipvlan mode 13
ip link add ipv2 link eth0 type ipvlan mode 13

# 把 IPvlan 接口放入 network namespace 中
ip link set ipv1 netns net1
ip link set ipv2 netns net2
ip netns exec net1 ip link set ipv1 up
ip setns exec net2 ip link set ipv2 up

# 给两个虚拟网卡接口配置不同网络 IP 地址, 并配置好路由
ip netns exec net1 ip addr add 10.0.1.10/24 dev ipv1
ip netns exec net2 ip addr add 192.168.1.10/24 dev ipv2
ip netns exec net1 ip route add default dev ipv1
ip netns exec net2 ip route add default dev ipv2

# 测试连通性
ip netns exec net1 ping -c 3 192.168.1.10
```