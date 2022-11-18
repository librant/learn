// 资源组、资源版本、资源

---
Pod 对象创建流程：

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

