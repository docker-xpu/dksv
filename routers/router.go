// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"dksv/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/api/test/", &controllers.MainController{})

	// 容器操作
	beego.Router("/api/container/run/", &controllers.Container{}, "post:Run")
	beego.Router("/api/container/stop/", &controllers.Container{}, "post:Stop")
	beego.Router("/api/container/remove/", &controllers.Container{}, "post:Remove")
	beego.Router("/api/container/list/", &controllers.Container{}, "get:List")
	beego.Router("/api/container/logs/", &controllers.Container{}, "get:Logs")
}
