package jsonscpt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
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
		"join":1,
		"get":1,
		"set":1,
		"exit":1,
		"trim":1,
		"compare": 1,
		"div":1,
		"mul":1,
	}

	funcs = map[string]Func{   // read only
		"append":apd,
		"len":lens,
		"split":split,
		"printf":printf,
		"print":printlnn,
		"sprintf":sprintf,
		"add":add,
		"delete": deleteFun,
		"isnil":isNil,
		"and":and,
		"or": or,
		"eq": eq,
		"gt": gt,
		"ge": ge,
		"le":le,
		"lt": lt,
		"not": not,
		"return":ret,
		"exit": exit,
		"in": in,
		"contains": contains,
		"join": join,
		"set": set,
		"get": get,
		"input": input,
		"trim": trim,
		"compare": compare,
		"div":div,
		"mul":mul,
		"new":news,
		"from_json":fromJson,
		"to_json":jsonMarshal,
		"number":numberFunc,
		"string":stringFunc,
		"bool":boolFunc,
		"int":intt,
		"sub":sub,
		"throw":throw,
	}
)
//可以用于初始化vm，减少vm创建开销。初始化执行
func SetInlineFunc (name string,f Func){
	if isSystemId(name){
		return
	}
	funcs[name] = f
}

func isSystemId(s string)bool  {
	if systemId[s]==1{
		return true
	}
	return false
}

type Context struct {
	table map[string]interface{}  // save all variables
	funcs map[string]Func
}

func NewVm() *Context {
	c:=&Context{
		table:map[string]interface{}{},
		funcs:funcs,    //内置函数分离，减少创建vm的开销
	}
	//c.init()
	return c
}
type ConcurrencyContext struct {
	*Context
	lock *sync.RWMutex
}

func NewConcurrencyVm() *ConcurrencyContext {
	c:=&ConcurrencyContext{
		Context:NewVm(),    //内置函数分离，减少创建vm的开销
		lock:&sync.RWMutex{},
	}
	//c.init()
	return c
}
func (ctx *ConcurrencyContext)Set(k string,v interface{})  {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()
	MarshalInterface(k,ctx.table,v)
}

func (ctx *ConcurrencyContext)Get(k string)interface{}  {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()
	v,_:=CachedJsonpathLookUp(ctx.table,k)
	return v
}


func (ctx *Context)Func(name string,params ...interface{})  {

}
//不能写到funcmap中，否则会
func (ctx *Context)SetFunc(name string,value Func)  {
	ctx.table[name] = value
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

func (ctx *Context)ExecFile(f string)error{

	b,err:=ioutil.ReadFile(f)
	if err!=nil{
		return err
	}
	return ctx.ExecJson(b)
}

func (c *Context)Execute(exp Exp) error {
	return exp.Exec(c)
}

func (c *Context)SafeExecute(exp Exp, fatalHandler func(err interface{})) error {
	defer func() {
		if err := recover();err !=nil{
			if fatalHandler !=nil{
				fatalHandler(err)
			}
		}
	}()
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
				//return  nil,errors.New("line:"+ifexp+" has no then block")
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

			if forrangeReg.MatchString(forexp){ // for range
				r:=forrangeReg.FindAllStringSubmatch(forexp,-1)
				if len(r)>0{
					if len(r[0])>3{
						key:=r[0][1]
						val:=r[0][2]
						if isSystemId(key){
							return nil,errors.New("system id cannot be variable:"+key)
						}
						if isSystemId(val){
							return nil,errors.New("system id cannot be variable:"+val)
						}
						forVal:=r[0][3]
						parsedForValue,err:=parseValue(forVal)
						if err !=nil{
							return nil,err
						}
						do,err:=CompileExpFromJsonObject(m["do"])
						if err !=nil{
							return nil,err
						}
						r:=&ForRangeExp{
							Value:parsedForValue,
							Do:do,
							SubIdx:key,
							SubValue:val,
						}
						return r,nil
					}
				}
				return nil,errors.New("invalid for range exp"+ String(v))

			}else{  // for bool
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


		}else if data,ok:=m["data"];ok{
			return &DataExp{
				Key:  String(m["key"]),
				Data: data,
			},nil
		}else if gofun,ok:=m["go"];ok{

			exp,err:=CompileExpFromJsonObject(gofun)
			if err !=nil{
				return nil,err
			}
			return &GoFunc{Exp:exp},nil
		}else if fun,ok:=m["func"].(string);ok{
			if params,ok:=m["do"];ok{
				return parseFunc(fun,params)
			}
		}else if try,ok:=m["try"];ok{
			excp,ok:=m["ecp"].(string)
			if !ok{
				return nil, errors.New("invalid ecp of try block")
			}
			if m["do"] ==nil{
				return nil,errors.New("invalid try exp, ecp or do block is nil")
			}
			tryExp,err:=CompileExpFromJsonObject(try)
			if err !=nil{
				return nil, err
			}
			doExp,err:=CompileExpFromJsonObject(m["do"])
			if err !=nil{
				return nil, err
			}
			return &TryCatchExpt{Try:tryExp,Do:doExp,Execption:excp}, nil
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

	return nil,errors.New("invalid exp:"+ String(v))
}
var funcReg = regexp.MustCompile(`(\w+)\((.*)\)`)
func parseFunc(s string,body interface{})(Exp,error){
	if !funcReg.MatchString(s){
		return nil, errors.New("invalid func define:" + s)
	}
	r:=funcReg.FindAllStringSubmatch(s,-1)
	if len(r)>0{
		v:=r[0]
		if len(v)>=2{
			fun:=v[1]
			if isSystemId(fun){
				return nil, errors.New("variable cannot be system Id"+fun)
			}
			do,err:=CompileExpFromJsonObject(body)
			if err !=nil{
				return nil, err
			}
			var params []string
			if len(v[2])>0{
				params = strings.Split(v[2],",")
				for i:=0;i< len(params);i++{
					params[i] = strings.Trim(params[i]," ")
					if isSystemId(params[i]){
						return nil, errors.New("variable cannot be system Id"+params[i])
					}
				}
			}
			fd:=&FuncDefine{
				FuncName:fun,
				Params:params,
				Body:do,
			}
			return fd, nil

		}
	}
	return nil, nil
}

