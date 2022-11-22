package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog"

	"github.com/librant/learn/librant/project/cluster-console/internal/pkg/config"
	"github.com/librant/learn/librant/project/cluster-console/internal/pkg/signals"
	"github.com/librant/learn/librant/project/cluster-console/internal/pkg/socket"
	"github.com/librant/learn/librant/project/cluster-console/internal/pkg/stream"
)

// IndexHandler index handler
func IndexHandler(c *gin.Context) {
	// 获取 GET 请求中的参数信息
	param, err := getParam(c.Request)
	if err != nil {
		setResponse(c, http.StatusBadRequest, err)
		return
	}
	klog.Infof("IndexHandler param: %v", param)

	restConfig, err := config.GetRestConfig(param.CurrentContext)
	if err != nil {
		setResponse(c, http.StatusBadRequest, err)
		return
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		setResponse(c, http.StatusBadRequest, err)
		return
	}
	execReq := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(param.Pod).
		Namespace(param.Namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: param.Container,
			Command:   []string{"/bin/bash", "-c", "export TMOUT=600; bash"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(restConfig, http.MethodPost, execReq.URL())
	if err != nil {
		setResponse(c, http.StatusBadRequest, err)
		return
	}

	ctx := signals.GracefulStopWithContext()
	conn, err := socket.InitSocket(ctx, c.Writer, c.Request)
	if err != nil {
		setResponse(c, http.StatusBadRequest, err)
		return
	}
	wsStream := stream.New(conn)
	opts := remotecommand.StreamOptions{
		Stdin:             wsStream,
		Stdout:            wsStream,
		Stderr:            wsStream,
		TerminalSizeQueue: wsStream,
		Tty:               true,
	}
	if err := exec.Stream(opts); err != nil {
		setResponse(c, http.StatusInternalServerError, err)
		return
	}
	return
}

// getParam 获取参数信息
func getParam(r *http.Request) (Param, error) {
	if err := r.ParseForm(); err != nil {
		return Param{}, err
	}
	return Param{
		CurrentContext: r.Form.Get("context"),
		Namespace:      r.Form.Get("namespace"),
		Pod:            r.Form.Get("pod"),
		Container:      r.Form.Get("container"),
	}, nil
}

func setResponse(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{
		"message": fmt.Sprintf("%v", err),
	})
}
