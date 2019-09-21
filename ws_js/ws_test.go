package ws_js

import (
	"fmt"
	jsonscpt "git.xfyun.cn/AIaaS/json_script"
	"testing"
)

func TestWebsocket(t *testing.T) {
	defer func() {
		if err := recover();err !=nil{
			fmt.Println(err)
		}
	}()
	vm := jsonscpt.NewConcurrencyVm()
	err := vm.ExecFile("ws.json")
	if err != nil {
		panic(err)
	}
}
