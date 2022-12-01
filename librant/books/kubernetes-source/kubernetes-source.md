
### 第一章 Kubernetes 架构

- 可移植
- 可扩展
- 自动化

Kubernetes 总体架构图：

**Master 组件：**  
- kube-apiserver: 集群的 HTTP REST API 接口，是集群的控制入口
- kube-scheduler：集群中 Pod 对象的调度服务
- kube-controller-manager：集群中所有资源对象的自动化控制中心

**Node 组件：**
- kubelet：负责管理节点上的容器创建，删除，启停等任务
- kube-proxy：负责 kubernetes 服务的通信及负载均衡服务
- containerd/dockerd：负责容器的基础服务关联

**其他组件：**
- kubectl：官方提供的命令行工具
- etcd
- kube-coredns

---
kubelet 实现三种开放接口：
- ContainerRuntimeInterface: (CRI)
  CRI 将 kubelet 组件与容器运行时进行解耦
- ContainerNetworkInterface: (CNI)  
  容器创建时，通过 CNI 插件配置网络
- ContainerStorageInterface: (CSI)
  容器创建时，通过 CSI 插件配置存储接口

---
Kubernetes Project Layout 设计
- cmd：存放可执行文件的入口，每个可执行文件对应一个 main 函数
- pkg：存放核心库代码，可以被项目内外直接引用
- vendor：存放项目依赖的库代码， 一般为第三方库代码
- api：存放 OpenAPI/Swagger 的 spec 文件， 包括 json，protobuf 的定义等
- build：存放与构建相关的脚本
- test：存放测试工具和测试数据
- docs：存放设计或用户使用文档
- hack：存放与构建、测试等相关的脚本
- third_party：存放第三方工具，代码或者其他组件
- plugin：存放 kubernetes 插件代码目录，例如认证，授权等相关插件
- staging：存放部分核心库的暂存目录
- translations：存放 i18n 国际化语言包相关文件

---
初始化过程：  
main --> options --> flag --> log --> options --> execute

---
### 第三章 Kubernetes 核心数据结构

Kubernetes 是一个完全以资源为中心的系统  
-- 注册、管理、调度资源并维护资源的状态

#### 3.1 Group、Version、Resource 核心数据结构

- Group：被称为资源组，在 Kubernetes API Server 中也可以称为 APIGroup
- Version：被称为资源版本，在 kubernetes API Server 中也可称其为 APIVersions
- Resource：被称为资源，在 Kubernetes API Server 中也可称其为 APIResource
- Kind：资源种类，描述 Resource 的种类，与 Resource 为同一级别

Kubernetes 系统支持多个 Group， 每个 Group 支持多个 Version，每个 Version 支持多个 Resource，
其中部分资源同时会拥有自己的子资源（SubResource）  

资源组、资源版本、资源、子资源的完整表现形式为：
```shell
<group>/<version>/<resource>/<subresource>
```

资源对象（Resource Object）： 由 资源组 + 资源版本 + 资源种类 组成
```shell
<group>/<version>, Kind=<kind>
apps/v1, Kind=Deployment
```

每一个资源都拥有一定数量的资源操作方法（Verbs），用于 etcd 中存储的资源对象的 增/删/改/查 操作；  
Kubernetes 系统支持 8 种资源操作方法：
- create
- delete
- deletecollection
- get
- patch
- update
- list
- watch

每一个资源至少有两个版本： 
- 外部版本（External Version）  
  对外暴露给用户请求的接口所使用的资源对象
- 内部版本（Internal Version）
  内部版本不对外暴露，仅在 Kubernetes API Server 内部使用  

Kubernetes 资源也可分为两种：
- 内置资源（Kubernetes Resource）
- 自定义资源（Custom Resource）  
  CRD 可以实现自定义资源，允许用户将自己定义的资源添加到 Kubernetes 系统中

Group/Version/Resource 核心数据结构存放在：
```shell
k8s.io/apimachinery/pkg/apis/meta/v1/types.go
```
```go
type APIResourceList struct {
	TypeMeta `json:",inline"`
	// groupVersion is the group and version this APIResourceList is for.
	GroupVersion string `json:"groupVersion" protobuf:"bytes,1,opt,name=groupVersion"`
	// resources contains the name of the resources and if they are namespaced.
	APIResources []APIResource `json:"resources" protobuf:"bytes,2,rep,name=resources"`
}
```

```shell
k8s.io/apimachinery/pkg/runtime/schema/group_version.go
```
```go
type GroupVersionResource struct {
	Group    string
	Version  string
	Resource string
}
```










