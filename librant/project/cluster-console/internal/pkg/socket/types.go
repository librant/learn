package socket

import "time"

const (
	ideaTimeout = time.Minute
)

// Message web socket message
type Message struct {
	Type int    // websocket 消息类型, 就是 websocket 管道中传输的是什么类型的数据
	Data []byte // websocket 消息体, 需要根据消息的类型进行解析
}
