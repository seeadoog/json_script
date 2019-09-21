package jsonscpt

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func init(){
	SetInlineFunc("open_file",open_file)
	SetInlineFunc("close",close)
	SetInlineFunc("read",read)
	SetInlineFunc("range_reader",rangeReader)
	SetInlineFunc("read_file",read_file)
	SetInlineFunc("write",write)
	SetInlineFunc("create_file",create_file)
}

type Err map[string]interface{}

func open_file(i ...interface{})interface{}{
	if len(i)==0{
		panic("open args should more than 1")
	}
	f,err:=os.OpenFile(String(i[0]),os.O_WRONLY,0666)
	if err !=nil{
		panic(err.Error())
	}
	return f
}

func create_file(i ...interface{})interface{}{
	if len(i)==0{
		panic("open args should more than 1")
	}
	f,err:=os.Create(String(i[0]))
	if err !=nil{
		panic(err.Error())
	}
	return f
}

func read_file(i ...interface{})interface{}{
	if len(i)==0{
		panic("read file args ==0")
	}
	b,err:=ioutil.ReadFile(String(i[0]))
	if err !=nil{
		panic(err.Error())
	}
	return b
}

func close(i ...interface{})interface{}{
	if len(i)==0{
		panic("open args should more than 1")
	}
	return i[0].(io.Closer).Close()
}

func read(i ...interface{})interface{}{
	if len(i)<2{
		panic("read args should more than 2")
	}
	len:=int(Number(i[1]))
	if len <=0{
		panic("read len cannot be zero:"+strconv.Itoa(len))
	}
	buf:=make([]byte,len)
	n,err:=i[0].(io.Reader).Read(buf)
	if err!=nil{
		if err ==io.EOF{
			return "EOF"
		}
		panic(err)
	}
	return n
}

func write (i ...interface{})interface{}{
	if len(i)<2{
		panic("write args <2")
	}
	_,err:=i[0].(io.Writer).Write(i[1].([]byte));
	if err !=nil{
		panic(err)
	}
	return nil

}

func rangeReader(i ...interface{})interface{}{
	if len(i)<3{
		panic("read args should more than 2")
	}

	f:=i[0].(io.ReadCloser)
	defer f.Close()
	len:=int(Number(i[1]))
	if len <=0{
		panic("read len cannot be zero:"+strconv.Itoa(len))
	}
	b:=make([]byte,len)
	for{
		l,err:=f.Read(b)
		if err!=nil{
			if err ==io.EOF{
				i[2].(Func)(b[:l],1)
				return "EOF"
			}
			panic(err)
		}
		if l<len{
			i[2].(Func)(b[:l],1)
			break
		}else{
			i[2].(Func)(b[:l],0)
		}
	}
	return nil
}


