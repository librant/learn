package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name string
	Age  int

	// 嵌入字段
	float32
	bool

	next *Student
}

func main() {
	rValue := reflect.ValueOf(Student{
		Name:    "librant",
		Age:     18,
		float32: 1.01,
		bool:    true,
		next: &Student{
			Name:    "hello",
			Age:     100,
			float32: 1.32,
			bool:    false,
			next:    &Student{},
		},
	})
	// 获取字段的数量
	fmt.Println("NumField:", rValue.NumField())
	// 获取每个字段中的值
	for i := 0; i < rValue.NumField(); i++ {
		filed := rValue.Field(i)
		fmt.Printf("Field[%d] %v %v\n", i, filed, filed.Type())
	}
	// 根据名字查找字段
	fmt.Printf("FieldByName(Age).Type: %v\n", rValue.FieldByName("Age").Type())

	// 根据索引查找值中next字段的name字段的类型
	fmt.Printf("FieldByIndex([]int{4, 0}).Type(): %v\n", rValue.FieldByIndex([]int{4, 0}).Type())
	fmt.Printf("FieldByIndex([]int{4, 0}): %v\n", rValue.FieldByIndex([]int{4, 0}))
}
