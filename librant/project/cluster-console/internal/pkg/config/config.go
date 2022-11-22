package config

import (
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/rest"

	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/tools/clientcmd"
	v1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

// clientSetMap key 为 contextName
var restConfigMap = make(map[string]*rest.Config)

// Init 初始化
func Init(kubeConfig string) error {

	// TODO: 支持多集群登录，遍历 Contexts, 设置 current-context, 再生成客户端

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return err
	}
	// 读取文件中的信息, 保存当前 current-context 的信息
	content, err := ioutil.ReadFile(kubeConfig)
	if err != nil {
		return err
	}
	var kubeClusterConfig v1.Config
	if err := yaml.Unmarshal(content, &kubeClusterConfig); err != nil {
		return err
	}
	restConfigMap[kubeClusterConfig.CurrentContext] = config
	return nil
}

// GetRestConfig 根据 CurrentContext 获取当前的 rest config 文件
func GetRestConfig(cur string) (*rest.Config, error) {
	if config, ok := restConfigMap[cur]; ok {
		return config, nil
	}
	return nil, fmt.Errorf("%s current context has not beed set", cur)
}
