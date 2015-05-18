package models

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"time"
)

type Userinfo struct {
	Uid           int64     `orm:"auto"`
	Username      string    `json:"uname"`
	Password      string    `json:"pwd"`
	CreatedTime   time.Time `json:"-"`
	LastLoginTime time.Time `orm:"null;auto_now;type(datetime)";json:"lastlogin"`
}

type Admininfo struct {
	Uid           int64     `orm:"auto"`
	Username      string    `json:"uname"`
	Password      string    `json:"pwd"`
	CreatedTime   time.Time `json:"-"`
	LastLoginTime time.Time `orm:"null;auto_now;type(datetime)";json:"lastlogin"`
}

type Deviceinfo struct {
	//Id               int64     `orm:"auto;pk"`
	Mac              string    `orm:"pk";json:"mac"`
	Guid             string    `orm:"null";json:"guid"`
	Mid              string    `orm:"null";json:"mid"`
	Psn              string    `orm:"null";json:"psn"`
	ImportTime       time.Time `json:"-"`
	InvalidTime      time.Time `json:"-"`
	RegistrationTime time.Time `orm:"null";json:"-"`
	UpdateTime       time.Time `orm:"null";json:"-"`
}

func RegisterDB() {
	//register all tables
	orm.RegisterModel(new(Deviceinfo), new(Userinfo), new(Admininfo))
	//register mysql driver
	err := orm.RegisterDriver("mysql", orm.DR_MySQL)
	if err != nil {
		beego.Critical(err)
	}
	//register default database lte_security
	err = orm.RegisterDataBase("default", "mysql", "root:way@tcp(192.168.15.155:3306)/lte_security?charset=utf8&loc=Asia%2FShanghai")
	//orm.RegisterDataBase("default", "mysql", "root:@/lte_security?charset=utf8&loc=Asia%2FShanghai")
	if err != nil {
		beego.Critical(err)
	}
}

func CheckAdmin(admin *Admininfo) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("admininfo").Filter("UserName", admin.Username).Filter("Password", admin.Password).Exist()
	return exist
}

func CheckAccount(user *Userinfo) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("userinfo").Filter("UserName", user.Username).Filter("Password", user.Password).Exist()
	return exist
}

func UpdateAdminStatus(admin *Admininfo) bool {
	o := orm.NewOrm()
	var u Admininfo
	err := o.QueryTable("admininfo").Filter("UserName", admin.Username).One(&u)
	if err != nil {
		beego.Error(err)
		return false
	}
	admin.Uid = u.Uid
	admin.CreatedTime = u.CreatedTime
	admin.LastLoginTime = time.Now()
	_, err = o.Update(admin)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func ImportDeviceCheck(device *Deviceinfo) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("deviceinfo").Filter("mac", device.Mac).Exist()
	return exist
}

func ImportDevices(devices *[]Deviceinfo) bool {
	o := orm.NewOrm()
	//successNum, err := o.InsertMulti(100, devices)
	_, err := o.InsertMulti(100, devices)
	if err != nil {
		beego.Error(err)
		return false
	}
	//fmt.Println("successNum=", successNum)
	return true
}

func RegisterDeivce(deviceinfo *Deviceinfo) (string, int64) {
	o := orm.NewOrm()
	guid := GetGuid()
	var device Deviceinfo
	err := o.QueryTable("deviceinfo").Filter("mac", deviceinfo.Mac).One(&device)
	if err != nil {
		//device doesnot  exsited, return error info
		beego.Error(err)
		return "", -2
	} else {
		//check if guid exist
		exist := o.QueryTable("deviceinfo").Filter("mac", device.Mac).Filter("registration_time__isnull", false).Exist()
		if exist {
			return "", -3
		}
		//check device register time is over time
		if device.InvalidTime.Before(time.Now()) {
			return "", -4
		}
		//device had exsited,  insert info
		device.Guid = guid
		device.Mid = deviceinfo.Mid
		device.Psn = deviceinfo.Psn
		device.RegistrationTime = time.Now()
		_, err := o.Update(&device)
		if err != nil {
			beego.Error(err)
			return "", -5
		}
	}
	return guid, 0
}

func VerifyDevice(deviceinfo *Deviceinfo) int64 {
	o := orm.NewOrm()
	exist := o.QueryTable("deviceinfo").Filter("mac", deviceinfo.Mac).Filter("guid", deviceinfo.Guid).Filter("mid", deviceinfo.Mid).Filter("psn", deviceinfo.Psn).Exist()
	if !exist {
		return -3 //mac/guid/mid/psn not match
	}
	return 0 //verify success
}

func ModifyDevice(deviceinfo *Deviceinfo, newmac string) bool {
	o := orm.NewOrm()
	var device Deviceinfo
	err := o.QueryTable("deviceinfo").Filter("mac", deviceinfo.Mac).One(&device)
	if err != nil {
		beego.Error(err)
		return false
	}
	//delete old record
	_, err = o.Delete(&device)
	if err != nil {
		beego.Error(err)
		return false
	}
	//insert new mac record
	device.Mac = newmac
	device.Guid = ""
	device.UpdateTime = time.Now()
	_, err = o.Insert(&device)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func GetDevices(start, offset int64) ([]*Deviceinfo, int64, bool) {
	o := orm.NewOrm()
	//get all devices
	devices := make([]*Deviceinfo, 0)
	num, err := o.QueryTable("deviceinfo").Limit(start, offset).Filter("registration_time__isnull", false).All(&devices)
	if err != nil {
		beego.Error(err)
		return nil, 0, false
	}
	return devices, num, true
}

func GetNoregisterDevices(start, offset int64) ([]*Deviceinfo, int64, bool) {
	o := orm.NewOrm()
	//get all devices of noregister
	devices := make([]*Deviceinfo, 0)
	num, err := o.QueryTable("deviceinfo").Limit(start, offset).Filter("registration_time__isnull", true).All(&devices)
	if err != nil {
		beego.Error(err)
		return nil, 0, false
	}
	return devices, num, true
}

func GetDevicesCount() int64 {
	o := orm.NewOrm()
	cnt, _ := o.QueryTable("deviceinfo").Filter("registration_time__isnull", false).Count()
	return cnt
}

func GetNoregisterDevicesCount() int64 {
	o := orm.NewOrm()
	cnt, _ := o.QueryTable("deviceinfo").Filter("registration_time__isnull", true).Count()
	return cnt
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetGuid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		beego.Error(err)
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}
