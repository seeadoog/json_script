package jsonscpt

import (
	"bytes"
	"strconv"
	"fmt"
	"strings"
)

type Func func(...interface{})interface{}

var lens Func = func(i ...interface{}) interface{} {
	if len(i)==0{
		return 0
	}
	v:=i[0]
	switch v.(type) {
	case string:
		return len(v.(string))
	case map[string]interface{}:
		return len(v.(map[string]interface{}))
	default:
		return 0
	}
}

var apd Func = func(i ...interface{}) interface{} {
	var bf bytes.Buffer
	for _, v := range i {
		bf.WriteString(ConvertToString(v))
	}
	return bf.String()
}


var split Func = func(i ...interface{}) interface{} {
	if len(i)>=2{
		if s,ok:=i[0].(string);ok{
			if sp,ok:=i[1].(string);ok{
				if len(i)>=3{
					if n,ok:=i[2].(float64);ok{
						return strings.SplitN(s,sp,int(n))
					}
				}else{
					return strings.Split(s,sp)
				}
			}
		}
	}
	return nil
}

func ConvertToString(v interface{}) string {
	switch v.(type) {
	case string:return v.(string)
	case float64:return strconv.FormatFloat(v.(float64),'f',-1,64)
	case bool:return strconv.FormatBool(v.(bool))
	case int: return strconv.Itoa(v.(int))
	default:
		return fmt.Sprintf("%v",v)
	}
}