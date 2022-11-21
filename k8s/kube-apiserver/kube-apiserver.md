// 资源组、资源版本、资源

---
#### Pod 对象创建流程：

1、kubectl --> kube-apiserver 提交 pod 资源清单；

2、kube-apiserver 验证请求 --> etcd 持久化；

3、kube-apiserver 基于 watch 机制 --> kube-scheduler

4、kube-scheduler 基于 预选和优选 调度算法 为 pod 选择最优节点 --> kube-apiserver

5、kube-apiserver 将最优节点 --> etcd
    将 资源配置清单中补充 nodeName 字段

6、kube-apiserver 基于 watch 机制 --> kubelet

7、kubelet --> CSI/CRI/CNI --> (dockerd/containerd) --> pod(container)

8、kubelet --> 上报容器的 pod/status --> kube-apiserver

9、kube-apiserver --> 更新 pod/status 子资源 --> etcd

---
#### HTTP 请求的完整生命周期：

1、用户向 kube-apiserver 发出 http 请求

2、kube-apiserver 接收到请求，启动 gorouting 处理接收到的请求

3、kube-apiserver 验证请求内容中的认证（auth）信息

4、kube-apiserver 解析请求内容

5、kube-apiserver 调用路由项对应的 handle 回调函数

6、kube-apiserver 获取 Handle 回调函数的数据信息

7、kube-apiserver 设置请求状态码，响应用户的请求

---
#### API Server 启动流程：

1、资源注册

2、Cobra 命令行参数解析

3、创建 API Server 通用参数配置

4、创建 APIExtensionsServer、KubeAPIServer、AggregatorServer

5、创建 GenericAPIServer

6、启动 Http、Https 服务

---
#### 资源注册

1、初始化 Scheme 资源注册表 <br>
2、注册 k8s 所支持的资源 <br>

```shell
pkg/api/legacyscheme/scheme.go

var Scheme = runtime.NewSchema()
var Codecs = serializer.NewCodecFactory(Scheme)
var ParameterCodec = runtime.NewParameterCodec(Scheme)
```

---
#### 创建 APIExtensionsServer

1、创建 GenericAPIServer 

2、实例化 CustomResourceDefinitions

3、实例化 APIGroupInfo

4、注册 APIGroup (InstallAPIGroup)

---




