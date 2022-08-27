package main

import (
	"fmt"
	"reflect"
)

type Enum int

const(
	Zero Enum = 0
)

func main() {
	typeOfZero := reflect.TypeOf(Zero)
	fmt.Println(typeOfZero.Name(), typeOfZero.Kind())
	// Enum int
}
