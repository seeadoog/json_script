package ws_js

import (
	"fmt"
	"github.com/seeadoog/schema"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestWebsocket(t *testing.T) {
	in:=`{"header":{"app_id":"4CC5779A","status":0,"uid":"77607d78ccc5424986cdbe4039879b14"},"parameter":{"iat":{"ent":"sms16k","iat_text":{"encoding":"gzip"}}},"payload":{"iat_audio":{"audio":"UklGRvQ8AgBXQVZFZm10IBAAAAABAAEAQB8AAIA+AAACABAATElTVBoAAABJTkZPSVNGVA4AAABMYXZmNTguMTcuMTAxAGRhdGGuPAIA6f8cAH4CgABNAJoAewCHAIUAeQB7AJUAggCRAIoAiwCDALcA/AAnAc8A7QAAAcIAlwC2AMMAxACBAIsAygCNAJAAXQBdAFoAaQBhAFEAbwBTAEwAQQBAAB8A/v8zAFQAIgD7//j/HQAbABkAEwAKABUAFQAeAEsAOgBFAFUAPgBaAGUAdQBDAD8ATQAVACgABgDZ/ysAOgACAAwAKgAXACUATwAmAOb/CwDf//7/CgD7/xwAPQBSAEAAMAAxAHoAfgBjAF8AfQBXAE4ARgBGAGoAXACDAJ4AlgCtAJYAhgBqAEAAPQBTAIwAdACHAGsAUgB1AF8AjACKAIcAfAA5AHsAkwA5ABUARgA2AG8AXwAgAD4AOABlAFgANgBhADkAOwBFAFkAfQA4ACwAOgAoADcADgAOABQABgAqAEcAPAB6AEwAFwAlADIAJAAIABsAKgAVABoAHgARAOb/4v+9/7v/3v/r//H/DwDd/7//n/+f/3j/S/89/1H/bf9n/5j/vP/a/+3/2//N/8b/nv+H/5L/p/+r/7//r/+6/7P/mv/F/9L/sP/I/9j/4f/W/5z/kf+i/5P/g/+K/5f/s/+Z/7X/hf+B/2X/YP9l/1n/b/92/4D/kf+A/yb/Af8V/xn/4f4i/0L/Mf9D/0L/Ff8q/zj/Wv9m/2T/jv9l/3f/XP9X/1L/af+k/4X/r//M/9D/2//7/wgACQAtAOv/tf+8/6f/rv+w/9L/6P/W/+v/2P/o//j/GwD+/wcA7v/G/+r/mP/C//n/7P/v/xIAHwAeAAsA/v8CAOP/DQDp//P//P/1/xIAJgAUABgAEQAIAAAADwD9/+n/2//o/9f/zP/G/9X/4f/Z/+H/EgAMAB4AOAAgADUAEADi//n/4/+7/9n//f8CAPT/LQANAAAACgC+/7j/2P/E/8H/sv/A/+b/4//l/xEA4v+z/9n/n/+H/4b/T/9T/2b/Tv9k/3f/ev9//5n/j/+X/23/Q/9r/2L/Tf80/yz/Ff8y/yr/P/86/zf/R////hv/F/9I/1//Tf+V/4P/gv+n/6b/j/+Q/6P/s//F/8P/zf+j/5r/tv+4/7//rf/m/w4ABwA1ABcAWwBSACwAUwA/AGMAKwAuADQAWQBpADgALwCAAHkAMwA5ABsANQA/ACcAFwD4//3/8f8HACMALwAyAA0ADQAWAAoADQDv/9r/9P/r/w0AVwCCAHEAeQB4AHAAgwB/AI4AogBsAF4AfQBtAEYAcABlAE8AdQBWAEUAXQBjAD4AVQByAG0AUwBIAEYAMwBKAFAAZwCKAIgAiwB9AD4AOgBMAEkAQQAfADoAQAA7ADMALQAeADMAGAAVAOn/8v/1/w0A+P/M//D/+v/a//3//P/s//j/4f+7//H//P/R/+v/0P+2/7T/g/+a/6n/tf/c/+P/1//J/+n/0//E/+f/6v/+/x8ABwAvAD4AMwA3AFAAZwB5ALIAaABAADkAPwBKAIIAdwCDAHIAUABIAFIAHwAIACoAJAA8ADgAOwA6AEwAKgBlAHoALQA7AO7/pP/H/6//y//H/9X/xf/y//P/0v8CAOD/vv/B/97/8v/k/9L/5f/y/zEAPAA=","encoding":"raw","sample_rate":"16000","seq":0,"status":0}}}
`
	out,_:=schema.GenerateSchemaFromJson([]byte(in))
	lua.ApiError{}
	fmt.Println(string(out))

}
