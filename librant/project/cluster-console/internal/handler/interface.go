package handler

import (
	"context"

	"k8s.io/client-go/kubernetes"
)

//go:generate mockgen -destination handler_mock.go -source handler.go -package handler

// ClusterClient 获取集群 client 接口
type ClusterClient interface {
	GetClientSet(ctx context.Context, curCtxName, namespace string) (kubernetes.Interface, error)
}

// ClusterConsole 集群容器登录接口
type ClusterConsole interface {

}


