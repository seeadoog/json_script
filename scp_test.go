package jsonscpt

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCachedJsonpathLookUp(t *testing.T) {
	var m = map[string]interface{}{
		"name":"lxiang",
		"son":map[string]interface{}{
			"name":"bs",
		},
	}
	MarshalInterface("lock.name",m,"hee")
	fmt.Println(CachedJsonpathLookUp(m,"son.name"))
	fmt.Println(CachedJsonpathLookUp(m,"lock.name"))

}

func TestReg(t *testing.T)  {
	fmt.Println(funv.MatchString("len(s,b,c)"))
	r:=funv.FindAllSubmatch([]byte("len(s,b,c)"),-1)
	fmt.Println(string(r[0][1]))
	fmt.Println(string(r[0][2]))
	//for _, v := range r {
	//	//for _, vv := range v {
	//	//	fmt.Println(string(vv))
	//	//}
	//	fmt.Println
	//}
}

func TestReg2(t *testing.T)  {
	v,err:=parseValue("append('hello ','world',6,len('mask'),name,' ',common.ent)")
	if err !=nil{
		panic(err)
	}
	set,err:=parseSetExp("common.ent=append('sms16k_en','_',rate)")
	if err !=nil{
		panic(err)
	}

	set2,err:=parseSetExp("common.names=split(common.ent,'_',2)")
	if err !=nil{
		panic(err)
	}
	vm:=NewVm()
	vm.Set("len",lens)
	vm.Set("append",apd)
	vm.Set("name"," zhaobiao")
	vm.Set("rate","16000")
	set.Exec(vm)
	set2.Exec(vm)
	fmt.Println(v.Get(vm))
	fmt.Println(vm.Get("common.ent"))
	fmt.Println(vm.Get("common.names[1]"))
	//fmt.Println(v.(*FuncValue).FuncName)
	//fmt.Println(v.(*FuncValue).Params)
	//for _, v := range r {
	//	//for _, vv := range v {
	//	//	fmt.Println(string(vv))
	//	//}
	//	fmt.Println
	//}
}

func TestParse(t*testing.T){
	var s = `
{
  "if":"a > b",
  "then":[{
    "if":""
  }],
  "else":["common.appid='123456'","business.auf='rate;16000'",{
    "if":"c > 6 & b == 5",
    "then":["common.appid1=''"]
  }]
}
`
var m = map[string]interface{}{}
json.Unmarshal([]byte(s),&m)
//p,_:= CompileExpFromJsonObject(m)
//fmt.Println(p.(*IfExp).Else[0].(*SetExps).Variable)
//fmt.Println(p.(*IfExp).Else[1].(*SetExps).Variable)
//fmt.Println(p.(*IfExp).Else[2].(*IfExp).Then[0].(*SetExps).Variable)
}

func TestCompile(t *testing.T) {
	var s = `
{
  "if":"a > b",
  "else":[{
    "if":""
  }],
  "then":["common.appid='123456'","business.auf='rate;16000'",{
    "if":"c > 6 & b == 5",
    "then":["common.appid1='hello'"]
  }]
}
`
fmt.Println(s)
var s2 = `[
{
	"if":"business.rate==16000",
	"then":["business.rate='16k'","business.name='lixiang'"]
},
"business.finalEnt=append(business.ent,'_',business.rate)"
]`
	var m interface{}
	json.Unmarshal([]byte(s2),&m)
	_,err:= CompileExpFromJsonObject(m)
	if err !=nil{
		panic(err)
	}
	busi:=map[string]interface{}{
		"ent":"sms",
		"rate":16000,
	}
	for i:=0;i<10000;i++{
		vm:=NewVm()
		vm.Set("business",busi)
		CompileExpFromJsonObject(m)
	}

	fmt.Println(busi)
	for i:=0;i<1;i++{
		v2:=NewVm()
		err=v2.ExecJsonObject(`name""=5`)
		fmt.Println("name=",err,v2.Get("name"))
	}

	//fmt.Println(vm.Get("common"))
}

func TestExec(t *testing.T) {
	vm:=NewVm()
	//vm.ExecJsonObject(`printf('%f %f %f %s',1,2,3,append(len(append('hello','world')),'nima'))`)
	vm.ExecJsonObject(`user.name='lixiang'`)
	vm.ExecJsonObject(`user.password='123456'`)
	vm.ExecJsonObject(`jsonStr=json_m(user)`)
	vm.ExecJsonObject("printf('%s',jsonStr)")
	vm.ExecJsonObject("array=split('1,2,3,4,5',',')")
	vm.ExecJsonObject("printf('%v',array)")
	vm.ExecJsonObject("printf('%s',array[0])")

	vm.ExecJsonObject("jsonpd=split(jsonStr,',')")
	fmt.Println(vm.ExecJsonObject("printf('%s',array[1])"))
	fmt.Println(vm.ExecJsonObject("printf('%v',jsonpd)"))
	fmt.Println(vm.ExecJsonObject("printf('len of splited array is %d',len(array))"))
	vm.SetFunc("show", func(i ...interface{}) interface{} {
		fmt.Println(i...)
		return ""
	})
	vm.ExecJsonObject("show('exec show function',1,2,3,4,5)")
	vm.ExecJsonObject("show(user)")
	vm.ExecJsonObject("key='name'")
	vm.ExecJsonObject("delete(user,key)")
	vm.ExecJsonObject("show(user)")
	vm.ExecJsonObject("show('nil value:',hal,isnil(user))")
	vm.ExecJsonObject("ar.len=len('123')")
	fmt.Println(vm.Get("ar"))

}

//== >= <= | && | ||






func TestBenchMark(t *testing.T)  {

	var param = map[string]interface{}{
		"name":"lixiang",
		"age":15,
		"desc":"",
	}

	for i:=0;i<1;i++{
		if param["name"].(string)=="lixiang" && param["age"].(int)>10{
			param["type"] = "old"
			param["desc1"] = fmt.Sprintf("name=%v,age=%v",param["name"],param["age"])
		}else{
			param["desc1"] = "a child"
		}

		param["desc"] = param["desc"].(string)+"js"
	}
	fmt.Println(param["desc"])
}

func TestBenchMark2(t *testing.T)  {

	var param = map[string]interface{}{
		"name":"lixiang",
		"age":1,
		"desc":"",
		"dsf":map[string]interface{}{
			"name":"lixiang",
			"age":1,
			"desc":"",

		},
	}

	vm:=NewVm()
	vm.Set("param",param)

	exp:=[]byte(`
[{
	"if":"and(eq(param.name,'lixiang'),gt(param.age,10))",
	"then":["param.type='old'","param.desc1=sprintf('name=%v,age=%v',param.name,param.age)"],
	"else":"param.desc1='a child'"
},

{
	"if":"in(param.name,'lixiang','zhaobiao')",
	"then":"printf('%v','yes niubi')"
}
]

`)
	cmd,err:= CompileExpFromJson(exp)
	if err !=nil{
		panic(err)
	}
	for i:=0;i<1;i++{
		vm.Execute(cmd)
	}
	fmt.Println(vm.Get("param.desc"))
}

func TestMarshal(t *testing.T) {
	var param = map[string]interface{}{
		"name":"lixiang",
		"age":1,
		"desc":"",
		"dsf":map[string]interface{}{
			"name":"lixiang",
			"age":1,
			"desc":"",
		},
	}

	for i:=0;i<100000;i++{
		json.Marshal(param)
	}
}

