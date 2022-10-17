
### VxLa 基本配置命令

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

