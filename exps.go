package jsonscpt

import (
	"unsafe"
	"strings"
	"errors"
	"fmt"
)

type Exp interface  {
	Exec(ctx *Context)error
}
//add ( 1 , len ( 'd' ) )
func parseSetExp(s string)(*SetExps, error){
	v:=strings.SplitN(s,"=",2)
	if len(v)!=2{
		if len(v)==1{
			val,err:=parseValue(v[0])
			if err !=nil{
				return nil,err
			}
			return &SetExps{Variable:"_",Value:val},nil
		}
		return nil,errors.New("invalid set exp:"+s)
	}

	if !checkRule(v[0]){
		return nil,errors.New("invalid setexp key:"+v[0])
	}
	if isSystemId(v[0]){
		return nil,errors.New("system identify can not be variable:"+s)
	}
	val,err:=parseValue(v[1])
	if err !=nil{
		return nil,err
	}
	return &SetExps{Variable:v[0],Value:val},nil
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
type BoolExp struct {
	//Op string  // > ,== ,< ,>= ,<=
	Value Value
}

func (b *BoolExp)Match(ctx *Context)bool  {
	return convertToBool(b.Value.Get(ctx))
}

//func (b *BoolExp)ExecJsonObject(ctx *Context)error{
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







func parseBoolExp( s string)(*BoolExp  ,error){
	v,err:=parseValue(s)
	if err !=nil{
		return nil,err
	}
	return &BoolExp{Value:v},nil
}






















type Expr interface {
	Match() bool
}

type Op interface {
	Equal(x,y Value,ctx *Context)bool
}

type BoolValue struct {
	X Value
	Op Op
	Y Value
}

func (b *BoolValue)Match(ctx *Context)bool  {
	return b.Op.Equal(b.X,b.Y,ctx)
}

func (b *BoolValue)Get(ctx *Context)interface{}  {
	return b.Op.Equal(b.X,b.Y,ctx)
}



//a==b
// a==b
type EqualOp struct {

}

func (o *EqualOp)Equal(x,y Value,ctx *Context)bool  {
	X:=x.Get(ctx)
	Y:=y.Get(ctx)
	//fmt.Println("bool op:",X,Y)
	return fmt.Sprintf("%v",X)==fmt.Sprintf("%v",Y)
}
// a && b
type AndOp struct {

}

func (o *AndOp)Equal(x,y Value,ctx *Context)bool  {
	X:=x.Get(ctx)
	Y:=y.Get(ctx)
	if xb,ok:=X.(bool);ok && xb{
		if yb,ok:=Y.(bool);ok && yb{
			return true
		}
	}
	return false
}


// a==b && c == d
func parseBoolExps( s string){
	//s  = strings.Trim(s," ")
	//token:=make([]byte,0, len(s))
	//for i:=0;i< len(s);i++{
	//	v:=s[i]
	//}
}

func parseOp(s string)Op  {
	switch s {
	case "==":
		return &EqualOp{}
	case "&&":
		return &AndOp{}

	}
	return nil
}

type ForExp struct {

}