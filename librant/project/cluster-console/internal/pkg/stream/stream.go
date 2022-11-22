package stream

import (
	"fmt"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/librant/learn/librant/project/cluster-console/internal/pkg/socket"
)

//go:generate mockgen -destination stream_mock.go -source stream.go -package stream

// WebSocket 操作接口
type WebSocket interface {
	WebsocketWrite(msgType int, data []byte) error
	WebsocketRead() (*socket.Message, error)
}

// Stream 流式处理
type Stream struct {
	wsConn      WebSocket
	resizeEvent chan remotecommand.TerminalSize
}

// New 生成实例
func New(wsConn WebSocket) *Stream {
	return &Stream{
		wsConn:      wsConn,
		resizeEvent: make(chan remotecommand.TerminalSize),
	}
}

// Next 处理页面变化
func (s *Stream) Next() *remotecommand.TerminalSize {
	ret := <-s.resizeEvent
	return &ret
}

// Write 将 k8s exec 的数据写入到 websocket 的管道中，返回给前端页面显示
func (s *Stream) Write(p []byte) (int, error) {
	size := len(p)
	copyData := make([]byte, size)
	copy(copyData, p)
	return size, s.wsConn.WebsocketWrite(websocket.BinaryMessage, copyData)
}

// Read 将前端页面的数据读入到 websocket 的管道中，写入 k8s exec
func (s *Stream) Read(p []byte) (int, error) {
	msg, err := s.wsConn.WebsocketRead()
	if err != nil {
		return 0, err
	}
	switch msg.Type {
	case websocket.BinaryMessage:
		// 如果是二进制数据，则表示消息，这里直接传输
		copy(p, msg.Data)

	case websocket.TextMessage:
		// 这里也可以切换成 text 的方式传输
		var xtermMsg xtermMessage
		if err := json.Unmarshal(msg.Data, &xtermMsg); err != nil {
			return 0, err
		}
		if xtermMsg.MsgType == xtermResizeType {
			s.resizeEvent <- remotecommand.TerminalSize{
				Width:  xtermMsg.Cols,
				Height: xtermMsg.Rows,
			}
			return len(xtermMsg.Input), nil
		}
		if xtermMsg.MsgType == xtermPingType {
			return 0, nil
		}
		if xtermMsg.MsgType == xtermInputType {
			// 当使用 text 的方式传输时，这里生效
			copy(p, msg.Data)
			return len(xtermMsg.Input), nil
		}

	case websocket.CloseMessage:
		return 0, fmt.Errorf("websocket has been closed")
	}
	return len(msg.Data), nil
}


