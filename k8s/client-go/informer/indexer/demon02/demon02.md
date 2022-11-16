[client-go 之 Indexer 的理解](https://cloud.tencent.com/developer/article/1692517)

```shell
IndexFunc：索引器函数，用于计算一个资源对象的索引值列表，上面示例是指定命名空间为索引值结果，当然我们也可以根据需求定义其他的，比如根据 Label 标签、Annotation 等属性来生成索引值列表。
Index：存储数据，对于上面的示例，我们要查找某个命名空间下面的 Pod，那就要让 Pod 按照其命名空间进行索引，对应的 Index 类型就是 map[namespace]sets.pod。
Indexers：存储索引器，key 为索引器名称，value 为索引器的实现函数，上面的示例就是 map["namespace"]MetaNamespaceIndexFunc。
Indices：存储缓存器，key 为索引器名称，value 为缓存的数据，对于上面的示例就是 map["namespace"]map[namespace]sets.pod
```

```shell
// Indexers 就是包含的所有索引器(分类)以及对应实现
Indexers: {  
  "namespace": NamespaceIndexFunc,
  "nodeName": NodeNameIndexFunc,
}
// Indices 就是包含的所有索引分类中所有的索引数据
Indices: {
 "namespace": {  //namespace 这个索引分类下的所有索引数据
  "default": ["pod-1", "pod-2"],  // Index 就是一个索引键下所有的对象键列表
  "kube-system": ["pod-3"]   // Index
 },
 "nodeName": {  //nodeName 这个索引分类下的所有索引数据(对象键列表)
  "node1": ["pod-1"],  // Index
  "node2": ["pod-2", "pod-3"]  // Index
 }
}
```