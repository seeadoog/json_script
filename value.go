package jsonscpt

import (
	"strings"
	"strconv"
	"regexp"
	"errors"
)

// const value

type Value interface {
	Get(ctx *Context) interface{}
}
type ConstValue struct {
	v interface{}
}

func (v *ConstValue) Get(ctx *Context) interface{} {
	return v.v
}

//this is variable ,the value from context
type VarValue struct {
	Key string
}

func (v *VarValue) Get(ctx *Context) interface{} {
	o, _ := CachedJsonpathLookUp(ctx.table, v.Key)
	return o
}

// value of function
type FuncValue struct {
	FuncName string  // func name
	params   []Value //params
}

func (v *FuncValue) Get(ctx *Context) interface{} {
	var params = make([]interface{}, 0, len(v.params))
	for _, v := range v.params {
		params = append(params, v.Get(ctx))
	}
	if fun, ok := ctx.table[v.FuncName].(Func); ok {
		return fun(params...)
	}
	return nil
}
//
var funv = regexp.MustCompile(`(\w+)\((.+)\)$`)
func parseValue(s string) (Value,error) {
	s = strings.Trim(s," ")
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
	//set from function
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
			return &FuncValue{FuncName:key,params:parsedValues},nil
		}
	}

	if !checkRule(s){
		return nil,errors.New("invalid variable key:"+s)
	}
	return &VarValue{Key:s},nil
//	return nil, nil
}

func toTokens(s string) []string {

	for i := 0; i < len(s); i++ {

	}
	return nil
}
