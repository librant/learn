package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"github.com/librant/learn/open-source/04-go-restful/04-beego/daemon-02/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
