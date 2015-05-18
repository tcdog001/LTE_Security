package controllers

import (
	"LTE_Security/models"
	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	this.TplNames = "login.html"
	this.Data["IsMatched"] = false
}

func (this *LoginController) Post() {
	//get form info
	uname := this.Input().Get("uname")
	pwd := this.Input().Get("pwd")

	admin := models.Admininfo{
		Username: uname,
		Password: pwd,
	}
	//check username and password
	ok := models.CheckAdmin(&admin)
	if ok {
		beego.Info("user/pwd matched!")

		this.SetSession("Admin", uname)
		if models.UpdateAdminStatus(&admin) {
			beego.Info("UpdateAdminStatus success!")
		} else {
			beego.Info("UpdateAdminStatus failed!")
		}

		beego.Info("Login success!")

		this.Data["IsMatched"] = true
		this.Redirect("/home", 301)
		return
	}

	beego.Info("Login failed! Once again!")
	this.Data["IsMatched"] = false
	this.Redirect("/login", 302)

	return
}
