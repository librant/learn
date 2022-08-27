package main

import (
	"fmt"
	"reflect"
)

type Cat struct {
	Name string `json:"name" name:"buou"`
	Type int    `json:"type" id:"100"`
}

func main() {
	ins := Cat{Name: "mimi", Type: 1}
	typeOfCat := reflect.TypeOf(ins)
	// 遍历结构体中的所有成员
	for i := 0; i < typeOfCat.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		fieldType := typeOfCat.Field(i)
		// 输出成员名和tag
		fmt.Printf("name: %v  tag: '%v'\n", fieldType.Name, fieldType.Tag)
	}
	// 通过字段名, 找到字段类型信息
	if catName, ok := typeOfCat.FieldByName("Name"); ok {
		// 从tag中取出需要的tag
		fmt.Println(catName.Tag.Get("json"), catName.Tag.Get("name"))
	}
	if catType, ok := typeOfCat.FieldByName("Type"); ok {
		fmt.Println(catType.Tag.Get("json"), catType.Tag.Get("id"))
	}
}
