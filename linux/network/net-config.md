
1、设置 ip_forward：
- 永久生效:
修改文件： /etc/sysctl.conf
![img_1.png](img_1.png)
- 临时生效：
```shell
echo 1 > /proc/sys/net/ipv4/ip_forward
```

