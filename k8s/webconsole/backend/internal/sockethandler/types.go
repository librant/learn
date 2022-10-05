package sockethandler

// WsMessage web socket message
type WsMessage struct {
	MessageType int    // websocket 消息类型
	Data        []byte // websocket 消息体
}


