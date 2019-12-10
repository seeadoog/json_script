package jsonscpt

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCachedJsonpathLookUp(t *testing.T) {
	var m = map[string]interface{}{
		"name": "lxiang",
		"son": map[string]interface{}{
			"name": "bs",
		},
	}
	MarshalInterface("lock.name", m, "hee")
	fmt.Println(CachedJsonpathLookUp(m, "son.name"))
	fmt.Println(CachedJsonpathLookUp(m, "lock.name"))

}

func TestReg(t *testing.T) {
	fmt.Println(funv.MatchString("len(s,b,c)"))
	r := funv.FindAllSubmatch([]byte("len(s,b,c)"), -1)
	fmt.Println(string(r[0][1]))
	fmt.Println(string(r[0][2]))
	//for _, v := range r {
	//	//for _, vv := range v {
	//	//	fmt.Println(string(vv))
	//	//}
	//	fmt.Println
	//}
}

func TestReg2(t *testing.T) {
	v, err := parseValue("append('hello ','world',6,len('mask'),name,' ',common.ent)")
	if err != nil {
		panic(err)
	}
	set, err := parseSetExp("common.ent=append('sms16k_en','_',rate)")
	if err != nil {
		panic(err)
	}

	set2, err := parseSetExp("common.names=split(common.ent,'_',2)")
	if err != nil {
		panic(err)
	}
	vm := NewVm()
	vm.Set("len", lens)
	vm.Set("append", apd)
	vm.Set("name", " zhaobiao")
	vm.Set("rate", "16000")
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

func TestParse(t *testing.T) {
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
	json.Unmarshal([]byte(s), &m)
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
	json.Unmarshal([]byte(s2), &m)
	_, err := CompileExpFromJsonObject(m)
	if err != nil {
		panic(err)
	}
	busi := map[string]interface{}{
		"ent":  "sms",
		"rate": 16000,
	}
	for i := 0; i < 10000; i++ {
		vm := NewVm()
		vm.Set("business", busi)
		CompileExpFromJsonObject(m)
	}

	fmt.Println(busi)
	for i := 0; i < 1; i++ {
		v2 := NewVm()
		err = v2.ExecJsonObject(`name""=5`)
		fmt.Println("name=", err, v2.Get("name"))
	}

	//fmt.Println(vm.Get("common"))
}

func TestExec(t *testing.T) {
	vm := NewVm()
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

//== >= <= | && | ||缺少

func TestBenchMark(t *testing.T) {

	var param = map[string]interface{}{
		"name": "lixiang",
		"age":  15,
		"desc": "",
	}

	for i := 0; i < 1; i++ {
		if param["name"].(string) == "lixiang" && param["age"].(int) > 10 {
			param["type"] = "old"
			param["desc1"] = fmt.Sprintf("name=%v,age=%v", param["name"], param["age"])
		} else {
			param["desc1"] = "a child"
		}

		param["desc"] = param["desc"].(string) + "js"
	}
	fmt.Println(param["desc"])
}

func TestBenchMark2(t *testing.T) {

	var param = map[string]interface{}{
		"name": "lixiang",
		"age":  1,
		"desc": "",
		"dsf": map[string]interface{}{
			"name": "lixiang",
			"age":  1,
		},
		"data": map[string]interface{}{
			"audio_src":   "s3",
			"s3_access":   "34",
			"s3_endpoint": "34",
			"s3_key":      "34",
			"s3_bucket":   "",
		},
	}

	vm := NewVm()
	vm.Set("$", param)
	SetInlineFunc("validate", func(i ...interface{}) interface{} {
		if len(i) > 1 {
			if m, ok := i[0].(map[string]interface{}); ok {
				if _, ok := m[String(i[1])]; ok {
					return &ErrorExit{Code: 10106, Message: fmt.Sprintf("parma %v is required", String(i[1]))}
				}
			}
		}
		return nil
	})

	exp := []byte(`
[
         
          {
            "if": "eq($.data.audio_src,'s3')",
            "then": [
              "validate($.data,'s3_access')",
              "validate($.data,'s3_endpoint')",
              "validate($.data,'s3_key')",
              "validate($.data,'s3_bucket')"
            ]
          }
        ]

`)
	cmd, err := CompileExpFromJson(exp)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100000; i++ {
		err = vm.Execute(cmd)
	}
	if err, ok := IsExitError(err); ok {
		fmt.Println("err->", err)
	}
	fmt.Println(vm.Get("param.desc"), err)
}

func TestMarshal(t *testing.T) {
	var param = map[string]interface{}{
		"name": "lixiang",
		"age":  1,
		"desc": "",
		"dsf": map[string]interface{}{
			"name": "lixiang",
			"age":  1,
			"desc": "",
		},
	}

	for i := 0; i < 100000; i++ {
		json.Marshal(param)
	}
}

func TestMap(t *testing.T) {
	var s = []string{}
	for k, _ := range systemId {
		s = append(s, k)
	}
	for i := 0; i < 10000000; i++ {
		if isSystemId("in") {

		}
	}
}

func TestMap2(t *testing.T) {
	a := [][]int{}
	sums(10, &a)
	fmt.Println(a)
}

//6
//1+5 2+4 3+3
//
func isContains(ss []string, s string) bool {

	for i := 0; i < len(ss); i++ {
		if ss[i] == s {
			return true
		}
	}
	return false
}

func sums(n int, t *[][]int) {
	sb := []int{}
	for i := 1; i <= n/2; i++ {
		sb = append(sb, i, n-i)
		*t = append(*t, sb)
		sb = sb[:0]
		if i > 1 {
			r := &[][]int{}
			sums(i, r)
			for _, v := range *r {
				*t = append(*t, append(v, n-i))
			}
		}
		if n-i > 1 {
			r := &[][]int{}
			sums(n-i, r)
			for _, v := range *r {
				*t = append(*t, append(v, i))
			}
		}
	}
}

func TestJson(t *testing.T) {
	//f, _ := ioutil.ReadFile(`C:\Users\admin\Documents\WeChat Files\lsj1581503100\FileStorage\File\2019-09\aa.txt`)
	var a  = map[string]interface{}{
		"business":map[string]interface{}{
			"personalization_multiple":map[string]interface{}{
				"key1":"v1",
				"key2":"v1",
			},
		},
	}

	vm := NewVm()
	vm.Set("$", a)
	err:=vm.ExecJson([]byte(`
[
  {
    "if": "not(isnil($.business.personalization_multiple))",
    "then": [
      "res=''",
      {
        "for": "k,v in $.business.personalization_multiple",
        "do": ["print(k,v)","res=append(res,k,'=',v,';')"]
      },
      "print('then',res)"
    ]
  },
  {
	"if":"isnil($.abg)",
	"then":"print('abg is nil')"
  },
	{
	"if":"$.business1",
	"then":"print('yes')"
}
]

`))
	fmt.Println(err)
	fmt.Println(vm.Get("$"))
	b, _ := json.Marshal(vm.Get("$"))
	fmt.Println(string(b))

}

func TestTryCatchExpt_Exec(t *testing.T) {
	js := []byte(`

{
  "common":{
    "app_id":"123456"
  },
  "business":{
    "ent":"intp65",
    "vcn":"xiaoyan",
    "volume":50,
    "speed":50
  },
  "data":{
    "status":2,
    "text":"exSI6ICJlbiIsCgkgICAgInBvc2l0aW9uIjogImZhbHNlIgoJf..."
  },
  "input":[
		[ {"name":"liix","status":2,"isok":true}]
	]
}

`)

	b, err := GenerateSchemaFromJson(js)
	fmt.Println(err)
	fmt.Println(string(b))

}

func gengrateSchema(o interface{}, t map[string]interface{}) {
	switch o.(type) {
	case map[string]interface{}:
		t["type"] = "object"
		p := map[string]interface{}{}
		t["properties"] = p
		for k, v := range o.(map[string]interface{}) {
			p[k] = map[string]interface{}{}
			gengrateSchema(v, p[k].(map[string]interface{}))
		}
	case string:
		t["type"] = "string"
	case float64,float32,int,int64,int32:
		t["type"] = "number"
	case []interface{}:
		t["items"] = map[string]interface{}{}
		t["type"] = "array"
		is := o.([]interface{})
		if len(is) > 0 {
			gengrateSchema(is[0], t["items"].(map[string]interface{}))
		}else{

		}
	case bool:
		t["type"] = "boolean"
	}
}

func GenerateSchemaFromJson(in []byte)([]byte,error){
	var i interface{}
	if err:=json.Unmarshal(in,&i);err != nil{
		return nil,err
	}
	sc:=make(map[string]interface{})
	gengrateSchema(i,sc)
	return json.Marshal(sc)
}

