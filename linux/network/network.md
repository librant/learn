# Network

**【Kubernetes 网络权威指南】**
- 将 P9-P11 中的 c 程序翻译成 go 程序；-- Open

1) netns 命令

```shell
ip netns help

ip netns exec
ip netns list
ip netns delete
```
![img.png](img.png)
```shell
# 添加并启动虚拟网卡tap设备
ip tuntap add dev tap0 mode tap 
ip tuntap add dev tap1 mode tap 
ip link set tap0 up
ip link set tap1 up
# 配置IP
ip addr add 10.0.0.1/24 dev tap0
ip addr add 10.0.0.2/24 dev tap1
# 添加netns
ip netns add ns0
ip netns add ns1
# 将虚拟网卡tap0，tap1分别移动到ns0和ns1中
ip link set tap0 netns ns0
ip link set tap1 netns ns1
```

```shell
# Linux 命令
route -n
iptables -L
```

