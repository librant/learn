package streamhandler

import (
	"encoding/json"
	"k8s.io/klog"
	"net/http"

	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/librant/learn/k8s/webconsole/backend/internal/client"
	"github.com/librant/learn/k8s/webconsole/backend/internal/sockethandler"
)

// StreamHandler 流式处理
type StreamHandler struct {
	WsConn      *sockethandler.WsConnection
	ResizeEvent chan remotecommand.TerminalSize
}

// Next 处理页面变化
func (s *StreamHandler) Next() *remotecommand.TerminalSize {
	ret := <-s.ResizeEvent
	return &ret
}

// Write 写处理
func (s *StreamHandler) Write(p []byte) (int, error) {
	size := len(p)
	copyData := make([]byte, size)
	copy(copyData, p)
	return size, s.WsConn.WsWrite(websocket.TextMessage, copyData)
}

// Read 读处理
func (s *StreamHandler) Read(p []byte) (int, error) {
	msg, err := s.WsConn.WsRead()
	if err != nil {
		return 0, err
	}
	// 解析对应的 msg 信息
	var xtermMsg xtermMessage
	if err := json.Unmarshal(msg.Data, &xtermMsg); err != nil {
		return 0, err
	}
	if xtermMsg.MsgType == "resize" {
		s.ResizeEvent <- remotecommand.TerminalSize{
			Width:  xtermMsg.Cols,
			Height: xtermMsg.Rows,
		}
	}
	if xtermMsg.MsgType == "input" {
		size := len(xtermMsg.Input)
		copy(p, xtermMsg.Input)
		return size, nil
	}
	return 0, nil
}

// WebConsoleLink web console 链接
func WebConsoleLink(kubeconfig string, w http.ResponseWriter, r *http.Request) error {
	cli, err := client.InitClient(kubeconfig)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(cli.RestConfig)
	if err != nil {
		return err
	}
	execReq := clientset.CoreV1().
		RESTClient().
		Post().
		Resource("pods").
		Name("pod-name").
		Namespace("ns-name").
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: "ctn-name",
			Command:   []string{"/bin/bash"},
			Stderr:    true,
			Stdout:    true,
			Stdin:     true,
			TTY:       true,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(cli.RestConfig, http.MethodPost, execReq.URL())
	if err != nil {
		return err
	}
	wsConn, err := sockethandler.InitWebSocket(w, r)
	if err != nil {
		return err
	}
	handler := &StreamHandler{
		WsConn:      wsConn,
		ResizeEvent: make(chan remotecommand.TerminalSize),
	}
	if err := exec.Stream(remotecommand.StreamOptions{
		Stdin:  handler,
		Stdout: handler,
		Stderr: handler,
		TerminalSizeQueue: handler,
		Tty: true,
	}); err != nil {
		klog.Errorf("stream failed: %v", err)
		wsConn.WsClose()
		return err
	}
	return nil
}
