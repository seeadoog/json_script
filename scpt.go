package jsonscpt

import (
	"encoding/json"
	"errors"
	"fmt"
)
var (
	systemId = map[string]int{
		"append":1,
		"len":1,
		"split":1,
		"printf":1,
		"print":1,
		"sprintf":1,
		"add":1,
		"json_m":1,
		"isnil":1,
		"delete":1,
		"and":1,
		"eq":1,
		"or":1,
		"true":1,
		"false":1,
		"gt":1,
		"ge":1,
		"lt":1,
		"le":1,
		"not":1,
		"return":1,
		"in":1,
		"contains":1,
		"index":1,
	}
)

func isSystemId(s string)bool  {
	if systemId[s]==1{
		return true
	}
	return false
}

type Context struct {
	table map[string]interface{}  // save all variables
}

func NewVm() *Context {
	c:=&Context{
		table:map[string]interface{}{},
	}
	c.init()
	return c
}


func (ctx *Context)init()  {
	ctx.SetFunc("append",apd)
	ctx.SetFunc("len",lens)
	ctx.SetFunc("split",split)
	ctx.SetFunc("printf",printf)
	ctx.SetFunc("print",printlnn)
	ctx.SetFunc("sprintf",sprintf)
	ctx.SetFunc("add",add)
	ctx.SetFunc("json_m", jsonMarshal)
	ctx.SetFunc("delete", deleteFun)
	ctx.SetFunc("isnil", isNil)
	ctx.SetFunc("and", and)
	ctx.SetFunc("or", or)
	ctx.SetFunc("eq", eq)
	ctx.SetFunc("gt", gt)
	ctx.SetFunc("ge", ge)
	ctx.SetFunc("le", le)
	ctx.SetFunc("lt", lt)
	ctx.SetFunc("not", not)
	ctx.SetFunc("return", ret)
	ctx.SetFunc("in", in)
	ctx.SetFunc("contains", contains)
	ctx.SetFunc("index", index)


}

func (ctx *Context)Func(name string,params ...interface{})  {

}

func (ctx *Context)SetFunc(name string,value Func)  {
	ctx.Set(name,value)
}

func (ctx *Context)Set(k string,v interface{})  {
	MarshalInterface(k,ctx.table,v)
}

func (ctx *Context)Get(k string)interface{}  {
	v,_:=CachedJsonpathLookUp(ctx.table,k)
	return v
}

func (ctx *Context) ExecJsonObject(v interface{}) error {
	//if stv,ok:=v.(string);ok{
	//	json.Unmarshal([]byte(stv),&v)
	//}
	cpd,err:= CompileExpFromJsonObject(v)
	if err !=nil{
		return err
	}
	return ctx.Execute(cpd)
}
func (ctx *Context)ExecJson(s []byte) error {
	var i interface{}
	err:=json.Unmarshal(s,&i)
	if err !=nil{
		return err
	}
	return ctx.ExecJsonObject(i)
}

func (c *Context)Execute(exp Exp) error {
	return exp.Exec(c)
}

func CompileExpFromJson(b []byte)(Exp,error){
	var i interface{}
	if err:=json.Unmarshal(b,&i);err !=nil{
		return nil,err
	}
	return CompileExpFromJsonObject(i)
}

func boolValid(s string) bool {
	return true
}


type IfExp struct {
	If *BoolValue
	Then Exp
	Else Exp
}

func (f *IfExp)Exec(ctx *Context)error{
	if f.If.Match(ctx){
		err:=f.Then.Exec(ctx)
		if err !=nil{
			return err
		}
	}else{
		if f.Else!=nil{
			err:=f.Else.Exec(ctx)
			if err !=nil{
				return err
			}
		}
	}
	return nil
}




func CompileExpFromJsonObject(v interface{}) (Exp,error) {
	if exp,ok:=v.(string);ok{
		if exp=="break"{   // parse break
			return &BreakExp{},nil
		}

		e,err:=parseSetExp(exp)
		if err !=nil{
			return nil, err
		}
		return e,nil
	}
	if m,ok:=v.(map[string]interface{});ok{

		if ifexp,ok:=m["if"].(string);ok{
			exp:=&IfExp{}
			efe,err:=parseBoolExp(ifexp);
			if err==nil{
				exp.If = efe
			}else{
				return nil,err
			}
			if then,ok:=m["then"];ok && then !=nil {
				var parsedExp,err = CompileExpFromJsonObject(then)
				if err !=nil{
					return nil,err
				}
				exp.Then = parsedExp
			}else{
				return  nil,errors.New("line:"+ifexp+" has no then block")
			}
			if el,ok:=m["else"];ok && el !=nil{
				var parsedExp,err = CompileExpFromJsonObject(el)
				if err !=nil{
					return nil,err
				}
				exp.Else = parsedExp
			}
			return exp,nil
		}else if forexp,ok:=m["for"].(string);ok{  //parse for
			exp:=&ForExp{}
			efe,err:=parseBoolExp(forexp);
			if err !=nil{
				return nil,err
			}
			exp.Addtion = efe
			if do,ok:=m["do"];ok && do !=nil{
				blexp,err:=CompileExpFromJsonObject(do)
				if err !=nil{
					return nil,err
				}
				exp.Do = blexp
			}else{
				return nil,errors.New("for do is nil")
			}
			return exp,nil
		}

		return nil,errors.New("invalid object:"+fmt.Sprintf("%v",v))
	}

	if eps,ok:=v.([]interface{});ok {
		var parsedExp = Exps{}
		for i := 0; i < len(eps); i++ {
			e,err:= CompileExpFromJsonObject(eps[i])
			if err !=nil{
				return nil,err
			}else{
				parsedExp = append(parsedExp,e)
			}
		}
		return parsedExp,nil
	}

	return nil,nil
}


