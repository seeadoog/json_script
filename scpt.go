package jsonscpt

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"
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
		"neq":neq,
		"addp":addp,
		"rem":rem,
		"orr":orr,
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

		e,err:=parseSetExp2(exp)
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

type Script struct {
	Exp
}

func (s *Script)UnmarshalJSON(data []byte)error{
	exp,err:=CompileExpFromJson(data)
	if err != nil{
		return err
	}
	s.Exp = exp
	return nil
}

func (s *Script)Exec(ctx *Context)(error){
	if s.Exp != nil{
		return s.Exp.Exec(ctx)
	}
	return nil
}

//use to replace key of $
const ReplaceKeyDollor ="jsp_replace_key_dollor_675gb__"
const ReplaceKeyReturn ="jsp_replace_key_return_675gb__"
const ReplaceKeyBreak ="jsp_replace_key_break_675gb__"
func preProcessExp(exp string)string{
	exp = strings.Replace(exp,"'","\"",-1)
	exp = strings.Replace(exp,"$",ReplaceKeyDollor,-1)
	exp = strings.Replace(exp,"return",ReplaceKeyReturn,-1)
	exp = strings.Replace(exp,"break",ReplaceKeyBreak,-1)
	return exp
}

func resumeProcessExp(exp string)string{
	exp = strings.Replace(exp,"\"","'",-1)
	exp = strings.Replace(exp,ReplaceKeyDollor,"$",-1)
	exp = strings.Replace(exp,ReplaceKeyReturn,"return",-1)
	exp = strings.Replace(exp,ReplaceKeyBreak,"break",-1)
	return exp
}

func ParseValue(exp string)(Value,error){
	return parseBoolExpFromAstString(exp)
}

func parseBoolExpFromAstString(exp string)(Value,error){
	exp = preProcessExp(exp)
	expr,err:=parser.ParseExpr(exp)
	if err != nil{
		return nil,err
	}
	e,err:=parseBoolExp2(expr)
	if err != nil{
		return nil,fmt.Errorf("invalid exp,%s,err=%s",exp,err.Error())
	}
	return e,nil
}

func parseParamsFromAst(exp *ast.BinaryExpr)([]Value,error){
	parms:=make([]Value,2)
	x,err:=parseBoolExp2(exp.X)
	if err != nil{
		return nil,err
	}
	parms[0] = x
	y,err:=parseBoolExp2(exp.Y)
	if err != nil{
		return nil,err
	}
	parms[1] = y
	return parms, nil
}

func parseFuncByName(exp *ast.BinaryExpr,fun string)(Value,error){
	v:= &FuncValue{
		FuncName: fun,
		Params:   nil,
	}
	ps,err:=parseParamsFromAst(exp)
	if err != nil{
		return nil, err
	}
	v.Params =ps
	return v,nil
}
var m = []int{
	1:1,
	3:4,
	10:10,
}

func parseBoolExp2(e ast.Expr)(Value,error){
	if exp,ok:=e.(*ast.BinaryExpr);ok{
		switch exp.Op {
		case token.LAND: // &&
			return parseFuncByName(exp,"and")
		case token.EQL: // ==
			return parseFuncByName(exp,"eq")
		case token.GTR: // >
			return parseFuncByName(exp,"gt")
		case token.GEQ://>=
			return parseFuncByName(exp,"ge")
		case token.LEQ://<=
			return parseFuncByName(exp,"le")
		case token.LSS:// <
			return parseFuncByName(exp,"lt")
		case token.LOR: // ||
			return parseFuncByName(exp,"orr")
		case token.NEQ://!=
			return parseFuncByName(exp,"neq")
		case token.ADD://+
			return parseFuncByName(exp,"addp")
		case token.MUL://*
			return parseFuncByName(exp,"mul")
		case token.SUB://-
			return parseFuncByName(exp,"sub")
		case token.QUO:// /
			return parseFuncByName(exp,"div")
		case token.REM:// %
			return parseFuncByName(exp,"rem")


		}

	}
	//
	if id,ok:=e.(*ast.Ident);ok{
		//fmt.Println(id.Name)
		return parseValue(resumeProcessExp(id.Name))
		//return nil, nil
	}

	//fmt.Println(reflect.TypeOf(e).String())
	if fun,ok:=e.(*ast.CallExpr);ok{
		params:=make([]Value,0, len(fun.Args))
		for _, v := range fun.Args {
			vv,err:=parseBoolExp2(v)
			if err != nil{
				return nil,err
			}
			params  =  append(params,vv)
		}
		return &FuncValue{
			FuncName: resumeProcessExp(fun.Fun.(*ast.Ident).Name),
			Params:   params,
		},nil
	}
	if bs,ok:=e.(*ast.BasicLit);ok{
		return parseValue(strings.Replace(bs.Value,"\"","'",-1))
	}
	if p,ok:=e.(*ast.ParenExpr);ok{
		if v,err:=parseBoolExp2(p.X);err ==nil{
			return v,nil
		}else{
			return nil,err
		}
	}

	if s,ok:=e.(*ast.SelectorExpr);ok{

		//fmt.Println("str=",selectToString(s))
		return &VarValue{Key:strings.Replace(selectToString(s),ReplaceKeyDollor,"$",-1)},nil
	}
	if idx,ok:=e.(*ast.IndexExpr);ok{
		//fmt.Println("str=",selectToString(idx))
		return &VarValue{Key:strings.Replace(selectToString(idx),ReplaceKeyDollor,"$",-1)},nil
	}
	if u,ok:=e.(*ast.UnaryExpr);ok{
		if u.Op == token.NOT{
			v,err:=parseBoolExp2(u.X)
			if err != nil{
				return nil,err
			}
			return &FuncValue{
				FuncName: "not",
				Params:   []Value{v},
			},nil
		}
	}
	return nil,fmt.Errorf("invalid exp,type=%s",reflect.TypeOf(e).String())
}

func selectToString(s ast.Expr)(string){
	if sl,ok:=s.(*ast.SelectorExpr);ok{
		return selectToString(sl.X)+"."+sl.Sel.Name
	}
	if id,ok:=s.(*ast.Ident);ok{
		return id.Name
	}

	if idx,ok:=s.(*ast.IndexExpr);ok{
		return fmt.Sprintf("%s[%s]",selectToString(idx.X),idx.Index.(*ast.BasicLit).Value)
	}

	return ""
}