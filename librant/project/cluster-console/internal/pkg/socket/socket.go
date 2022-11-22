package socket

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"k8s.io/klog"
)

var wsUpGrader = websocket.Upgrader{}

// Connection websocket 链接
type Connection struct {
	wsSocket  *websocket.Conn
	inChan    chan *Message
	outChan   chan *Message
	mu        *sync.Mutex
	isClosed  bool
	closeChan chan byte // websocket 关闭通知
	ideaTime  *time.Timer
}

func InitSocket(ctx context.Context, w http.ResponseWriter, r *http.Request) (*Connection, error) {
	conn, err := wsUpGrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	wsSocket := &Connection{
		wsSocket:  conn,
		inChan:    make(chan *Message),
		outChan:   make(chan *Message),
		mu:        &sync.Mutex{},
		closeChan: make(chan byte),
		ideaTime:  time.NewTimer(ideaTimeout),
	}

	// 消息读取协程（页面 --》 k8s）
	go func(ctx context.Context) {
		defer wsSocket.WebsocketClose()
		if err := wsSocket.readLoop(ctx); err != nil {
			klog.Errorf("wsReadLoop failed: %v", err)
		}
	}(ctx)

	go func(ctx context.Context) {
		defer wsSocket.WebsocketClose()
		if err := wsSocket.writeLoop(ctx); err != nil {
			klog.Errorf("wsWriteLoop failed: %v", err)
		}
	}(ctx)

	return wsSocket, nil
}

// readLoop 从页面读取通道中的信息
func (c *Connection) readLoop(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("readLoop ctx error: %v", err)
		}
		msgType, data, err := c.wsSocket.ReadMessage()
		if err != nil {
			return fmt.Errorf("ReadMessage failed: %v", err)
		}
		msg := &Message{
			Type: msgType,
			Data: data,
		}
		select {
		case <-ctx.Done():
			return nil

		case c.inChan <- msg:
			// 从前端正确读取消息，继续等待下一次读取
			c.ideaTime.Reset(ideaTimeout)

		case <- c.closeChan:
			klog.Infof("readLoop chan closed")
			return nil
		}
	}
}

// writeLoop 循环读取 k8s 集群中的消息，发送到前端
func (c *Connection) writeLoop(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("writeLoop ctx error: %v", err)
		}
		select {
		case <-ctx.Done():
			return nil

		case msg := <-c.outChan:
			// 从 k8s exec 中读取消息，发送到前端
			if err := c.wsSocket.WriteMessage(msg.Type, msg.Data); err != nil {
				return fmt.Errorf("WriteMessage failed: %v", err)
			}
		case <-c.closeChan:
			klog.Infof("writeLoop chan closed")
			return nil
		}
	}
}

// WebsocketWrite 写消息到前端
func (c *Connection) WebsocketWrite(msgType int, data []byte) error {
	select {
	case c.outChan <- &Message{Type: msgType, Data: data}:
		return nil
	case <-c.closeChan:
		return fmt.Errorf("write websocket closed")
	}
}

// WebsocketRead 从前端读消息
func (c *Connection) WebsocketRead() (*Message, error) {
	select {
	case msg := <-c.inChan:
		return msg, nil
	case <-c.closeChan:
		return nil, fmt.Errorf("read websocker closed")
	}
}

// WebsocketClose 关闭 websocket 通道
func (c *Connection) WebsocketClose() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isClosed {
		c.isClosed = true
		close(c.closeChan)
		close(c.inChan)
		close(c.outChan)
		_ = c.wsSocket.Close()
	}
}
