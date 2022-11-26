
[Macvlan和IPvlan基础知识](https://mp.weixin.qq.com/s/r_CuqjypaaMRDZfW-RHjxw)

1、创建 macvlan 设备
```shell
ip link add <macvlna-name> link eth0 type macvlan mode bridge
```
![img.png](img.png)


2、启动 macvlan 设备
```shell
ifconfig <macvlan-new> up
```
![img_1.png](img_1.png)

---
Macvlan 工作模式：
- bridge
- VEPA
- Private
- Passthru
- Source

Macvlan 可以让一块物理网卡上可以绑定多个 IP 地址，每个 IP 地址都有自己的 MAC 地址。

