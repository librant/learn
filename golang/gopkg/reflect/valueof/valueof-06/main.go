package main

import (
	"fmt"
	"reflect"
)

// 通过反射调用函数

func add(a, b int) int {
	return a + b
}

// 通过反射调用方法

type MyMath struct {
	Pi float64
}

func (myMath MyMath) Sum(a, b int) int {
	return a + b
}

func (myMath MyMath) Dec(a, b int) int {
	return a - b
}

func main() {
	//将函数包装为反射值对象
	funcValue := reflect.ValueOf(add)
	//构造函数参数，传入两个整形值
	paramList := []reflect.Value{reflect.ValueOf(2), reflect.ValueOf(3)}
	//反射调用函数
	retList := funcValue.Call(paramList)
	fmt.Println(retList[0].Int())

	var myMath = MyMath{Pi:3.14159}
	//获取myMath的值对象
	rValue := reflect.ValueOf(myMath)
	// 获取有多少个方法
	for i := 0; i < rValue.NumMethod(); i++ {
		fmt.Printf("method[%d] type: %v\n", i, rValue.Method(i).Type())
		paramList := []reflect.Value{reflect.ValueOf(30*(i+1)), reflect.ValueOf(20*(i+1))}
		result := rValue.Method(i).Call(paramList)
		fmt.Println(result[0].Int())
	}
}
