package jsonscpt

import (
	"reflect"
	"testing"
)

func Test_number(t *testing.T) {
	vm:=NewVm()
	vm.ExecJson([]byte(`
[{
  "func": "pow(i)",
  "do": [
    "return(i*i)"
  ]
}
,
  "print(pow(5)+pow(5)-pow(2*2+4))"
]`))

}



func TestConvertToString(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(tt.args.v); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lens(t *testing.T) {
	type args struct {
		i []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
		{args:args{i:[]interface{}{"hello world"}},want:11},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lens(tt.args.i...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lens() = %v, want %v", got, tt.want)
			}
		})
	}
}
