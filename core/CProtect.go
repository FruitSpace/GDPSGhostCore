package core

import (
	"encoding/json"
	"math"
	"os"
	"strings"
	"time"
)

type CProtect struct {
	DB MySQLConn
	LevelModel ProtectModel
	Savepath string
	DisableProtection bool
}

type ProtectModel struct {
	MaxStars int
	MaxLevelUpload int
	PeakLevelUpload int
	Stats map[string]int
}


func (protect *CProtect) LoadModel() {
	model, err:= os.ReadFile(protect.Savepath+"/levelModel.json")
	if err!=nil{
		os.Mkdir(protect.Savepath,0777)
		os.WriteFile(protect.Savepath+"/levelModel.json",[]byte("{}"),0755)
		return
	}
	json.Unmarshal(model,&protect.LevelModel)
}

func (protect *CProtect) FillLevelModel() {

	//Calculate LevelModel
	date:=time.Now()
	stats:=make(map[string]int)
	total:=0
	for i:=0; i<7; i++ {
		current:=strings.Split(date.AddDate(0,0,-1*i).Format("2006-01-02 15:04:05")," ")[0]
		currentIndex:=strings.Split(date.AddDate(0,0,-(i+1)*i).Format("2006-01-02 15:04:05")," ")[0]
		var count int
		protect.DB.ShouldQueryRow("SELECT count(*) as cnt FROM actions WHERE type=4 AND date<? AND date>? AND data LIKE '%Upload%'",
			current,currentIndex).Scan(&count)
		stats[currentIndex]=count
		if count>protect.LevelModel.PeakLevelUpload {
			protect.LevelModel.PeakLevelUpload=count
		}
		total+=count
	}
	if total<10 {
		protect.LevelModel.MaxLevelUpload=10
	}else{
		protect.LevelModel.MaxLevelUpload=int(math.Round(float64(total/7)))+protect.LevelModel.PeakLevelUpload
	}

	//Calculate total stars allowed
	var count2, count1 int
	protect.DB.ShouldQueryRow("SELECT SUM(starsGot) as stars FROM levels").Scan(&count1)
	protect.DB.ShouldQueryRow("SELECT SUM(packStars) as stars FROM levelpacks").Scan(&count2)
	protect.LevelModel.MaxStars=200+count1+count2

	//Dump
	data,err:=json.Marshal(protect.LevelModel)
	if err!=nil {
		data=[]byte("{}")
	}
	protect.DB.logger.Must(os.WriteFile(protect.Savepath+"/levelModel.json",data,0755))
}

func (protect *CProtect) ResetUserLimits() {
	protect.DB.ShouldQuery("UPDATE users SET protect_levelsToday=0")
	protect.DB.ShouldQuery("UPDATE users SET protect_todayStars=stars")
}

func (protect *CProtect) DetectLevelModel(uid int) bool {
	if protect.DisableProtection {return true}
	var lvlCnt int
	protect.DB.ShouldQueryRow("SELECT protect_levelsToday as cnt FROM users WHERE uid=?",uid).Scan(&lvlCnt)
	if lvlCnt>=protect.LevelModel.MaxLevelUpload {
		protect.DB.ShouldQuery("UPDATE users SET isBanned=2 WHERE uid=?",uid)
		RegisterAction(ACTION_BAN_BAN,0,uid, map[string]string{"type":"Ban:LevelAuto"},protect.DB)
		return false
	}
	protect.DB.ShouldQuery("UPDATE users SET protect_levelsToday=protect_levelsToday+1 WHERE uid=?",uid)
	return true
}

func (protect *CProtect) DetectStats(uid int, stars int, diamonds int, demons int, coins int, ucoins int) bool {
	if protect.DisableProtection {return true}
	if stars<0 || diamonds<0 || demons<0 || coins<0 || ucoins<0 {
		protect.DB.ShouldQuery("UPDATE users SET isBanned=2 WHERE uid=?",uid)
		protect.DB.ShouldQuery("DELETE FROM levels WHERE uid=?",uid)
		protect.DB.ShouldQuery("DELETE FROM actions WHERE type=4 AND uid=?",uid)
		RegisterAction(ACTION_BAN_BAN,0,uid, map[string]string{"type":"Ban:StatsNegative"},protect.DB)
		return false
	}
	var starCnt int
	protect.DB.ShouldQuery("SELECT protect_todayStars as cnt FROM users WHERE uid=?",uid).Scan(&starCnt)
	if stars-starCnt>protect.LevelModel.MaxStars {
		protect.DB.ShouldQuery("UPDATE users SET isBanned=2 WHERE uid=?",uid)
		RegisterAction(ACTION_BAN_BAN,0,uid, map[string]string{"type":"Ban:StarsLimit"},protect.DB)
		return false
	}
	return true
}

func (protect *CProtect) GetMeta(uid int) map[string]int{
	meta:=make(map[string]int)
	var sMeta string
	protect.DB.ShouldQueryRow("SELECT protect_meta FROM users WHERE uid=?",uid).Scan(&sMeta)
	json.Unmarshal([]byte(sMeta),&meta)
	return meta
}

func (protect *CProtect) DetectMessages(uid int) bool {
	if protect.DisableProtection {return true}
	meta:=protect.GetMeta(uid)
	t :=int(time.Now().Unix())
	if t-meta["msg_time"]<120 {return false}
	meta["msg_time"]= t
	data,_:=json.Marshal(meta)
	protect.DB.ShouldQuery("UPDATE users SET protect_meta=? WHERE uid=?",string(data), uid)
	return true
}

func (protect *CProtect) DetectPosts(uid int) bool {
	if protect.DisableProtection {return true}
	meta:=protect.GetMeta(uid)
	t :=int(time.Now().Unix())
	if t-meta["post_time"]<900 {return false}
	meta["post_time"]= t
	data,_:=json.Marshal(meta)
	protect.DB.ShouldQuery("UPDATE users SET protect_meta=? WHERE uid=?",string(data), uid)
	return true
}

func (protect *CProtect) DetectComments(uid int) bool {
	if protect.DisableProtection {return true}
	meta:=protect.GetMeta(uid)
	t :=int(time.Now().Unix())
	if t-meta["comm_time"]<120 {return false}
	meta["comm_time"]= t
	data,_:=json.Marshal(meta)
	protect.DB.ShouldQuery("UPDATE users SET protect_meta=? WHERE uid=?",string(data), uid)
	return true
}


func CheckIPBan(IPAddr string, config ConfigBlob) bool {
	return InArray(config.SecurityConfig.BannedIPs, IPAddr)
}