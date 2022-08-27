package main

import (
	"fmt"
	"reflect"
)

func main() {
	var a int = 1024

	//获取变量a的反射值对象
	valueOfA := reflect.ValueOf(a)
	//获取interface{}类型的值，通过类型断言转换
	var getA int = valueOfA.Interface().(int)
	var getB int = int(valueOfA.Int())
	fmt.Printf("getA = %d getB = %d\n", getA, getB)
}
