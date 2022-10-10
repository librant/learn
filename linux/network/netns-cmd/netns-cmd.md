### 常见 netns 操作

1、创建 add
```shell
ip netns add <netns-name>
```
![img_3.png](img_3.png)

创建完成，会在 /var/run/netns 下生成挂载点：
![img_4.png](img_4.png)

2、查看 list
```shell
ip netns list
```
![img_2.png](img_2.png)

3、进入 exec
```shell
ip netns exec <netns-name> ip link list
```
![img_5.png](img_5.png)

从上图可以看出，自带的 lo 设备状态还是 DOWN:
```shell
ip netns exec <netns-name> ping 127.0.0.1
```
![img.png](img_1.png)

4、删除 delete
```shell
ip netns delete <netns-name>
```
![img_6.png](img_6.png)


