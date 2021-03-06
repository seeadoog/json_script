package jsonscpt

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// const value

type Value interface {
	Get(ctx *Context) interface{}
	GetName()string
}
// a=5 ,b='124',c=true
type ConstValue struct {
	v interface{}
}

func (v *ConstValue) Get(ctx *Context) interface{} {
	return v.v
}
func (v *ConstValue)GetName()string  {
	return ""
}
//this is variable ,the value from context
// ent = o.ent
type VarValue struct {
	Key string
}

func (v *VarValue) Get(ctx *Context) interface{} {
	o, err := CachedJsonpathLookUp(ctx.table, v.Key)
	if  err !=nil{
		if f,ok:=ctx.funcs[v.Key];ok{
			return f
		}
		return nil
		//return ctx.funcs[v.Key]
	}
	return o
}

func (v *VarValue)GetName()string  {
	return v.Key
}

// value of function
// a = add(1,3)
type FuncValue struct {
	FuncName string  // func name
	Params   []Value //Params
}

func (v *FuncValue) Get(ctx *Context) interface{} {
	var params = make([]interface{}, 0, len(v.Params))
	for _, v := range v.Params {
		ev:=v.Get(ctx)
		if err,ok:=ev.(*ErrorExit);ok{
			return err
		}
		params = append(params,ev)
	}
	if fun, ok := ctx.funcs[v.FuncName]; ok {
		return fun(params...)
	}else if fun, ok := ctx.table[v.FuncName].(Func); ok{
		return fun(params...)
	}
	panic("func:"+v.FuncName+" does not exists")
	return nil
}

func (v *FuncValue)GetName()string  {
	return v.FuncName
}



var funv = regexp.MustCompile(`(\w+)\((.+)*\)$`)

var eqexp =regexp.MustCompile(`(.+)==(.+)`)
func parseValue(s string) (Value,error) {
	s = strings.Trim(s," ")
	//if !isLegalBlock(s){
	//	return nil,errors.New("illegal block exp:"+s)
	//}
	//s= TrimBlock(s)
	if strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'") {
	 	v:=&ConstValue{
	 		v:strings.Trim(s,"'"),
		}
		return v ,nil
	}
	if iv,err:= strconv.ParseFloat(s,64);err ==nil{
		//fmt.Println("d")
		return &ConstValue{iv}  ,nil
	}

	if s=="nil"{
		return &ConstValue{nil},nil
	}

	if s=="true"{
		return &ConstValue{true},nil
	}

	if s=="false"{
		return &ConstValue{false},nil
	}
	//if v,err:= parseBoolValues(s);err ==nil{
	//	return v,nil
	//}

	if funv.MatchString(s){
		sub:=funv.FindAllSubmatch([]byte(s),-1)
		if len(sub)>0{
			key:=string(sub[0][1])
			values:=string(sub[0][2])
			var vals []string
			var token = make([]byte,0, len(values))
			var cs = 0
			var bl = 0
			for i:=0;i< len(values);i++{
				v:=values[i]
				token=append(token,v)
				if v !='\'' && cs ==1{
					continue
				}
				if v=='('{
					bl ++
					continue
				}

				if v==')'{
					if bl==0{
						return nil,errors.New("invalid exp:"+s)
					}
					bl--
					continue
				}

				if v=='\'' && cs==0{
					cs = 1
					continue
				}

				if v=='\'' && cs ==1{
					cs = 0
					continue
				}

				if  v==',' && len(token)>0 && cs==0 && bl==0{
					vals = append(vals,string(token[0:len(token)-1]))
					token = token[:0]
				}
			}
			if bl>0{
				return nil,errors.New("right ')' is required:"+s)
			}
			if cs>0{
				return nil,errors.New(" ' is not valid:"+s)
			}
			if len(token)>0{
				vals =  append(vals,string(token))
			}
			parsedValues:=make([]Value, 0,len(vals))
			for _, vs := range vals {
				  pv,err:=parseValue(vs)
				  if err !=nil{
				  	return nil,err
				  }
				  parsedValues = append(parsedValues,pv)
			}
			return &FuncValue{FuncName:key, Params:parsedValues},nil
		}
	}

	if eqexp.MatchString(s){
		v,err :=parseEqexp(s)
		if err !=nil{
			return nil,err
		}
		return v,nil
	}

	if !checkRule(s){
		return nil,errors.New("invalid variable key:"+s)
	}
	//fmt.Println(s)
	if s=="$root"{
		//todo
	}
	return &VarValue{Key:s},nil
//	return nil, nil
}

//todo this function has bug
func TrimBlock(s string)string{
	si:=0
	for i:=0;i< len(s);i++{
		v:=s[i]
		if v=='('{
			si++
		}else if v==')'{
			si--
		}else if v==' '{

		}else{
			break
		}
	}
	return trimBlock(s,si)
}
//合法
func isLegalBlock(s string)bool  {
	si:=0
	stack:=make([]byte, len(s))
	for i:=0;i< len(s);i++{
		v:=s[i]
		if v=='('{
			stack[si]= v
			si++
		}
		if v==')'{
			si--
			if si<0 {
				return  false
			}
		}
	}
	return si==0
}

func trimBlock (s string,i int)string{
	if (i<=0){
		return s
	}
	if strings.HasPrefix(s,"(") && strings.HasSuffix(s,")"){
		return trimBlock(s[1:len(s)-1],i-1)
	}else{
		return  s
	}
	return s
}




func parseEqexp(s string)(Value,error){
	r:=strings.SplitN(s,"==",2)
	v1,err:=parseValue(r[0])
	if err != nil{
		return nil,err
	}

	v2,err:=parseValue(r[1])
	if err != nil{
		return nil,err
	}

	return &FuncValue{
		FuncName: "eq",
		Params:   []Value{v1,v2},
	},nil
}