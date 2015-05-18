package controllers

import (
	"LTE_Security/models"
	"github.com/astaxie/beego"
	"strconv"
)

var (
	devicesCount int64 = 0  //default device numbers in database
	totalPages   int64 = 0  //default page numbers to show
	curPage      int64 = 1  //default current page
	listCount    int64 = 10 //show numbers of deviceinfo every page
	curCount     int64 = 0  //total numbers of deviceinfo had been showed
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	//check islogin
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.TplNames = "home.html"

	//recive request listcount info
	listcount := this.Input().Get("listcount")
	if !(listcount == "") {
		count, _ := strconv.Atoi(listcount)
		switch count {
		case 20:
			listCount = 20
			curCount = 0
			curPage = 1
		case 30:
			listCount = 30
			curCount = 0
			curPage = 1
		default:
			listCount = 10
			curCount = 0
			curPage = 1
		}
	}
	beego.Debug("listcount=", listCount)

	//calc variables to show
	devicesCount = models.GetDevicesCount()
	beego.Debug("devicesCount=", devicesCount)
	if devicesCount%listCount > 0 {
		totalPages = devicesCount/listCount + 1
	} else {
		totalPages = devicesCount / listCount
	}

	//recive request op info
	ope := this.Input().Get("op")
	beego.Debug("ope=", ope)
	switch ope {
	case "firstpage":
		curCount = 0
		curPage = 1
	case "prepage":
		if curPage > 1 {
			curCount -= listCount
			curPage -= 1
		}
	case "nextpage":
		if curPage < totalPages {
			curCount += listCount
			curPage += 1
		}
	case "lastpage":
		curCount = listCount * (totalPages - 1)
		curPage = totalPages
	}

	devices, _, ok := models.GetDevices(listCount, curCount)
	if ok {
		beego.Info("GetDevices success!")
		this.Data["Devices"] = devices
		this.Data["DevicesNum"] = devicesCount
		this.Data["CurPage"] = curPage
		this.Data["TotalPages"] = totalPages
	} else {
		this.Data["NoInfo"] = "There have no devices!"
	}
}

func (this *HomeController) Post() {
	//check islogin
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

	this.TplNames = "login.html"
}
