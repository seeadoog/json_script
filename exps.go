package jsonscpt

import (
	"errors"
	"unsafe"
	"strings"
)

type Exp interface  {
	Exec(ctx *Context)error
}
//add ( 1 , len ( 'd' ) )

func trimspace(s string)string  {
	return strings.Trim(s," ")
}
func parseSetExp(s string)(Exp, error){
	v,err:=splitSetExp(s)
	if err !=nil{
		return nil,err
	}
	if len(v)!=2{
		if len(v)==1{
			val,err:=parseValue(trimspace(v[0]))
			if err !=nil{
				return nil,err
			}
			return &FuncExp{Value:val},nil
		}
		return nil,errors.New("invalid set exp:"+s)
	}
	v[0] =trimspace(v[0])
	v[1] =trimspace(v[1])
	if !checkRule(v[0]){
		return nil,errors.New("invalid setexp key:"+v[0])
	}
	if isSystemId(v[0]){
		return nil,errors.New("variable name cannot be system identify :"+s)
	}
	val,err:=parseValue(v[1])
	if err !=nil{
		return nil,err
	}
	return &SetExps{Variable:v[0],Value:val},nil
}

func splitSetExp(s string)([]string,error)  {
	var tokens = []string{}
	var token = make([]byte,0, len(s))
	cs:=0
	for i:=0;i< len(s);i++{
		v:=s[i]
		token = append(token,v)
		if v== '\'' && cs == 0{
			cs = 1
			continue
		}
		if v=='\'' && cs ==1{
			cs=0
			continue
		}
		if v=='=' && cs == 0 {
			if len(token)<1{
				return nil,errors.New("invalid set token:"+s)
			}
			tokens = append(tokens,string(token[:len(token)-1]))
			token=token[:0]
		}
	}
	if cs==1{
		return nil,errors.New("invalid set exp:"+s)
	}
	if len(token)>0{
		tokens = append(tokens,string(token))
	}
	return tokens,nil
}

func toString(b []byte)string  {
	return *(*string)(unsafe.Pointer(&b))
}

type SetExps struct{
	Variable string
	Value Value
}

func (e *SetExps)Exec(ctx *Context)error{
	if fv,ok:=e.Value.(*FuncValue);ok{
		if fv.FuncName=="return"{
			err:=e.Value.Get(ctx)
			if re,ok:=err.(*ErrorReturn);ok{
				return re
			}
		}
	}
	return MarshalInterface(e.Variable,ctx.table,e.Value.Get(ctx))
}
//a>b && (c>d)
type BoolValue struct {
	//Op string  // > ,== ,< ,>= ,<=
	Value Value
}

func (b *BoolValue)Match(ctx *Context)bool  {
	return convertToBool(b.Value.Get(ctx))
}


//func (b *BoolValue)ExecJsonObject(ctx *Context)error{
//
//}


func convertToBool(v interface{})bool  {
	switch v.(type) {
	case bool:
		return v.(bool)
	case string:
		if len(v.(string))>0{
			return true
		}
		return false
	case float64:
		if int(v.(float64))>0{
			return true
		}
		return false
	}
	return false
}


type IfExp struct {
	If *BoolValue
	Then Exp
	Else Exp
}

func (f *IfExp)Exec(ctx *Context)error{
	if f.If.Match(ctx){
		if f.Then !=nil{
			return f.Then.Exec(ctx)
		}
		return nil
	}else{
		if f.Else!=nil{
			return  f.Else.Exec(ctx)
		}
	}
	return nil
}

// a collection of exps
type Exps []Exp

func (es Exps)Exec(ctx *Context)error  {
	for i:=0;i< len(es) ;i++{
		err:=(es)[i].Exec(ctx)
		if err !=nil{
			return err
		}
	}
	return nil
}

//for

type ForExp struct {
	Addtion *BoolValue
	Do      Exp
}

func (f *ForExp)Exec(ctx *Context)error  {
	for f.Addtion.Match(ctx){
		if err:=f.Do.Exec(ctx);err !=nil{
			if err == breakError{
				break
			}
			return err
		}
	}
	return nil
}



type BreakExp struct {

}

func (b *BreakExp)Exec(ctx *Context)error  {
	return breakError
}



type DataExp struct {
	Data interface{}
	Key string
}

func (b *DataExp)Exec(ctx *Context)error  {
	ctx.table[b.Key] = b.Data
	return nil
}
//执行函数
type FuncExp struct{
	Value Value
}
func (b *FuncExp)Exec(ctx *Context)error  {
	b.Value.Get(ctx)
	return nil
}

//////
func parseBoolExp( s string)(*BoolValue,error){
	v,err:=parseValue(s)
	if err !=nil{
		return nil,err
	}
	return &BoolValue{Value: v},nil
}

/////
type Expr interface {
	Match() bool
}

type Op interface {
	Equal(x,y Value,ctx *Context)bool
}



//type BoolValue struct {
//	X Value
//	Op Op
//	Y Value
//}
//
//func (b *BoolValue)Match(ctx *Context)bool  {
//	return b.Op.Equal(b.X,b.Y,ctx)
//}
//
//func (b *BoolValue)Get(ctx *Context)interface{}  {
//	return b.Op.Equal(b.X,b.Y,ctx)
//}
//
//
//
////a==b
//// a==b
//type EqualOp struct {
//
//}
//
//func (o *EqualOp)Equal(x,y Value,ctx *Context)bool  {
//	X:=x.Get(ctx)
//	Y:=y.Get(ctx)
//	//fmt.Println("bool op:",X,Y)
//	return fmt.Sprintf("%v",X)==fmt.Sprintf("%v",Y)
//}
//// a && b
//type AndOp struct {
//
//}
//
//func (o *AndOp)Equal(x,y Value,ctx *Context)bool  {
//	X:=x.Get(ctx)
//	Y:=y.Get(ctx)
//	if xb,ok:=X.(bool);ok && xb{
//		if yb,ok:=Y.(bool);ok && yb{
//			return true
//		}
//	}
//	return false
//}
//
//
//// a==b && c == d
//func parseBoolExps( s string){
//	//s  = strings.Trim(s," ")
//	//token:=make([]byte,0, len(s))
//	//for i:=0;i< len(s);i++{
//	//	v:=s[i]
//	//}
//}
//
//func parseOp(s string)Op  {
//	switch s {
//	case "==":
//		return &EqualOp{}
//	case "&&":
//		return &AndOp{}
//
//	}
//	return nil
//}
//
