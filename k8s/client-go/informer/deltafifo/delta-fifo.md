### DeltaFIFO

1、增量的本地队列，记录对象的变化过程
2、Delta 两个属性：<br>
1）Type <br>
2) Object <br>

3、
```go
items map[string]Deltas: 每个 key 的变化
queue []string
keyFunc KeyFunc

func MetaNamespaceKeyFunc(obj interface{}) (string, error) {」
```

