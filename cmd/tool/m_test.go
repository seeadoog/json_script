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
	_,err:=vm.Compile("",`
	data=[0,1,2,3,4,5,6,7,8,9]
function bsearsh(arr,n){
	lo =0
	hi=arr.length
	for(;lo<=hi;){
		mid=(lo+hi)/2
		if (arr[mid]===n){
			return mid
		}
		if (arr[mid]<n){
			lo=mid
		}else{
			hi=mid+1
		}
	}
	return -1
}
bsearsh(data,8)
`)
	st:=time.Now()
	for i:=0;i<10000;i++{
		vm=otto.New()
		//vm.Run(cp)
	}
	fmt.Println(err,time.Since(st))
}

// a= 5
// b = 6
// a[0]  = 5
//a [1] = 5
//m]

func TestJson(t *testing.T){
	b:=[]byte(`
[
  {
    "data":[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,18,19,20,21,22,23],
    "key":"arr"
  }
  ,
  {
    "func":"bsearch(arr,n)",
    "do":[
      "lo=0",
      "hi=add(len(arr),-1)",
      {
        "for":"le(lo,hi)",
        "do":[
          "mid=div(add(lo,hi),2,-1)",
          {
            "if":"eq(get(arr,mid),n)",
            "then":[
              "return(mid)"
            ]
          },
          {
            "if":"lt(get(arr,mid),n)",
            "then":"lo=add(mid,1)",
            "else":"hi=mid"
          }
        ]
      },
      "return(-1)"
    ]
  },
  "bsearch(arr,0)"

]
`)
	b=[]byte(`
[
	"$.ent='s'",
	"$.ent='s'",
	"$.ent='s'",
	"$.ent='s'"
]

`)
	cp,err:=jsonscpt.CompileExpFromJson(b)
	st:=time.Now()
	vm:=jsonscpt.NewVm()
	vm.Set("$",map[string]interface{}{
		"ent":"sms16k",
	})
	var obj =  map[string]interface{}{
		"name":"lixiang",
	}
	var a interface{}
	for i:=0;i<5000;i++{
		jsonscpt.MarshalInterface("$.name",obj,"name")
		//if na,ok:=obj["$"].(map[string]interface{});ok{
		//	na["name"] = "name"
		//}else{
		//	obj["$"] = map[string]interface{}{}
		//}
		//	vm:=jsonscpt.NewVm()
		//	vm.SafeExecute(cp,nil)

		//jsonscpt.CachedJsonpathLookUp(obj,"name")
		//a = obj["name"]
	}
	vm.SafeExecute(cp,nil)
	fmt.Println(err,time.Since(st))
	fmt.Println(vm.Get("ast"))
	fmt.Println(obj)
	fmt.Println(a)
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

func TestEx(t *testing.T){
	for i:=0;i<20000000;i++{
		bser([]int{0,1,2,3,4,5,6,7,8,9},8)
	}
}

func bser(a []int,n int)int{
	lo:=0
	hi:= len(a)
	for lo<=hi{
		mid:=(lo+hi)/2
		if a[mid]==n{
			return mid
		}
		if a[mid]<n{
			lo=mid
		}else{
			hi=mid+1
		}

	}
	return -1
}


