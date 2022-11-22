package console

import (
	"context"
	"net/http"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/librant/learn/librant/project/cluster-console/internal/pkg/config"
)

// Console login 参数
type Console struct {
}

// LinkContainer 容器登陆
func (c *Console) LinkContainer(ctx context.Context, param Param, w http.ResponseWriter, r *http.Request) error {
	restConfig, err := config.GetRestConfig(param.CurrentContext)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
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
		return err
	}

	conn, err :=


	return nil
}
