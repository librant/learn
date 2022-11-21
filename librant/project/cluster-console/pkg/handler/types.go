package handler

// Param container login info
type Param struct {
	CurrentContext string // 支持多集群登陆
	Namespace      string // 命名空间
	Pod            string // pod name
	Container      string // container name
}
