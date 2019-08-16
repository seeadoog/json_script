package jsonscpt

import (
	"errors"
	"fmt"
	"encoding/json"
)

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
	ctx.SetFunc("sprintf",sprintf)
	ctx.SetFunc("add",add)
	ctx.SetFunc("json_m", jsonMarshal)
	ctx.SetFunc("delete", deleteFun)
	ctx.SetFunc("isnil", isNil)
	ctx.SetFunc("and", and)
	ctx.SetFunc("eq", eq)
	ctx.SetFunc("or", or)
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

func (ctx *Context)Exec(v interface{}) error {
	//if stv,ok:=v.(string);ok{
	//	json.Unmarshal([]byte(stv),&v)
	//}
	cpd,err:=ComplieExp(v)
	if err !=nil{
		return err
	}
	return ctx.CompliedExec(cpd)
}
func (ctx *Context)ExecJson(s []byte) error {
	var i interface{}
	err:=json.Unmarshal(s,&i)
	if err !=nil{
		return err
	}
	return ctx.Exec(i)
}

func (ctx *Context)CompliedExec(v interface{})error {
	switch v.(type) {
	case *IfExp:
		ifexp:=v.(*IfExp)
		if ifexp.If!=nil{
			if ifexp.If.Match(ctx){
				if err:=ctx.CompliedExec(ifexp.Then);err !=nil{
					return err
				}
			}else{
				if ifexp.Else==nil{
					return nil
				}
				if err:=ctx.CompliedExec(ifexp.Else);err !=nil{
					return err
				}
			}
		}
	case *SetExps:
		setexp:=v.(*SetExps)
		err:=setexp.Exec(ctx)
		//fmt.Println(setexp.Variable,setexp.Value.Get(ctx))
		if err !=nil{
			return err
		}
	case []interface{}:
		in:=v.([]interface{})
		for _, vv := range in {
			err:=ctx.CompliedExec(vv)
			if err !=nil{
				return err
			}
		}
	default:
		return errors.New("invalid complied type:"+fmt.Sprintf("%v",v))
	}
	return nil
}


func boolValid(s string) bool {
	return true
}


type IfExp struct {
	If *BoolExp
	Then interface{}
	Else interface{}
}



func ComplieExp(v interface{}) (interface{},error) {
	if exp,ok:=v.(string);ok{
		e,err:=parseSetExp(exp)
		//fmt.Println("-----------------")
		if err !=nil{
			return nil, err
		}
		return e,nil
	}
	if m,ok:=v.(map[string]interface{});ok{
		exp:=&IfExp{}
		if ifexp,ok:=m["if"].(string);ok{
			efe,err:=parseBoolExp(ifexp);
			if err==nil{
				exp.If = efe
			}else{
				return nil,err
			}
		}
		if then,ok:=m["then"];ok && then !=nil {
			var parsedExp,err = ComplieExp(then)
			if err !=nil{
				return nil,err
			}
			exp.Then = parsedExp
		}
		if el,ok:=m["else"];ok && el !=nil{
			var parsedExp,err = ComplieExp(el)
			if err !=nil{
				return nil,err
			}
			exp.Else = parsedExp
		}
		return exp,nil
	}

	if eps,ok:=v.([]interface{});ok {
		var parsedExp = make([]interface{},0, len(eps))
		for i := 0; i < len(eps); i++ {
			e,err:= ComplieExp(eps[i])
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


