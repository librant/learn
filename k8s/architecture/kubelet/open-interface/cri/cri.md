#### CRI（Container Runtime Interface）

CRI 主要分为三部分：
- CRI Client
- CRI Server
- OCI （Open Container Initiative） Runtime
  - runc：OCI 标准的参考实现，直接依赖 cgroup/namespace kernel 等进行交互，
    负责为容器配置 cgroup/namespace 等启动容器所需的环境，创建容器启动的相关进程

---
Docker 的架构调整：

Docker Engine --> containerd  
--> containerd-shim --> runC
--> containerd-shim --> runC

kubelet <-- CRI <--> CRI-containerd <--> containerd --> container

cri-o：cri 和 oci 之间的一座桥梁




