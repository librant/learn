package main

import (
	"fmt"
	"reflect"
)

func main() {
	var a int = 1024
	fmt.Printf("aPtr: %v\n", &a)

	rValuePtr := reflect.ValueOf(&a)
	if rValuePtr.Elem().CanSet() {
		rValuePtr.Elem().SetInt(100)
	}
	fmt.Printf("a = %v\n", rValuePtr.Elem())
	if rValuePtr.Elem().CanAddr() {
		fmt.Printf("rValuePtr: %v\n", rValuePtr)
	}
}
