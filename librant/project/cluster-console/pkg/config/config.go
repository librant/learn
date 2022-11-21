package config

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

var clientSet kubernetes.Interface

// Init  初始化配置
func Init(kubeConfig string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	clientSet = clientset
	return nil
}

// GetClientSet 获取 clientSet 信息
func GetClientSet() kubernetes.Interface {
	if clientSet == nil {
		klog.Fatalln(fmt.Sprintf("clientSet not init"))
	}
	return clientSet
}