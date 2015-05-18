package main

import (
	"LTE_Security/models"
	_ "LTE_Security/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/session"
)

var GlobalSessions *session.Manager

func init() {
	models.RegisterDB()
	GlobalSessions, _ = session.NewManager("memory", `{
												"cookieName":"lte_securityid",
												"enableSetCookie,omitempty": true,
												"gclifetime":3600,
												"maxLifetime": 3600,
												"secure": true,
												"sessionIDHashFunc": "sha1",
												"sessionIDHashKey": "",
												"cookieLifeTime": 3600,
												"providerConfig": ""}`)
	go GlobalSessions.GC()
	beego.SetLogger("file", `{"filename":"logs/server.log"}`)
	beego.SetLevel(beego.LevelDebug)
}

func main() {
	//Open orm debug mode
	orm.Debug = false

	//Auto create tables
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		beego.Critical("sycndb error! Error:", err)
	}

	//Start app
	beego.Trace("LTE_Security start running...")
	beego.Run()
}
