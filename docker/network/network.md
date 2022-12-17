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

---
1) 查看节点上的 docker 网络
```shell
docker network ls
```
![img_4.png](img_4.png)

2) 查看网络信息
```shell
docker network inspect network_id
```

![img_5.png](img_5.png)

3) 运行一个 nginx 容器
```shell
docker run --name nginx-test -p 8080:80 -d nginx
```

4) 查看当前运行的容器
```shell
docker ps
```
![img_6.png](img_6.png)

5) 查看节点上的路由
```shell
route -n
```
![img_7.png](img_7.png)



