package jsonscpt

import (
	"testing"
	"fmt"
	"encoding/json"
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
	//fmt.Println(v.(*FuncValue).params)
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
//p,_:= ComplieExp(m)
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
	_,err:=ComplieExp(m)
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
		ComplieExp(m)
	}

	fmt.Println(busi)
	for i:=0;i<1;i++{
		v2:=NewVm()
		err=v2.Exec(`name""=5`)
		fmt.Println("name=",err,v2.Get("name"))
	}

	//fmt.Println(vm.Get("common"))
}

func TestExec(t *testing.T) {
	vm:=NewVm()
	//vm.Exec(`printf('%f %f %f %s',1,2,3,append(len(append('hello','world')),'nima'))`)
	vm.Exec(`user.name='lixiang'`)
	vm.Exec(`user.password='123456'`)
	vm.Exec(`jsonStr=json_m(user)`)
	vm.Exec("printf('%s',jsonStr)")
	vm.Exec("array=split('1,2,3,4,5',',')")
	vm.Exec("printf('%v',array)")
	vm.Exec("printf('%s',array[0])")

	vm.Exec("jsonpd=split(jsonStr,',')")
	fmt.Println(vm.Exec("printf('%s',array[1])"))
	fmt.Println(vm.Exec("printf('%v',jsonpd)"))
	fmt.Println(vm.Exec("printf('len of splited array is %d',len(array))"))
	vm.SetFunc("show", func(i ...interface{}) interface{} {
		fmt.Println(i...)
		return ""
	})
	vm.Exec("show('exec show function',1,2,3,4,5)")
	vm.Exec("show(user)")
	vm.Exec("key='name'")
	vm.Exec("delete(user,key)")
	vm.Exec("show(user)")
	vm.Exec("show('nil value:',hal,isnil(user))")
	vm.Exec("ar.len=len('123')")
	fmt.Println(vm.Get("ar"))

}

//== >= <= | && | ||




func TestScript(t *testing.T) {
	vm:=NewVm()
	var param = map[string]interface{}{
		"channel":"ens",
		"language":"en_us",
		"domain":"iat",
		"sample_rate":"16000",
		"appid":"123456",
	}

	vm.Set("param",param)
	b:=[]byte(`
[
  {
    "if": "not(param.language)",
    "then": "return(10137,'param language is required')"
  },
  {
    "if": "not(param.domain)",
    "then": "return(10137,'param domain is required')"
  },
  {
    "if": "not(param.sample_rate)",
    "then": "return(10137,'param sample_rate is required a valid value is between [16000,16k,8000,8k]')"
  },
  {
    "if": "or(not(param.language),not(param.domain),not(param.sample_rate))",
    "then": [
      "return(10137,'param convert error')"
    ],
    "else": [
      {
        "if": "and(eq(param.language,'en_us'),or(eq(param.sample_rate,'16000'),eq(param.sample_rate,'16k')))",
        "then": "param.ent='mardin_16k'"
      }
    ]
  },
  {
    "if": "and(not(param.dwa),eq(param.appid,'123456'))",
    "then": "param.dwa='wpgs'",
    "else": "param.dwa=''"
  },
  "param.cc=5"
]`)
	var err error
	for j:=0;j<0;j++{
		err=vm.ExecJson(b)
	}

	var i interface{}
	json.Unmarshal(b,&i)
	compiledCode,err:=ComplieExp(i)
	if err !=nil{
		panic(err)
	}
	for j:=0;j<100000;j++{
		//ComplieExp(i)
		err = vm.CompliedExec(compiledCode)
	}
	fmt.Println(err)
	fmt.Println(param)
	fmt.Println(vm.Get("tmp"))
}