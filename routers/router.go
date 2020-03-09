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
	beego.Router("/api/container/run/", &controllers.ContainerController{}, "post:Run")
	beego.Router("/api/container/stop/", &controllers.ContainerController{}, "post:Stop")
	beego.Router("/api/container/remove/", &controllers.ContainerController{}, "post:Remove")
	beego.Router("/api/container/list/", &controllers.ContainerController{}, "get:List")
	beego.Router("/api/container/logs/", &controllers.ContainerController{}, "get:Logs")
	beego.Router("/api/container/commit/", &controllers.ContainerController{}, "post:Commit")

	// 镜像操作
	beego.Router("/api/image/list/", &controllers.ImageController{}, "get:List")
	beego.Router("/api/image/pull/", &controllers.ImageController{}, "post:Pull")
	beego.Router("/api/image/remove/", &controllers.ImageController{}, "post:Remove")

	// 网络操作
	beego.Router("/api/network/create/", &controllers.NetworkController{}, "post:Create")
	beego.Router("/api/network/list/", &controllers.NetworkController{}, "get:List")
	beego.Router("/api/network/remove/", &controllers.NetworkController{}, "post:Remove")
}
