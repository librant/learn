package main

import (
	"errors"
	"fmt"
)

type ErrorString struct {
	s string
}

func (e *ErrorString) Error() string {
	return e.s
}

func main() {
	var targetErr *ErrorString
	err := fmt.Errorf("new error:[%w]", &ErrorString{s: "target err"})
	fmt.Println(errors.As(err, &targetErr))
}
