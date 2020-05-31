package jsonscpt


type Switch struct {
	swh Value
	cases map[Value]Exp
	defaulte Exp
}

func (s Switch) Exec(ctx *Context) error {
	v:=s.swh.Get(ctx)
	for value, exp := range s.cases {
		if String(value.Get(ctx)) == String(v){
			return exp.Exec(ctx)
		}
	}
	if s.defaulte != nil{
		return s.defaulte.Exec(ctx)
	}
	return nil
}


