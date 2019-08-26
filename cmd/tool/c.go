package main

import (
	"fmt"
	"git.xfyun.cn/AIaaS/json_script"
)

func main1(){
	rule:=[]byte(`
[
  {
	"if":"eq('a','a')"
  }
]
`)
	scp,err:=jsonscpt.CompileExpFromJson(rule)
	if err !=nil{
		panic(err)
	}
	user :=map[string]interface{}{
		"name":"bob",
		"age":"16",
	}
	var vm *jsonscpt.Context
	vm =jsonscpt.NewVm()
	for i:=0;i<100000;i++{
		err =vm.SafeExecute(scp,nil)
	}


	fmt.Println("userMap:",user)
}