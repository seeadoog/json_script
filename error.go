package jsonscpt

import "fmt"

type ErrorReturn struct {
	Message string
	Code int
}

func (e *ErrorReturn)Error()string  {
	return fmt.Sprintf("code=%d,message=%s",e.Code,e.Message)
}

var breakError = &BreakError{}
//break error
type BreakError struct {

}

func (e *BreakError)Error()string  {
	return "break"
}
