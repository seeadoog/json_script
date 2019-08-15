package jsonscpt

import "errors"

type Context struct {
	table map[string]interface{}
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
}

func (ctx *Context)Func(name string,params ...interface{})  {

}

func (ctx *Context)SetFunc(name string,value Func)  {
	ctx.Set(name,value)
}

func (ctx *Context)Set(k string,v interface{})  {
	ctx.table[k] = v
}

func (ctx *Context)Get(k string)interface{}  {
	v,_:=CachedJsonpathLookUp(ctx.table,k)
	return v
}
func (ctx *Context)CompliedExec(v interface{})error {
	switch v.(type) {
	case *IfExp:
		ifexp:=v.(*IfExp)
		if ifexp.If!=nil{
			if ifexp.If.Match(){
				if err:=ctx.CompliedExec(ifexp.Then);err !=nil{
					return err
				}
			}else{
				if err:=ctx.CompliedExec(ifexp.Else);err !=nil{
					return err
				}
			}
		}
	case *SetExps:
		setexp:=v.(*SetExps)
		err:=setexp.Exec(ctx)
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
		return errors.New("invalid complied type")
	}
	return nil
}


func boolValid(s string) bool {
	return true
}


type IfExp struct {
	If *BoolExp
	Then []interface{}
	Else []interface{}
}



func ComplieExp(v interface{}) (interface{},error) {
	if exp,ok:=v.(string);ok{
		e,err:=parseSetExp(exp)
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
		if then,ok:=m["then"].([]interface{});ok {
			var parsedExp,err = ComplieExp(then)
			if err !=nil{
				return nil,err
			}
			exp.Then = parsedExp.([]interface{})
		}
		if el,ok:=m["else"].([]interface{});ok {
			var parsedExp,err = ComplieExp(el)
			if err !=nil{
				return nil,err
			}
			exp.Else = parsedExp.([]interface{})
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


