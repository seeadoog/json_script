package jsonscpt

import (
	"strconv"
	"strings"
)

type JValue interface {
	String() string
	Bool() bool
	Float() float64
	Int()  int
}


type Table struct {
	table []Value
}

type JNumber float64

func (n *JNumber)Get()interface{}{

	return *n
}

type JString string

func (s  *JString)String()string{
	return string(*s)
}

func (s  *JString)Float()float64{
	return Number(*s)
}

func (s *JString)Bool()bool{
	return Bool(*s)
}

func (s *JString)Int()int{
	return int(Number(*s))
}

type JFloat float64

func ParseJFloat(v float64)JValue{
	jv:=JFloat(v)
	return &jv
}

func (f *JFloat)String()string{
	return strconv.FormatFloat(float64(*f),'f',-1,64)
}

func(f *JFloat)Bool()bool{
	if *f>0{
		return true
	}
	return false
}

func(f *JFloat)Int()int{
	return int(*f)
}

func(f *JFloat)Float()float64{
	return float64(*f)
}

type JFunc func(... JValue)JValue

func (s  JFunc)String()string{
	return ""
}

func (s  JFunc)Float()float64{
	return 0
}

func (s JFunc)Bool()bool{
	return false
}

func (s JFunc)Int()int{
	return 0
}

func jadd(vs... JValue)JValue{
	if len(vs) >=2{
		r:=JFloat(vs[0].Float()+vs[1].Float())
		return &r
	}
	return ParseJFloat(0)
}

type JMap map[string]JValue

func (m JMap)Set(k string,v JValue){
	m[k] = v
}
func (m JMap)Get(k string)JValue{
	return m[k]
}
func (m JMap)ExecFunc(name string,values ...JValue)JValue{
	if fun,ok:=m[name].(JFunc);ok{
		return (fun)(values...)
	}
	return nil
}

func (m JMap)String()string{
	return ""
}

func (m JMap)Int()int{
	panic("cannot use map as int value")
}

func (m JMap)Float()float64{
	panic("cannot use map as float value")
}

func (m JMap)Bool()bool{
	if m ==nil{
		return true
	}
	return false
}

// ValueAbstract

type AbstractValue interface{
	Get(c *JContext)JValue
}
// j const value
type JConstValue struct{
	Val JValue
}

func (v *JConstValue)Get(c *JContext)JValue{
	return v.Val
}

type JVariableValue struct{
	Name string
}
func (v *JVariableValue)Get(c *JContext)JValue{
	return c.Get(v.Name)
}

type JFuncValue struct{
	Name string
	Params []AbstractValue
}

func (v *JFuncValue)Get(c *JContext)JValue{
	p:=make([]JValue,0,len(v.Params))
	for i:=0 ;i< len(v.Params);i++{
		p = append(p,v.Params[i].Get(c))
	}
	fun:=c.GetFunc(v.Name)
	if fun !=nil{
		return fun(p...)
	}
	return nil
}
//JContext

type JContext struct{

}

func (c *JContext)Get(s string)JValue{
	return nil
}
func (c *JContext)GetFunc(s string)JFunc{
	return nil
}

func (c *JContext)Set(s string,v JValue){
	return
}

//exp

type JExp interface{
	Exec(ctx *JContext) error
}

type JSetExp struct{
	Name string
	Value AbstractValue
}

func (e *JSetExp)Exec(ctx *JContext) error{
	ctx.Set(e.Name,e.Value.Get(ctx))
	return nil
}

type JExps []JExp

func (e *JExps)Exec(ctx *JContext) error{
	for i:=0 ;i<len(*e);i++{
		err:=(*e)[i].Exec(ctx)
		if err !=nil{
			return err
		}
	}
	return nil
}

type JFuncExp struct{
	Value JFuncValue
}

func (e *JFuncExp)Exec(ctx *JContext) error{
	v:=e.Value.Get(ctx)
	if err,ok:=v.(error);ok{
		return err
	}
	return nil
}

//error
func Contains(tokens []string,s string)(string,bool){
	for _, v := range tokens {
		if strings.Contains(s,v){
			return v,true
		}
	}
	return "", false
}

func parseTokens(s string,tokenSet []string){

}




type TableValue struct {
	Table JMap
	Key string
}

