package main

import (
	"log"
	"unsafe"
)

// slice 切片的定义
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}

func main() {
	// slice 变量本身的地址不会变, slice 中的 array 指针会因为 cap 容量不够后会重新分配；
	a := make([]string, 0, 100)
	log.Printf("[1]&a = %p", &a) // a 变量底层对应的是数组，切片 a 的地址不随切片的长度变动而变动；
	log.Printf("[1]a = %p", a)   // a 的值对应底层数组的值，在 cap 容量内是不变的，超过 cap 时，由 append 改变
	for i := 0; i < 100; i++ {
		a = append(a, "I do not know ")
	}
	log.Printf("[2]&a = %p", &a)
	log.Printf("[2]a = %p", a)
	for i := 0; i < 1; i++ {
		a = append(a, "I do not know ")
	}
	log.Printf("[3]&a = %p", &a)
	log.Printf("[3]a = %p", a)
	return
}
