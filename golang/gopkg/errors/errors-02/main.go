package main

// Unwrap() 拆开一个被包装的 error

import (
	"errors"
	"fmt"
)

/*
// errors 包中的错误结构体信息

type WarpError struct {
	msg string
	err error
}

func (e *WarpError) Error() string {
	return e.msg
}

func (e *WarpError) Unwrap() error {
	return e.err
}

*/

func main() {
	err1 := errors.New("new error")
	err2 := fmt.Errorf("err2: [%w]", err1)
	err3 := fmt.Errorf("err3: [%w]", err2)

	fmt.Println(errors.Unwrap(err3))
	// err2: [new error]
	fmt.Println(errors.Unwrap(errors.Unwrap(err3)))
	// new error
}
