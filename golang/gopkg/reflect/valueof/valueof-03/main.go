package main

import (
	"fmt"
	"reflect"
)

func main() {
	var a *int
	fmt.Printf("var a *int: %v\n", reflect.ValueOf(a).IsNil())

	fmt.Printf("nil: %v\n", reflect.ValueOf(nil).IsValid())

	// *int类型的空指针
	fmt.Printf("(*int)(nil): %v\n", reflect.ValueOf((*int)(nil)).Elem().IsValid())

	s := struct{}{}
	// 尝试从结构体中查找一个不存在的字段
	fmt.Println("不存在的结构体成员:", reflect.ValueOf(s).FieldByName("").IsValid())
	// 尝试从结构体中查找一个不存在的方法
	fmt.Println("不存在的方法:", reflect.ValueOf(s).MethodByName("").IsValid())

	//实例化一个map
	m := map[string]int{
		"name": 1,
		"age":  2,
	}
	//尝试从map中查找一个存在的键
	fmt.Println("不存在的键:", reflect.ValueOf(m).MapIndex(reflect.ValueOf("name")).IsValid())
	fmt.Println("不存在的键:", reflect.ValueOf(m).MapIndex(reflect.ValueOf("")).IsValid())
}
