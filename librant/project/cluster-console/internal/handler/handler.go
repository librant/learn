package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"

	"github.com/librant/learn/librant/project/cluster-console/internal/pkg/console"
)

// Handler 容器登录处理
type Handler struct {
	client ClusterClient
}

// New 生成实例
func New(client ClusterClient) *Handler {
	return &Handler{
		client: client,
	}
}

// IndexHandler index handler
func IndexHandler(c *gin.Context) {
	// 获取 GET 请求中的参数信息
	param, err := getParam(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("get param failed: %v", err),
		})
		return
	}
	klog.Infof("IndexHandler param: %v", param)

	return
}

// getParam 获取参数信息
func getParam(r *http.Request) (console.Param, error) {
	if err := r.ParseForm(); err != nil {
		return console.Param{}, err
	}
	return console.Param{
		CurrentContext: r.Form.Get("context"),
		Namespace:      r.Form.Get("namespace"),
		Pod:            r.Form.Get("pod"),
		Container:      r.Form.Get("container"),
	}, nil
}
