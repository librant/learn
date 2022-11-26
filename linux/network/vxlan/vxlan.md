
### VxLan 基本配置命令

1、创建 VxLan 接口
```shell
ip link add VxLan-new0 type vxlan id 42 group 239.1.1.1 dev eth0 dstport 4789
```
![img.png](img.png)


2、删除 VxLan 接口
```shell
ip link delete <vxlan-name>
```

3、查看 VxLan 接口信息
```shell
ip -d link show <vxlan-name>
```
![img_1.png](img_1.png)

4、查看 VxLan 接口转发表
```shell
bridge fdb show dev <vxlan-name>
```
![img_2.png](img_2.png)

- 配置 IP:
```shell
ip addr add 172.17.1.2/24 dev <vxlan-name>
```
- 启用：
```shell
ip link set <vxlan-name> up
```
![img_3.png](img_3.png)

- 查看路由：
```shell
ip route
```
![img_4.png](img_4.png)

```shell
bridge fdb
```
![img_5.png](img_5.png)

---

### VxLan 协议原理

在 VXLAN 网络的每个端点都有一个 VTEP 的设备，负责 VXLAN 报文的封包和解包；

- VTEP: VXLAN 的边缘设备，用来进行 VXLAN 报文的处理（封包和解包）
- VNI: 是每个 VXLAN 的标，VNI 相同的机器逻辑上处于同一个二层网络中
- tunnel: 隧道是一个逻辑概念，在 VXLAN 模型中没有具体的物理实体相对应

VXLAN 封包格式：
- Ethernet Header (14 bytes)
- IP Header (20 bytes)
- UDP Header (8 bytes)
- VXLAN Header (8 bytes)

---
### 多播模式下 VXLAN 的通信全过程：

1）主机 1 的 vxlan0 发送 ping 报文到主机 2 的 172.17.1.3 地址，内核发现源地址和目的IP地址是在同一个局域网内。
需要知道对方的 AMC 地址，而本地又没有缓存，故先发送一个 ARP 查询报文；
2） 
