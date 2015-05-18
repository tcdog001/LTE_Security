package controllers

import (
	"LTE_Security/models"
	"github.com/astaxie/beego"
)

type ModifyController struct {
	beego.Controller
}

func (this *ModifyController) Get() {
	//check islogin
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

	this.TplNames = "modify.html"
	mac := this.Input().Get("mac")
	this.Data["Mac"] = mac
	this.Data["HadModified"] = false
}

func (this *ModifyController) Post() {
	///check islogin
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	//get form info
	mac := this.Input().Get("mac")
	newmac := this.Input().Get("newmac")

	deviceinfo := models.Deviceinfo{
		Mac: mac,
	}

	//check mac address format
	if !CheckMacFormat(newmac) {
		beego.Error(newmac, "not a mac address!")
		router := "/modify?mac="
		router += mac
		this.Redirect(router, 302)
		return
	}
	//write to db
	ok := models.ModifyDevice(&deviceinfo, newmac)
	if ok {
		beego.Info("ModifyDevice success!")
		this.Redirect("/noregister", 301)
		return
	}
	beego.Info("ModifyDevice failed! Once again")
	this.Data["HadModified"] = false
}
