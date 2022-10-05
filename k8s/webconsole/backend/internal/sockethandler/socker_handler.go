package sockethandler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"k8s.io/klog"
)

var wsUpGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsConnection websocket 链接
type WsConnection struct {
	wsSocket  *websocket.Conn
	inChan    chan *WsMessage
	outChan   chan *WsMessage
	mu        *sync.Mutex
	isClosed  bool
	closeChan chan byte
}

// InitWebSocket 初始化 web socket
func InitWebSocket(w http.ResponseWriter, r *http.Request) (*WsConnection, error) {
	wsConn, err := wsUpGrader.Upgrade(w, r, nil)
	if err != nil {
		klog.Errorf("Upgrade failed: %v", err)
		return nil, err
	}
	wsSocket := &WsConnection{
		wsSocket:  wsConn,
		inChan:    make(chan *WsMessage),
		outChan:   make(chan *WsMessage),
		mu:        &sync.Mutex{},
		closeChan: make(chan byte),
	}

	// 读协程
	go wsSocket.wsReadLoop()

	// 写协程
	go wsSocket.wsWriteLoop()

	return wsSocket, nil
}

//wsReadLoop 读取循环
func (wsConn *WsConnection) wsReadLoop() {
	for {
		msgType, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			klog.Errorf("read message failed: %v", err)
			wsConn.WsClose()
			return
		}
		msg := &WsMessage{
			MessageType: msgType,
			Data:        data,
		}
		select {
		case wsConn.inChan <- msg:
			// 继续读取下一个消息
		case <-wsConn.closeChan:
			// 这里就直接退出
			return
		}
	}
}

// wsWriteLoop 发送循环
func (wsConn *WsConnection) wsWriteLoop() {
	for {
		select {
		case msg := <-wsConn.outChan:
			if err := wsConn.wsSocket.WriteMessage(msg.MessageType, msg.Data); err != nil {
				klog.Errorf("write message failed: %v", err)
				wsConn.wsSocket.Close()
				return
			}
		case <-wsConn.closeChan:
			return
		}
	}
}

// WsWrite 写消息到前端
func (wsConn *WsConnection) WsWrite(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &WsMessage{MessageType: messageType, Data: data}:
		return nil
	case <-wsConn.closeChan:
		return fmt.Errorf("write websocket closed")
	}
}

// WsRead 从前端读消息
func (wsConn *WsConnection) WsRead() (*WsMessage, error) {
	select {
	case msg := <-wsConn.inChan:
		return msg, nil
	case <-wsConn.closeChan:
		return &WsMessage{}, fmt.Errorf("read websocker closed")
	}
}

// WsClose 关闭 websocket
func (wsConn *WsConnection) WsClose() {
	if err := wsConn.wsSocket.Close(); err != nil {
		klog.Errorf("wsConn.wsSocket close failed: %v", err)
		return
	}
	wsConn.mu.Lock()
	defer wsConn.mu.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}
