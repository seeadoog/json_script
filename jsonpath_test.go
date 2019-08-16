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
  "else":["common.appid='100IME'","business.auf='rate;16000'",{
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
  "then":["common.appid='100IME'","business.auf='rate;16000'",{
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

func TestTrim(t *testing.T) {
	//fmt.Println(trimBlock("(a=b)(c==d)"))
	//fmt.Println(trimBlock("(()())"))
	fmt.Println(TrimBlock("((a=b)(c==d))"))
	fmt.Println(TrimBlock("(a==b)"))
}


//== >= <= | && | ||
func TestCompiled_String(t *testing.T) {
	b:=[]byte(`name=='lixiang'`)
	//b=[]byte(`((3==3)&& (2==2)) && (3==3)`)
	//b=[]byte(`isnil(name,age)`)
	_,err:=parseBoolV(string(b))
	if err !=nil{
		panic(err)
	}
	vm:=NewVm()
	vm.Exec("name='lixiang1'")
	var v interface{}
	err=json.Unmarshal([]byte(`[{
	"if":"name=='lixiang'",
	"then":"printf('%s',name)",
	"else":"printf('%s','nima')"
}]`),&v)
	fmt.Println(vm.Exec(v),err)
	vm.Exec("printf('%f',1)")

	//fmt.Println(v.Get(vm))
	//for _, v := range r {
	//	for k, vv := range v {
	//		fmt.Println(k,string(vv))
	//	}
	//}
}

func TestJsonPathLookup(t *testing.T) {
	vm:=NewVm()
	vm.Set("common.app_id","100IME")
	vm.Set("common.uid","123456")
//	vm.Set("business.eos","123456")
	vm.Set("data.text","123456")
	fmt.Println(vm.Exec(`a=b`))
	fmt.Println(vm.Get("common.app_id"))
	err:=vm.ExecJson([]byte(`
[
  {
    "if": "or(and(eq(common.app_id,'100IME1'),eq(2,2)),eq(2,2))",
    "then": [
      {
        "if": "eq(common.app_id,'100IME')",
        "then": "printf('%s is niubi',common.app_id)",
        "else": "printf('invalid uid')"
      },
      {
        "if":"isnil(business.eos)",
        "then":"business.eos=12000"

      }
    ],
    "else": "printf('invalid appid')"
  },
  "business.datelen=len(data.text)",
  "len1=5",
  {
	"if":"true",
	"then":"printf('true is true')"
  }
]
`))
	fmt.Println(err)
	fmt.Println(vm.Get("business"))
}

func TestCompiled_Lookup(t *testing.T) {
	vm:=NewVm()
	err:=vm.ExecJson([]byte(`
[
  	"a='hello'",
	{
		"if":"a=='hello'",
  		"then":"b='is bello'",
		"else":"b='is not hello'"
	},
	"bi.r=b",
	"bi.len=len(b)",
	"bi.len23=len('123')"
]

`))
	fmt.Println(err)
	fmt.Println(vm.Get("a"))
	fmt.Println(vm.Get("b"))
	fmt.Println(vm.Get("bi"))
}

func TestBoolExp_Match(t *testing.T) {
	vm:=NewVm()
	vm.Exec("printf('%v',eq(2,2))")
	b:=[]byte(`
[
  "user.name='lixiang'",
  "user.age=1",
  {
    "if":"gt(user.age,10)",
    "then":"show('%s is an old man',user.name)",
    "else":[
      "show('%s is an child',user.name)",
      "user.isold=false"
    ]
  },
  "show('%v',user)"
]
`)
	var i interface{}
	err:=json.Unmarshal(b,&i)
	if err !=nil{
		panic(err)
	}
	var cd interface{}
	for j:=0;j<1;j++{
		cd,err =ComplieExp(i)
		if err !=nil{
			panic(err)
		}
	}

	for i:=0;i<100000;i++{
		err:=vm.CompliedExec(cd)
		if err !=nil{
			panic(err)
		}
	}

	fmt.Println(err)
}