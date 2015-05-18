package controllers

import (
	"LTE_Security/models"
	"github.com/astaxie/beego"
	"github.com/tealeg/xlsx"
	"os"
	"regexp"
	"strconv"
	"time"
)

const (
	regular = "^([0-9a-fA-F]{2})(([:][0-9a-fA-F]{2}){5})$"
)

type ImportController struct {
	beego.Controller
}

func (this *ImportController) Get() {
	//check islogin
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.Data["IsImported"] = false
	this.TplNames = "import.html"
}

func (this *ImportController) Post() {
	//check islogin
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

	saveFileName := "temp" + strconv.FormatInt(int64(time.Now().UnixNano()), 10) + ".xlsx"
	beego.Debug("saveFileName=", saveFileName)

	this.SaveToFile("file", saveFileName)

	devices := make([]models.Deviceinfo, 0)

	xlFile, err := xlsx.OpenFile(saveFileName)
	if err != nil {
		beego.Error("Open excel file!", err)
	}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				device := models.Deviceinfo{
					Mac:         cell.String(),
					ImportTime:  time.Now(),
					InvalidTime: time.Now().Add(3 * 30 * 24 * time.Hour),
				}
				//check mac address format
				if !CheckMacFormat(device.Mac) {
					beego.Info(device.Mac, "not a mac address!")
					continue
				}
				//if the mac address had existed,skip it
				if models.ImportDeviceCheck(&device) {
					beego.Info(device.Mac, "had been imported!")
					continue
				}
				devices = append(devices, device)
			}
		}
	}
	//delete used file
	err = os.Remove(saveFileName)
	if err != nil {
		beego.Error("Remove temp excel file failed!", err)
	}

	ok := models.ImportDevices(&devices)
	if ok {
		beego.Info("ImportDevices success!")
		this.Redirect("/home", 301)
		this.Data["IsImported"] = true
		return
	} else {
		beego.Info("ImportDevices failed! once again")
		this.Redirect("/import", 302)
		return
	}
}

func CheckMacFormat(mac string) bool {
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mac)
}
