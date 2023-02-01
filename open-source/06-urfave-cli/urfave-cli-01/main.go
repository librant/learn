package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	// 实例化一个命令行程序
	oApp := cli.NewApp()

	oApp.Name = "GoTool"
	oApp.Usage = "To save the world"
	oApp.Version = "1.0.0"

	// 程序实际执行体
	oApp.Action = func(c *cli.Context) error {
		fmt.Println("Test")
		return nil
	}

	// 命令执行
	if err := oApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
