package jsonscpt

import (
	"strings"
	"strconv"
	"regexp"
	"errors"
	"fmt"
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
	if !isLegalBlock(s){
		return nil,errors.New("illegal block exp:"+s)
	}
	s= TrimBlock(s)
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

	//set from function
	//if boooreg.Match([]byte(s)){
	//	//r:=boooreg.FindAllSubmatch([]byte(s),-1)
	//	//fmt.Println(len(r))
	//	bv,err:=parseBoolV(s)
	//	if err !=nil{
	//		//goto other
	//	}
	//	return bv,nil
	//}
//other:
	//value from functon
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
	//fmt.Println(s)

	return &VarValue{Key:s},nil
//	return nil, nil
}

func parseBoolV(s string) (Value,error) {
	s = strings.Trim(s," ")
	bl:=0 //括号匹配
	eof:=0
	token :=make([]byte, 0,len(s))
	vals:=[]string{}
	//fmt.Println("-----------",s)
	for i:=0;i< len(s);i++{
		v:=s[i]
		token = append(token,v)
		if v=='('{
			bl++
			continue
		}
		if v==')'{
			if bl==0{
				return nil,errors.New("invalid bool exp==>"+s)
			}
			bl--
			continue
		}
		if bl==0{
			if (v=='&' || v=='>' || v=='<' || v=='='){
				//op = append(op,v)
				if eof ==0{
					if len(token)==0{
						return nil,errors.New("nivalid bool exp,start with eof:"+s)
					}
					vals = append(vals,string(token[0:len(token)-1]))
					token = token[len(token)-1:]
					eof = 1
				}
			}else{
				if eof == 1{
					eof = 0
					if len(token)==0{
						return nil,errors.New("invalid bool exp,start with eof:"+s)
					}
					vals = append(vals,string(token[0:len(token)-1]))
					token = token[len(token)-1:]
					//终结符扫描完成
				}
			}
		}
	}

	//括号不匹配
	if bl !=0{
		return nil,errors.New("invald bool exp,need ')':"+s)
	}

	if len(token)>0{
		vals = append(vals,string(token))
	}

	if len(vals)>3{
		return nil,errors.New("invalid bool exp:too many sub exp:"+s)
	}
	if len(vals) == 3{
		op:=parseOp(vals[1])
		if op==nil{
			return nil,errors.New("invalid bool op:"+vals[1])
		}
		x,err:=parseValue(vals[0])
		if err !=nil{
			return nil,err
		}
		y,err:=parseValue(vals[2])
		if err !=nil{
			return nil,err
		}
		return &BoolValue{X:x,Y:y,Op:op},nil
	}
	if len(vals) == 1{
		return parseValue(vals[0])
	}
	fmt.Println(vals)
	return nil,errors.New("invalid bool exp,"+s)
}

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
	for _, v := range s {
		if v=='('{
			si++
		}
		if v==')'{
			si--
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

var boooreg = regexp.MustCompile(`(.+)((==)|(&&)|(>=))(.+)`)
func toTokens(s string) []string {

	for i := 0; i < len(s); i++ {

	}
	return nil
}
