package main

// 它递归调用 Unwrap() 并判断每一层的 err 是否相等，如果有任何一层 err 和传入的目标错误相等，则返回 true。

import (
	"errors"
	"fmt"
)

func main() {
	err1 := errors.New("new error")
	err2 := fmt.Errorf("err2: [%w]", err1)
	err3 := fmt.Errorf("err3: [%w]", err2)

	fmt.Println(errors.Is(err3, err2))
	fmt.Println(errors.Is(err3, err1))
}
