package controllers

import (
	"LTE_Security/models"
	"github.com/astaxie/beego"
	"strconv"
)

var (
	NR_devicesCount int64 = 0  //default device numbers in database
	NR_totalPages   int64 = 0  //default page numbers to show
	NR_curPage      int64 = 1  //default current page
	NR_listCount    int64 = 10 //show numbers of deviceinfo every page
	NR_curCount     int64 = 0  //total numbers of deviceinfo had been showed
)

type NoregisterController struct {
	beego.Controller
}

func (this *NoregisterController) Get() {
	//check islogin
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

	this.TplNames = "noregister.html"

	//recive request listcount info
	listcount := this.Input().Get("listcount")
	beego.Debug("listcount=", listcount)
	if !(listcount == "") {
		count, _ := strconv.Atoi(listcount)
		switch count {
		case 20:
			NR_listCount = 20
			NR_curCount = 0
			NR_curPage = 1
		case 30:
			NR_listCount = 30
			NR_curCount = 0
			NR_curPage = 1
		default:
			NR_listCount = 10
			NR_curCount = 0
			NR_curPage = 1
		}
	}

	//calc variables to show
	NR_devicesCount = models.GetNoregisterDevicesCount()
	beego.Debug("devicesCount=", NR_devicesCount)
	if NR_devicesCount%NR_listCount > 0 {
		NR_totalPages = NR_devicesCount/NR_listCount + 1
	} else {
		NR_totalPages = NR_devicesCount / listCount
	}

	//recive request op info
	ope := this.Input().Get("op")
	beego.Debug("ope=", ope)
	switch ope {
	case "firstpage":
		NR_curCount = 0
		NR_curPage = 1
	case "prepage":
		if NR_curPage > 1 {
			NR_curCount -= listCount
			NR_curPage -= 1
		}
	case "nextpage":
		if NR_curPage < NR_totalPages {
			NR_curCount += listCount
			NR_curPage += 1
		}
	case "lastpage":
		NR_curCount = NR_listCount * (NR_totalPages - 1)
		NR_curPage = NR_totalPages
	}

	devices, _, ok := models.GetNoregisterDevices(NR_listCount, NR_curCount)
	if ok {
		beego.Info("GetDevices success!")
		this.Data["Devices"] = devices
		this.Data["DevicesNum"] = NR_devicesCount
		this.Data["CurPage"] = curPage
		this.Data["TotalPages"] = NR_totalPages
	} else {
		this.Data["NoInfo"] = "There have no devices!"
	}
}

func (this *NoregisterController) Post() {
	//check islogin
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

	this.TplNames = "login.html"
}
