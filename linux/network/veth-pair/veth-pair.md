
1、创建一对虚拟的 veth-pair 以太网卡
```shell
ip link add veth0 type veth peer name veth1
```
![img.png](img.png)

2、将 veth 的一端放入另一个 netns
```shell
ip link set veth1 netns <netns-name>
```
当前的 netns 中将看不到 veth1:
![img_1.png](img_1.png)
在放入的另一个 netns 中可以看到 veth1:
![img_2.png](img_2.png)

