package console

import (
	"context"
	"net/http"

	"k8s.io/client-go/kubernetes"
)

// Console login 参数
type Console struct {
	clientSet kubernetes.Interface
}

// LinkContainer 容器登陆
func (c *Console) LinkContainer(ctx context.Context, param Param, w http.ResponseWriter, r *http.Request) error {
	return nil
}
