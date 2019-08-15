package jsonscpt

import (
	"unsafe"
	"strings"
	"errors"
)

type Exp interface  {
	Exec(ctx *Context)error
}
//add ( 1 , len ( 'd' ) )
func parseSetExp(s string)(*SetExps, error){
	v:=strings.Split(s,"=")
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
	return MarshalInterface(e.Variable,ctx.table,e.Value.Get(ctx))
}
//a>b && (c>d)
type BoolExp struct {
	Op string  // > ,== ,< ,>= ,<=
	Left Value
	Right Value
	Parent *BoolExp
	LeftChild *BoolExp
	RightChild *BoolExp
}

func (b *BoolExp)Match()bool  {
	return true
}

func parseBoolExp( s string)(*BoolExp  ,error){
	return &BoolExp{},nil
}
