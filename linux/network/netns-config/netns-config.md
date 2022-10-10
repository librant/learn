### 常见 netns 配置

1、将 netns 设备设置成 UP:
![img.png](img.png)

```shell
ip netns exec <netns-name> ip link set dev lo up
```
经过设置成 UP 之后，就可以 ping 通 127.0.0.1:
![img_2.png](img_2.png)


