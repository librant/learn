package main

import (
	"log"

	"github.com/mitchellh/go-ps"
)

func main() {
	processes, err := ps.Processes()
	if err != nil {
		log.Fatal(err)
	}

	// 打印当前系统正在执行的命令
	for _, process := range processes {
		log.Printf("%s\n", process.Executable())
	}
}
