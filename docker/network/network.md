[一文彻悟容器网络通信](https://mp.weixin.qq.com/s/Hr9qpkfTWP9jxYR2sNOeFA)

### Docker 网络模式
- bridge
- host
- container
- none

bridge 模式：
![img.png](img.png)
![img_3.png](img_3.png)

1、docker 启动后默认创建的 docker0 网桥
```shell
brctl show
```
![img_1.png](img_1.png)

查看路由
```shell
route -n
```
![img_2.png](img_2.png)

2、常用的 Docker 网络技巧
1）查看容器 IP
```shell
docker inspect -f "{{ .NetworkSettings.IPAddress }}" <containerID or Name>
```

