package client

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client 客户端
type Client struct {
	RestConfig *rest.Config
}

// InitClient 初始化 client
func InitClient(kubeconfig string) (*Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return &Client{
		RestConfig: config,
	}, nil
}
