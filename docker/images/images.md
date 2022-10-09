
### containerd

```shell
crictl images ls: 查看镜像

// 拉取镜像
crictl image pull <images-name>

// 将镜像导入本地
crictl image export <output-filename> <image-name>

// 将本地镜像导入
crictl -n=k8s.io images import <filename-from-previous-step>
```

#### 镜像制作

