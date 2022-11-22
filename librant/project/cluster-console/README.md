### cluster-console 实现多集群登陆

Cobra + Gin + websocket

1、读取 kubeconfig 文件信息，根据传入的 clusterName 切换指定的 context；

2、初始化 kubeconfig 中的集群列表，生成对应的 clientset 客户端；

3、这里需要根据传入参数：
- context: 上下文
- namespace: 命名空间
- pod: pod 名
- container: 容器名

