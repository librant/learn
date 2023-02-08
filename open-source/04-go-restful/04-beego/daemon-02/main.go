package main

import (
	beego "github.com/beego/beego/v2/server/web"

	_ "github.com/librant/learn/open-source/04-go-restful/04-beego/daemon-02/routers"
)

func main() {
	beego.Run()
}
