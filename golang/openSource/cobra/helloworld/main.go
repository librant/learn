package main

import (
	"log"

	"github.com/librant/learn/golang/openSource/cobrademon/helloworld/cmd"
)

func main() {
	log.Printf("cobra hello world init...")
	cmd.Execute()
}
