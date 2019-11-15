package ws_js

import (
	"encoding/json"
	"fmt"
	jsonscpt "git.xfyun.cn/AIaaS/json_script"
	"io/ioutil"
	"testing"
)

func TestWebsocket(t *testing.T) {
	defer func() {
		if err := recover();err !=nil{
			fmt.Println(err)
		}
	}()
	vm := jsonscpt.NewVm()
	vm.Set("$",map[string]interface{}{
		"app_id":"100IME",
	})
	b,err:=ioutil.ReadFile(`ws.json`)
	if err !=nil{
		panic(err)
	}
	var s  = &jsonscpt.Script{}
	err =json.Unmarshal(b,s)
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Exec(vm))

}
