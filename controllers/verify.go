package controllers

import (
	"LTE_Security/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

type VerifyData struct {
	Code int64 `json:"code"`
}

type VerifyController struct {
	beego.Controller
}

func (this *VerifyController) Get() {
	//check islogin
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.TplNames = "login.html"
}

/* return value
(0)verify success
(-1)uname and pwd not match
(-2)input content error
(-3)verify failed
*/
func (this *VerifyController) Post() {
	data := VerifyData{}
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
	//get client data
	mac := this.GetString("mac")
	guid := this.GetString("guid")
	mid := this.GetString("mid")
	psn := this.GetString("psn")

	if mac == "" || guid == "" || mid == "" || psn == "" {
		beego.Info(" input content error!")
		data.Code = -2
		writeContent, _ := json.Marshal(data)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	//verify from deviceinfo table
	deviceinfo := &models.Deviceinfo{
		Mac:  mac,
		Mid:  mid,
		Psn:  psn,
		Guid: guid,
	}

	data.Code = models.VerifyDevice(deviceinfo)
	writeContent, _ := json.Marshal(data)
	this.Ctx.WriteString(string(writeContent))
	return
}
