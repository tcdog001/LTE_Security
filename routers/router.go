package routers

import (
	"LTE_Security/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.LoginController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/home", &controllers.HomeController{})
	beego.Router("/register", &controllers.RegisterController{})
	beego.Router("/verify", &controllers.VerifyController{})
	beego.Router("/modify", &controllers.ModifyController{})
	beego.Router("/noregister", &controllers.NoregisterController{})
	beego.Router("/import", &controllers.ImportController{})
}
