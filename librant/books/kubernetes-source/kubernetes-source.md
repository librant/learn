
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
  - CRI client：kubelet 实现 grpc 的客户端
  - CRI server：dockerd/containerd 的服务端
  - OCI Runtime：真正容器执行 拉起、销毁 等动作
    - runc：OCI 标准的一个具体实现
      - 直接与 cgroup/namespace kernel 进行交互
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
- deletecollection：批量删除资源对象，kubectl 不支持该 verb， 在 rbac 中使用
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

```shell
staging/src/k8s.io/apimachinery/pkg/runtime/schema/group_version.go
```
- GroupVersionResource
- GroupVersion
- GroupResource
- GroupVersionKind
- GroupKind
- GroupVersions

Group/Version/Resource 核心数据结构详情：
- APIGroup
  - typeMeta
  - Name string：资源名称
  - Versions []GroupVersionForDiscovery：资源组下支持的资源版本
  - PreferredVersion GroupVersionForDiscovery：首选版本
  - ServerAddressByClientCIDRs []ServerAddressByClientCIDR

- 将众多资源按照功能划分成不同的资源组，并允许单独启用和禁用资源组
- 支持不同的资源组中拥有不同的资源版本
- 支持同名的资源种类（Kind）存在不同的资源组中
- 资源组和资源版本通过 API Server 对外暴露，允许开发者通过 HTTP 协议进行交互并通过动态客户端进行资源发现
- 支持 CRD 自定义资源扩展
- kubectl 可以不用填资源组名称

- APIVersions
  - TypeMeta
  - Versions []string：所支持的资源版本列表
  - ServerAddressByClientCIDRs []ServerAddressByClientCIDR

kubernetes 的资源版本控制可以分为 3 种：
- Alpha：内部测试
- Beta：特定用户群来进行测试
- Stable：稳定运行版本

- APIResource
  - Name string：资源名称
  - SingularName string：资源的单数名称
  - Namespaced bool：资源是否拥有所属的命名空间
  - Group string：资源所在的资源组名称
  - Version string：资源所在的资源版本
  - Kind String：资源种类
  - Verbs ：资源可操作的资源列表
  - ShortNames []string：资源的简称
  - Categories []string
  - StorageVersionHash string

资源控制系统 -- 管理、调度资源并维护资源状态

Resource Object: 一个资源实例后会表达为一个资源对象，--》Entity
- 持久性实体（Persistent Entity）：在资源对象被创建后，会持久确保该资源对象的存在
  - Deployment
- 短暂性实体（Ephemeral Entity）：在资源对象被创建后，如果出现故障或者调度失败，不会重新创建该资源对象
  - Pod

---
资源外部版本与内部版本（pkg/apis/）
- External Version：
  - External Object：外部资源对象，用于给外部用户请求的接口所使用的资源对象
```shell
<group>/<version>/apps/{v1, v1beta1}/
```  
- Internal Version：
  - Internal Object: 内部版本，不对外暴露，在 API Server 内部使用，用于多资源版本的转换
```shell
<group>/apps/__internal/

staging/src/k8s.io/apimachinery/pkg/runtime/interfaces.go
APIVersionInternal = "__internal" // 内部版本标识
```

---
资源对象描述定义：
- Group/Version
  - apiVersion：指定资源对象的资源组和资源版本
- Kind
  - kind：指定创建资源的种类
- Metadata
  - metadata：描述创建资源对象的元数据信息（名称，命名空间等）
- Spec
  - spec：包含有关资源对象的核心信息
- Status
  - status：正在运行的资源对象的状态信息

---
kubernetes 内置资源：
- apiextensions.k8s.io
  - CustomResourceDefinition：自定义资源类型（APIExtensionServer 负责管理）
- apiregistration.k8s.io
  - APIService：聚合资源类型（AggregatorServer 负责管理）
- admissionregistration.k8s.io
  - MutatingWebhookConfiguration：变更准入控制器资源类型
  - ValidatingWebhookConfiguration：验证准入控制器资源类型

---
runtime.Object 类型系统基石
```shell
staging/src/k8s.io/apimachinery/pkg/runtime/interfaces.go
```
```go
type Object interface {
	GetObjectKind() schema.ObjectKind
	DeepCopyObject() Object
}
```
```shell
staging/src/k8s.io/apimachinery/pkg/runtime/schema/interfaces.go
```
```go
type ObjectKind interface {
	// SetGroupVersionKind sets or clears the intended serialized kind of an object. Passing kind nil
	// should clear the current setting.
	SetGroupVersionKind(kind GroupVersionKind)
	// GroupVersionKind returns the stored group, version, and kind of an object, or an empty struct
	// if the object does not expose or provide these fields.
	GroupVersionKind() GroupVersionKind
}
```

---
**Scheme 资源注册表**  
kubernetes 系统中的资源类型，提供统一的注册，存储，查询，管理等机制；
- 支持注册多种资源类型，包括内部版本和外部版本
- 支持多种版本转换机制
- 支持不同资源序列化/反序列化机制

Scheme 支持两种资源类型（Type） 的注册
- UnversionedType：无版本资源类型
  - scheme.AddUnversionedTypes() 方法注册
- KnownType：有版本资源类型
  - scheme.AddKnownTypes() 方法注册
  - scheme.AddKnownTypesWithName()：注册 KnownType 资源类型，须指定资源的 kind 资源种类名称

Scheme 资源注册表数据结构
- gvkToType：存储 GVK 与 Type 的映射关系
- typeToGVK：存储 Type 和 GVK 的映射关系
- unversionedTypes：存储 UnversionedType 与 GVK 的映射关系
- unversionedKinds：存储 Kind 名称与 UnversionedType 映射关系

```shell
staging/src/k8s.io/apimachinery/pkg/runtime/scheme.go
```
```go
type Scheme struct {
  // gvkToType allows one to figure out the go type of an object with
  // the given version and name.
  gvkToType map[schema.GroupVersionKind]reflect.Type
  
  // typeToGVK allows one to find metadata for a given go object.
  // The reflect.Type we index by should *not* be a pointer.
  typeToGVK map[reflect.Type][]schema.GroupVersionKind
  
  // unversionedTypes are transformed without conversion in ConvertToVersion.
  unversionedTypes map[reflect.Type]schema.GroupVersionKind
  
  // unversionedKinds are the names of kinds that can be created in the context of any group
  // or version
  // TODO: resolve the status of unversioned types.
  unversionedKinds map[string]reflect.Type
}
```

Scheme 资源注册表的查询方法：
- scheme.KnownTypes：查询注册表中指定 GV 下的资源类型
- scheme.AllKnownTypes：查询注册表中所有 GVK 下的资源类型
- scheme.ObjectKinds：查询资源对象所对应的 GVK， 一个资源对象可能会存在多个 GVK

---
**Codec 编解码器**  
编解码器 与 序列化器
- Serializer  
  - 序列化器：包含序列化和反序列化操作
- Codec
  - 编解码器：包含编码器和解码器

```shell
staging/src/k8s.io/apimachinery/pkg/runtime/interfaces.go
```

Codec 包含三种序列化器  
- jsonSerializer：json 格式序列化/反序列化器 
  - ContentType: application/json
- yamlSerializer：yaml 格式序列化/反序列化器
  - ContentType: application/yaml
- protobufSerializer: pb 格式序列化/反序列化器
  - ContentType: application/vnd.kubernetes.protobuf

每一种序列化器都对资源对象的 metav1.TypeMeta（APIVersion,Kind 字段） 进行验证

---
Converter 资源版本转换器  
在 kubernetes 系统中，同一个资源拥有多个资源版本，kubernetes 系统允许同一资源的不同资源版本进行转换；  
-- kubectl convert 命令进行转换
```shell
v1alpha1 --> __internal --> v1beta1/v1
```

Converter 数据结构  
```shell
staging/src/k8s.io/apimachinery/pkg/conversion/converter.go
```
```go
// Converter knows how to convert one type to another.
type Converter struct {
	// Map from the conversion pair to a function which can do the conversion.
	conversionFuncs          ConversionFuncs
	generatedConversionFuncs ConversionFuncs

	// Set of conversions that should be treated as a no-op
	ignoredUntypedConversions map[typePair]struct{}
}
```
- conversionFuncs：默认转换函数
- generatedConversionFuncs：自动生成的转换函数
- ignoredUntypedConversions：若资源对象注册此字段，则忽略此资源对象的转换操作

Converter 注册转换函数  
- scheme.AddIgnoredConversionType：注册忽略的资源类型，不会执行转换操作，忽略资源对象的转换操作
- scheme.AddConversionFuncs：注册多个 Conversion Func 转换函数
- scheme.AddConversionFunc：注册单个 Conversion Func 转换函数
- scheme.AddGeneratedConversionFunc: 注册自动生成的转换函数
- scheme.AddFieldLabelConversionFunc：注册字段标签（Field Label）的转换函数

Scheme 资源注册表可以通过两种方式进行版本转换  
- scheme.ConvertToVersion：将传入的（in）资源对象转换成目标（target）资源版本
- scheme.UnsafeConvertToVersion：在转换过程中不深复制资源对象，而是直接对资源对象进行转换操作

```shell
convertVersion() 
--> reflect.TypeOf(in)：获取资源对象的反射类型
--> s.typeToGVK[t]：从资源注册表中查找传入的资源对象的 GVK(kinds)
--> target.KindForGroupVersionKinds(kinds)：从多个 GVK 中选出与目标资源对象匹配的 GVK 
--> s.unversionedTypes[t]：判断传入的资源对象是否属于 Unversioned 类型 
```
