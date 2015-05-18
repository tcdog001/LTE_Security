package controllers

import (
	"LTE_Security/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

type RegisterData struct {
	Mac  string `json:"mac"`
	Mid  string `json:"mid"`
	Psn  string `json:"psn"`
	Guid string `json:"guid"`
	Code int64  `json:"code"`
}

type RegisterController struct {
	beego.Controller
}

func (this *RegisterController) Get() {
	//check islogin
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

	this.TplNames = "login.html"
}

/* ret
(0)deivce register success
(-1)uname and pwd not match
(-2)invaild mac address
(-3)repeated registration
(-4)device register over time
(-5)insert db failed
*/

func (this *RegisterController) Post() {
	data := RegisterData{}
	//get client data
	data.Mac = this.GetString("mac")
	data.Mid = this.GetString("mid")
	data.Psn = this.GetString("psn")
	//check auth
	uname, pwd, ok := this.Ctx.Request.BasicAuth()
	if !ok {
		beego.Info("get client  Request.BasicAuth failed!")
		data.Code = -1
		writeContent, _ := json.Marshal(data)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	user := models.Userinfo{
		Username: uname,
		Password: pwd,
	}
	ok = models.CheckAccount(&user)
	if !ok {
		beego.Info("user/pwd not matched!")
		data.Code = -1
		writeContent, _ := json.Marshal(data)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	//register
	deviceinfo := models.Deviceinfo{
		Mac: data.Mac,
		Mid: data.Mid,
		Psn: data.Psn,
	}
	data.Guid, data.Code = models.RegisterDeivce(&deviceinfo)
	writeContent, _ := json.Marshal(data)
	this.Ctx.WriteString(string(writeContent))
}
