package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("test the os package")
	fmt.Println("Args len", len(os.Args))

	for i, arg := range os.Args {
		fmt.Println("arg", i)
		fmt.Println("arg", arg)
		fmt.Println(filepath.Base(arg))
	}

	fileInfo, err := os.Stat("./test.txt")
	if err != nil {
		fmt.Println(err)
	}

	Date := fileInfo.ModTime().String()
	fmt.Println("Date", Date)
}
