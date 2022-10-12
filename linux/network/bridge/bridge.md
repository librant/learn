
### bridge: 
多个 network namespace 中间进行连接（二层网络设备）

1、创建网桥
```shell
ip link add name <bridge-name> type bridge
```
![img.png](img.png)

2、启动网桥
```shell
ip link set <bridge-name> up
```

---
---
1、通过 brctl 工具进行创建网桥：（工具在 bridge-utils 包中）
```shell
brctl addbr br-new
```
![img_1.png](img_1.png)

2、查看当前网桥上网络设备
```shell
bridge link
```
![img_2.png](img_2.png)

```shell
brctl show
```
![img_3.png](img_3.png)

3、删除网桥上的 veth 设备
```shell
brctl delif <br-name> <veth-name>
```
