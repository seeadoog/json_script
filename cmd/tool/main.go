package main

import (
	"flag"
	"crypto/sha256"
	"encoding/hex"
	"bytes"
	"errors"
	"fmt"
	"time"
)
var file = flag.String("f","","")
func main(){
	bm:=&Bm{
		blocks:map[string]*Block{},
		difficult:10,
	}
	b:=createBlock(bm)
	fmt.Println(b.CalculateHash())
	//flag.Parse()
	//vm:=jsonscpt.NewVm()
	//if *file==""{
	//	sc:=bufio.NewScanner(os.Stdin)
	//	for sc.Scan(){
	//		func (){
	//			defer func() {
	//				if err := recover();err !=nil{
	//					fmt.Println(err)
	//				}
	//			}()
	//			vm.ExecJsonObject(vm.ExecJsonObject(sc.Text()))
	//		}()
	//	}
	//	return
	//}
	//b,err:=ioutil.ReadFile(*file)
	//if err !=nil{
	//	fmt.Println(err.Error())
	//	return
	//}
	//cmd,err:=jsonscpt.CompileExpFromJson(b)
	//if err !=nil{
	//	fmt.Println(err)
	//	return
	//}
	//start:=time.Now()
	//for i:=0;i<1;i++{
	//	if err:=vm.SafeExecute(cmd, func(err interface{}) {
	//		fmt.Println(err)
	//	});err !=nil{
	//		fmt.Println(err)
	//	}
	//}
	//
	//fmt.Println("total cost=>",time.Since(start))

}


func hash(b []byte)string{
	h :=sha256.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

type Block struct{
	PreHash string
	Nonce uint32
	Base string
	Difficult byte
	PK string
	Time int32
}
func int2bytes(i uint32)[]byte{
	b:=make([]byte,4)
	b[0] = byte((i>>24) & 0xff)
	b[1] = byte((i>>16) & 0xff)
	b[2] = byte((i>>8) & 0xff)
	b[3] = byte((i) & 0xff)
	return b
}
func (b *Block)CalculateHash()string{
	bf:=&bytes.Buffer{}
	bf.WriteString(b.PreHash)
	bf.Write(int2bytes(b.Nonce))
	bf.WriteRune(b.Time)
	bf.WriteString(b.Base)
	bf.WriteByte(b.Difficult)
	bf.WriteString(b.PK)
	return hash(bf.Bytes())
}

func validate(b *Block){
	b.CalculateHash()
}


func checkDifficulty(s string)byte{
	b,_:=hex.DecodeString(s)
	var d byte = 0
	for i:=0;i< len(b);i++{
		if b[i]==0{
			d++
		}else{
			return d
		}

	}
	return d
}
type Bm struct {
	blocks map[string]*Block
	difficult byte
}

func (b *Bm)GetLen(hash string)int{
	n := 0
	h:=hash
	for {
		 bl:= b.blocks[h]
		 if bl == nil{
		 	return n
		 }
		 h = bl.PreHash
	}
}

func (b *Bm)AddBlock(block *Block)error{
	pre:=b.blocks[block.PreHash]
	if pre ==nil && len(b.blocks)>0{
		return errors.New("invalid block")
	}
	hsh:=block.CalculateHash()
	if checkDifficulty(hsh)<b.difficult{
		return errors.New("difficult is not enough")
	}

	b.blocks[hsh] = block
	return nil
}

func (b *Bm)GetLongestChain()*Block{
	len:=0
	var block *Block
	for _, v := range b.blocks {
		l := b.GetLen(v.CalculateHash())
		if l>len{
			len = l
			block = v
		}
	}
	return block
}




func createBlock(bm *Bm)*Block{
	lb:=bm.GetLongestChain()
	pre:=""
	if lb!=nil{
		pre = lb.CalculateHash()
	}
	b:=&Block{
		PreHash:pre,
		Base:"15",
		PK:"123456",
		Difficult:10,
	}
	b.Time = int32(time.Now().UnixNano())
	for i:=0;i<0x7ffffff;i++{
		b.Nonce = uint32(i)
		hash:=b.CalculateHash()
		if len(hash)!=64{
			panic(hash)
		}
		df:=checkDifficulty(hash)
		//fmt.Println(hash)
		if df >= b.Difficult{
			return b
		}

	}
	panic("cannot find block")
	return nil
}
