package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name string
	Age  int
}

func main() {
	var stu Student
	typeOfStu := reflect.TypeOf(stu)
	valueOfStu := reflect.ValueOf(stu)
	fmt.Println(typeOfStu.Name(), typeOfStu.Kind())
	fmt.Println(valueOfStu.Kind())
	// Student struct

	stuPtr := &Student{
		Name: "librant",
		Age:  30,
	}
	typeOfStuPtr := reflect.TypeOf(stuPtr)
	fmt.Println(typeOfStuPtr.Name(), typeOfStuPtr.Kind())
	// 空 ptr

	newStuPtr := typeOfStuPtr.Elem()
	fmt.Printf("stuPtr: %p typeOfStuPtr: %p\n", stuPtr, newStuPtr)
	//显示反射类型对象的名称和种类
	fmt.Printf("element name: '%v', element kind: '%v'\n", newStuPtr.Name(), newStuPtr.Kind())

}
