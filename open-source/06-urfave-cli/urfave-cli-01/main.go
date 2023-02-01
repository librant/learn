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

	oApp.Commands = []cli.Command{
		{
			Name: "lang", 				// 命令全称
			Aliases: []string{"l"},		// 命令简称
			Usage:"Setting language",	// 命令详细描述
			Action: func(c *cli.Context) {
				// 通过c.Args().First()获取命令行参数
				fmt.Printf("language=%v \n",c.Args().First())
			},
		},
		{
			Name: "encode",
			Aliases: []string{"e"},
			Usage:"Setting encoding",
			Action: func(c *cli.Context) {
				fmt.Printf("encoding=%v \n",c.Args().First())
			},
		},
	}

	// 命令执行
	if err := oApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
