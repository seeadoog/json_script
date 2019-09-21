package ws_js


import (
	"encoding/json"
	"fmt"
	jsonscpt "git.xfyun.cn/AIaaS/json_script"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"time"
)
func init(){
	jsonscpt.SetInlineFunc("websocket_open", websocket_open)
	jsonscpt.SetInlineFunc("websocket_write", websocket_write)
	jsonscpt.SetInlineFunc("websocket_close", websocket_close)
	jsonscpt.SetInlineFunc("sleep", sleep)
	jsonscpt.SetInlineFunc("boce_info", boce_info)
}
func newWebsocket(url string,onOpen,onMessage,onClose jsonscpt.Func)*websocket.Conn{
	d:=websocket.Dialer{}
	con,rsp,err:=d.Dial(url,nil)
	if err !=nil{
		panic(err.Error()+readResp(rsp))
		return nil
	}
	onOpen(readResp(rsp))
	go func() {
		defer func() {
			if err :=recover();err !=nil{
				onClose(err)
			}
		}()
		for{
			_,msg,err:=con.ReadMessage()
			if err !=nil{
				onClose(err)
				return
			}
			var r interface{}
			err=json.Unmarshal(msg,&r)
			if err !=nil{
				onMessage(msg)
				continue
			}
			onMessage(r)
		}
	}()

	return con
}

func websocket_open(i ...interface{})interface{}{
	if len(i)<4{
		panic("open websocket ,args  < 4")
	}
	return newWebsocket(jsonscpt.String(i[0]),i[1].(jsonscpt.Func),i[2].(jsonscpt.Func),i[3].(jsonscpt.Func))
}

func websocket_write(i ...interface{})interface{}{
	if len(i)<2{
		panic("websocket_write args <2")
	}
	err:=i[0].(*websocket.Conn).WriteJSON(i[1])
	if err !=nil{
		panic("websocket write error")
	}
	return nil
}
func websocket_close(i ...interface{})interface{}{
	if len(i)<1{
		panic("websocket_write args <1")
	}
	i[0].(*websocket.Conn).Close()

	return nil
}
func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}

func sleep(i ...interface{})interface{}{
	if len(i)<1{
		panic("sleep args <1 ")
	}
		time.Sleep(time.Duration(jsonscpt.Number(i[0]))*time.Millisecond)
	return nil
}

func boce_info(i ...interface{})interface{}{
	if len(i) <3{
		panic("print_info args <3")
	}
	in:= map[string]interface{}{
		"code":i[0],
		"info":i[1],
		"time":i[2],
	}
	b,_:=json.Marshal(in)
	return string(b)
}