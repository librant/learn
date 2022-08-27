package main

import (
	"fmt"
	"reflect"
)

type Dog struct {
	legCount int
}

type Cat struct {
	LegCount int
}

func main() {
	//获取dog实例的反射值对象
	dog := reflect.ValueOf(&Dog{
		legCount: 4,
	}).Elem()
	//获取legCount字段的值
	dogLegCount := dog.FieldByName("legCount")
	fmt.Printf("dogLegCount: %v\n", dogLegCount)
	if dogLegCount.CanSet() {
		// 非导出字段， 不支持设置
		dogLegCount.SetInt(3)
		fmt.Printf("set dogLegCount: %v\n", dogLegCount)
	}
	//获取 cat 实例的反射值对象
	cat := reflect.ValueOf(&Cat{
		LegCount: 4,
	}).Elem()
	//获取legCount字段的值
	catLegCount := cat.FieldByName("LegCount")
	fmt.Printf("catLegCount: %v\n", catLegCount)
	if catLegCount.CanSet() {
		// 导出字段，支持修改
		catLegCount.SetInt(3)
		fmt.Printf("set catLegCount: %v\n", catLegCount)
	}
}
