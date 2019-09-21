package jsonscpt

import (
	"encoding/base64"
)
func init(){
	SetInlineFunc("base64decode",base64decode)
	SetInlineFunc("base64encode",base64encode)
}
func base64decode(i ...interface{}) interface{} {
	if len(i)>=1{
		b,err:=base64.StdEncoding.DecodeString(String(i[0]))
		if err !=nil{
			panic(err.Error())
		}
		return b
	}
	return nil
}

func base64encode(i ...interface{}) interface{} {
	if len(i)>=1{
		return base64.StdEncoding.EncodeToString(i[0].([]byte))
	}
	return nil
}