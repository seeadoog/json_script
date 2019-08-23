package jsonscpt

import (
	"fmt"
	"testing"
)

func TestTokenized3(t *testing.T)  {

	fmt.Println(tokeniz3("$.name\\\\.age[3].gg3\\.id"))
	fmt.Println(checkRule("$.name\\\\.age[e3].gg3\\.id"))

}

