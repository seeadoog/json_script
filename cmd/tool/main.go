package main

import (
	"git.xfyun.cn/AIaaS/json_script"
	"flag"
	"io/ioutil"
	"fmt"
	"time"
)
var file = flag.String("f","input.json","")
func main(){
	flag.Parse()
	vm:=jsonscpt.NewVm()
	b,err:=ioutil.ReadFile(*file)
	if err !=nil{
		fmt.Println(err.Error())
		return
	}
	cmd,err:=jsonscpt.CompileExpFromJson(b)
	if err !=nil{
		fmt.Println(err)
		return
	}
	start:=time.Now()
	if err:=vm.Execute(cmd);err !=nil{
		fmt.Println(err)
	}
	fmt.Println("total cost=>",time.Since(start))

}

