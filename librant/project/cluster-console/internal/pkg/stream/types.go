package stream

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	xtermResizeType = "resize"
	xtermPingType   = "ping"
	xtermInputType  = "input"
)

// xtermMessage xterm 消息格式
type xtermMessage struct {
	MsgType string `json:"type"`
	Input   string `json:"input"`
	Rows    uint16 `json:"rows"`
	Cols    uint16 `json:"cols"`
}
