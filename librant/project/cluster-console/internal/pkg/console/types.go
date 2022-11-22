package console

// Param container login info
type Param struct {
	CurrentContext string // 根据 current-context 来切换多集群登录
	Namespace      string // 命名空间
	Pod            string // pod name
	Container      string // container name
}
