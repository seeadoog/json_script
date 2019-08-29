#### a json script for go

#### use

install
````
go get https://github.com/seeadoog/json_script
````

````
	rule:=[]byte(`
[
  {
    "if": "lt(user.age,15)",
    "then": "user.generation='yong'"
  },
  {
    "if": "lt(user.age,30)",
    "then": [
      "user.generation='old'",
      "user.hasChild=true"
    ]
  },
  {
    "for":"k,v in user",
    "do":"print('k==',k,'v==',v)"
  },
  {
    "func":"show(u)",
    "do":"printf('name=%s,age=%v,generation=%s',u.name,u.age,u.generation)"
  },
  "show(user)"
]
`)
	scp,err:=jsonscpt.CompileExpFromJson(rule)
	if err !=nil{
		panic(err)
	}
	vm:=jsonscpt.NewVm()
	user :=map[string]interface{}{
		"name":"bob",
		"age":16,
	}
	vm.Set("user", user)
	err =vm.SafeExecute(scp,nil)
	if err !=nil{
		panic(err)
	}
	fmt.Println("userMap:",user)
````


#### inline functions

funs|use
---|---
append(str...)string|拼接字符串
split(str,sep,n)[]string|按照给定的切割符号 sep切分字符串 
len(str)int|求字符串的长度
sprintf(str,obj...)string|obj...|格式化字符串
add(nubmer ...)number|number array|数字求和
isnil(object)bool|判断对象是否为空
and(bool...)bool|对一系列的bool变量做 and操作 相当于 && ，任何变量都可以视为bool值<br>
or(bool ...)bool|对一系列的bool变量做 or 操作 相当于 '\\' ，任何变量都可以视为bool值<br>
eq(objectA,objectB)bool|判断两者转换为字符串后是否相等
gt(numberA,numberB)bool| 大于 > ,numberA>numberB 返回true，否则返回false
ge(numberA,numberB)bool| 大于等于 >= ,numberA>=numberB 返回true，否则返回false
lt(numberA,numberB)bool| 小于 < ,numberA<numberB 返回true，否则返回false
le(numberA,numberB)bool| 小于等于 <= ,numberA<=numberB 返回true，否则返回false
not(bool)bool|非 ！  如果bool=true 返回false 否则返回true
in(str,strs...)bool |如果str 在 strs中一个，返回true ，否则返回false  
contains(str,sub)| 如果str包含sub片段 ，返回true，否则返回false
join(str ... ,sep)string| 拼接字符串，用sep分割
get(object,key)obj| 返回数组的第key个元素，越界会panic，或者 map的key
set(object,key,value)| 设置数组的第key个元素，越界会panic，或者为map 设置键和值
exit(code,msg)|终止script的执行并返回一个error，包含code，和message，类型为*ErrorExit
return(obj)|放在函数中会作为返回值


#### syntax 语法

1. if else 如果if语句中的bool值为真执行then 部分，否则执行else部分

````
{
    "if":"eq(a,5)"
    "then":"print('a == 5')",
    "else":"print('a != 5')"

}
````

2. foreach 遍历map或者数组。
````
{
    "for":"k,v in user",
    "do":"print('k==',k,'v==',v)"
 },
````

3. 定义function  函数名(参数名....)

````
{
    "func":"show(u)",
    "do":"printf('name=%s,age=%v,generation=%s',u.name,u.age,u.generation)"
}
````

4. bool 值的判定

任何基础类型都可以视为bool值

string  ： 长度大于0 为 true 否则 false
number :  >0 为true  <=0 为false
bool 值：bool值

5. 赋值
//值得类型有三种，bool，string，number
````
[
"va=123", //给va 赋值为 number 类型
"vstr='hello world'" // 给vstr赋值为 string类型
"vb=true"   // 给vb赋值为bool类型，值为true
"param.a='hello'" // 给param map中的a赋值string =  'hello'
]


````