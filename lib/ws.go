package lib

import (
	jsonscpt "git.xfyun.cn/AIaaS/json_script"
	"golang.org/x/net/websocket"
)

func newWebsocket(url string,onOpen,onMessage,onClose,onFail jsonscpt.Func){
	websocket.Dial()
}
