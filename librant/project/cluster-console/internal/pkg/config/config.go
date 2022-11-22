package config

import (
	"io/ioutil"

	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

// clientSetMap key 为 contextName
var clientSetMap map[string]kubernetes.Interface

// Init 初始化
func Init(kubeConfig string) error {
	// 读取文件中的信息
	content, err := ioutil.ReadFile(kubeConfig)
	if err != nil {
		return err
	}
	var kubeClusterConfig v1.Config
	if err := yaml.Unmarshal(content, &kubeClusterConfig); err != nil {
		return err
	}
	clientSetMap = make(map[string]kubernetes.Interface)
	// 支持多集群登录，遍历 Contexts, 设置 current-context, 再生成客户端
	for _, ctx := range kubeClusterConfig.Contexts {
		kubeClusterConfig.CurrentContext = ctx.Name
		yaml.Unmarshal()
	}

}
