package main

import (
	"os"

	"k8s.io/component-base/cli"

	"github.com/librant/learn/k8s/client-go/demon/cobra-rest-client/app"
)

func main() {
	command := app.NewRestClientCommand()
	code := cli.Run(command)
	os.Exit(code)
}
