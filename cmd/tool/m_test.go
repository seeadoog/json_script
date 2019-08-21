package main

import (
	"fmt"
	"git.xfyun.cn/AIaaS/json_script"
	"github.com/robertkrimen/otto"
	"testing"
	"time"
)
var paramMap=map[string]interface{}{
	"tts":"cont",
	"td":"tre",
	"tdds":"tre",
	"a":"e",
	"b":"s",
	"subappid":"34345",
	"ent":"s",
}
func TestJs(t *testing.T){

	vm:=otto.New()

	vm.Set("$",paramMap)
	cp,err:=vm.Compile("",`
	if ($['tts']==='cont'){
		$['val']='volmap'
		$['single']=10000
		$['sdk']=$['a']+$['b']+'_16k'
	}
	if ($['ent']==''){

	}else{
	$['appid']=$['subappid']
	$['appid']=$['subappid']
	$['appid']=$['subappid']
	$['appid']=$['subappid']
	$['appid']=$['subappid']
	$['appid']=$['subappid']
	$['appid']=$['subappid']
	$['appid']=$['subappid']
	$['appid']=$['subappid']
	$['appid']=$['subappid']
	$['appid']=$['subappid']
	$['appid']=$['subappid']
	$['appid']=$['subappid']
	$['appid']=$['subappid']
	$['appid']=$['subappid']
}
`)
	st:=time.Now()
	for i:=0;i<100000;i++{

		vm.Run(cp)
	}
	fmt.Println(vm.Get("str"))
	fmt.Println(err,time.Since(st))
	fmt.Println(paramMap)
}



func TestJson(t *testing.T){

	vm:=jsonscpt.NewVm()

	vm.Set("$",paramMap)
	cp,err:=jsonscpt.CompileExpFromJson([]byte(`
	[
  {
    "if": "eq($.tts,'cont')",
    "then": [
      "$.val='valmap'",
      "$.single=10000",
      "$.sdk=append($.a,$.b,'_16k')"
    ]
  },
	{
		"if":"not($.ent)",
		"then":"return(0,'ent')"
	},
	"$.appid=$.subappid",
	"$.appid=$.subappid",
	"$.appid=$.subappid",
	"$.appid=$.subappid",
	"$.appid=$.subappid",
	"$.appid=$.subappid",
	"$.appid=$.subappid",
	"$.appid=$.subappid",
	"$.appid=$.subappid",
	"$.appid=$.subappid",
	"$.appid=$.subappid",
	"$.appid=$.subappid",
	"$.appid=$.subappid",
	"$.appid=$.subappid",
	"$.appid=$.subappid",
	"$.appid=$.subappid"
]
`))
	st:=time.Now()
	for i:=0;i<100000;i++{

		vm.SafeExecute(cp,nil)
	}
	fmt.Println(vm.Get("str"))
	fmt.Println(err,time.Since(st))
	fmt.Println(paramMap)
}

func TestSo(t *testing.T){
	for i:=0;i<100000;i++{
		if paramMap["tts"]=="cont"{
			paramMap["val"]="volmap"
			paramMap["single"]=10000
			paramMap["sdk"]=paramMap["a"].(string)+paramMap["b"].(string)+"_16k"
		}
		if paramMap["ent"]==""{
			return
		}else{
			paramMap["appid"]=paramMap["subappid"]
			paramMap["appid"]=paramMap["subappid"]
			paramMap["appid"]=paramMap["subappid"]
			paramMap["appid"]=paramMap["subappid"]
			paramMap["appid"]=paramMap["subappid"]
			paramMap["appid"]=paramMap["subappid"]
			paramMap["appid"]=paramMap["subappid"]
			paramMap["appid"]=paramMap["subappid"]
			paramMap["appid"]=paramMap["subappid"]
			paramMap["appid"]=paramMap["subappid"]
			paramMap["appid"]=paramMap["subappid"]
			paramMap["appid"]=paramMap["subappid"]
			paramMap["appid"]=paramMap["subappid"]
			paramMap["appid"]=paramMap["subappid"]
			paramMap["appid"]=paramMap["subappid"]
		}
	}

	fmt.Println(paramMap)
}