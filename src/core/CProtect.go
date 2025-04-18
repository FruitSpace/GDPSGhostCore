package core

import (
	"encoding/json"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type CProtect struct {
	DB                *MySQLConn
	LevelModel        ProtectModel
	Savepath          string
	DisableProtection bool
}

type ProtectModel struct {
	MaxStars        int
	MaxLevelUpload  int
	PeakLevelUpload int
	Stats           map[string]int
}

func (protect *CProtect) LoadModel(config *GlobalConfig, blob ConfigBlob) {
	model, err := os.ReadFile(protect.Savepath + "/levelModel.json")
	if err != nil {
		os.Mkdir(protect.Savepath, 0777)
		protect.FillLevelModel()
		//req, err := http.Get(config.ApiEndpoint + "?id=" + blob.ServerConfig.SrvID + "&key=" + blob.ServerConfig.SrvKey + "&action=getModel")
		//data, err := io.ReadAll(req.Body)
		//if err != nil || req.StatusCode != 200 {
		//	protect.DisableProtection = true
		//	return
		//}
		//os.WriteFile(protect.Savepath+"/levelModel.json", data, 0755)
		//protect.LoadModel(config, blob)
	}
	json.Unmarshal(model, &protect.LevelModel)
	if protect.LevelModel.MaxLevelUpload == 0 {
		protect.LevelModel.MaxLevelUpload = 10
	}
}

func (protect *CProtect) FillLevelModel() {

	//Calculate LevelModel
	date := time.Now()
	stats := make(map[string]int)
	total := 0
	for i := 0; i < 7; i++ {
		current := strings.Split(date.AddDate(0, 0, -1*i).Format("2006-01-02 15:04:05"), " ")[0]
		currentIndex := strings.Split(date.AddDate(0, 0, -1*(i+1)).Format("2006-01-02 15:04:05"), " ")[0]
		var count int
		protect.DB.ShouldQueryRow("SELECT count(*) as cnt FROM #DB#.actions WHERE type=4 AND date<? AND date>? AND data LIKE '%Upload%'",
			current, currentIndex).Scan(&count)
		stats[currentIndex] = count
		if count > protect.LevelModel.PeakLevelUpload {
			protect.LevelModel.PeakLevelUpload = count
		}
		total += count
	}
	if total < 10 {
		protect.LevelModel.MaxLevelUpload = 10
	} else {
		protect.LevelModel.MaxLevelUpload = int(math.Round(float64(total/7))) + protect.LevelModel.PeakLevelUpload
	}

	//Calculate total stars allowed
	var count2, count1 int
	protect.DB.ShouldQueryRow("SELECT SUM(starsGot) as stars FROM #DB#.levels").Scan(&count1)
	protect.DB.ShouldQueryRow("SELECT SUM(packStars) as stars FROM #DB#.levelpacks").Scan(&count2)
	protect.LevelModel.MaxStars = 200 + count1 + count2

	protect.LevelModel.Stats = stats
	//Dump
	data, err := json.Marshal(protect.LevelModel)
	if err != nil {
		data = []byte("{}")
	}
	protect.DB.logger.Must(os.WriteFile(protect.Savepath+"/levelModel.json", data, 0755))
}

func (protect *CProtect) ResetUserLimits() {
	protect.DB.ShouldExec("UPDATE #DB#.users SET protect_levelsToday=0")
	protect.DB.ShouldExec("UPDATE #DB#.users SET protect_todayStars=stars")
}

func (protect *CProtect) DetectLevelModel(uid int) bool {
	//FIXME
	if 1 == 1 {
		return true
	}
	if protect.DisableProtection {
		return true
	}
	var lvlCnt int
	protect.DB.ShouldQueryRow("SELECT protect_levelsToday as cnt FROM #DB#.users WHERE uid=?", uid).Scan(&lvlCnt)
	if lvlCnt >= protect.LevelModel.MaxLevelUpload {
		protect.DB.ShouldExec("UPDATE #DB#.users SET isBanned=2 WHERE uid=?", uid)
		RegisterAction(ACTION_BAN_BAN, 0, uid, map[string]string{"type": "Ban:LevelAuto"}, protect.DB)
		SendMessageDiscord("[" + protect.Savepath + "] User " + strconv.Itoa(uid) + " has been banned for uploading too many levels (" + strconv.Itoa(lvlCnt) + "/" + strconv.Itoa(protect.LevelModel.MaxLevelUpload) + ") in a day.")
		return false
	}
	protect.DB.ShouldExec("UPDATE #DB#.users SET protect_levelsToday=protect_levelsToday+1 WHERE uid=?", uid)
	return true
}

func (protect *CProtect) DetectStats(uid int, stars int, diamonds int, demons int, coins int, ucoins int) bool {
	if protect.DisableProtection {
		return true
	}
	if stars < 0 || diamonds < 0 || demons < 0 || coins < 0 || ucoins < 0 {
		protect.DB.ShouldExec("UPDATE #DB#.users SET isBanned=2 WHERE uid=?", uid)
		protect.DB.ShouldExec("DELETE FROM #DB#.levels WHERE uid=?", uid)
		protect.DB.ShouldExec("DELETE FROM #DB#.actions WHERE type=4 AND uid=?", uid)
		RegisterAction(ACTION_BAN_BAN, 0, uid, map[string]string{"type": "Ban:StatsNegative"}, protect.DB)
		SendMessageDiscord("User " + strconv.Itoa(uid) + " has been banned for having negative stats.")
		return false
	}
	//FIXME
	if 1 == 1 {
		return true
	}
	if protect.LevelModel.MaxStars == 0 {
		protect.LevelModel.MaxStars = 200
	}
	var starCnt int
	protect.DB.ShouldQueryRow("SELECT protect_todayStars FROM #DB#.users WHERE uid=?", uid).Scan(&starCnt)
	if (stars - starCnt) > protect.LevelModel.MaxStars {
		protect.DB.ShouldExec("UPDATE #DB#.users SET isBanned=2 WHERE uid=?", uid)
		RegisterAction(ACTION_BAN_BAN, 0, uid, map[string]string{"type": "Ban:StarsLimit"}, protect.DB)
		SendMessageDiscord("User " + strconv.Itoa(uid) + " has been banned for having too many stars (" + strconv.Itoa(stars) + "+" + strconv.Itoa(starCnt) + "/" + strconv.Itoa(protect.LevelModel.MaxStars) + ").")
		return false
	}
	return true
}

func (protect *CProtect) GetMeta(uid int) map[string]int {
	meta := make(map[string]int)
	var sMeta string
	protect.DB.ShouldQueryRow("SELECT protect_meta FROM #DB#.users WHERE uid=?", uid).Scan(&sMeta)
	json.Unmarshal([]byte(sMeta), &meta)
	return meta
}

func (protect *CProtect) DetectMessages(uid int) bool {
	if protect.DisableProtection {
		return true
	}
	meta := protect.GetMeta(uid)
	t := int(time.Now().Unix())
	if t-meta["msg_time"] < 60 {
		return false
	}
	meta["msg_time"] = t
	data, _ := json.Marshal(meta)
	protect.DB.ShouldExec("UPDATE #DB#.users SET protect_meta=? WHERE uid=?", string(data), uid)
	return true
}

func (protect *CProtect) DetectPosts(uid int) bool {
	if protect.DisableProtection {
		return true
	}
	meta := protect.GetMeta(uid)
	t := int(time.Now().Unix())
	if t-meta["post_time"] < 60 {
		return false
	}
	meta["post_time"] = t
	data, _ := json.Marshal(meta)
	protect.DB.ShouldExec("UPDATE #DB#.users SET protect_meta=? WHERE uid=?", string(data), uid)
	return true
}

func (protect *CProtect) DetectComments(uid int) bool {
	if protect.DisableProtection {
		return true
	}
	meta := protect.GetMeta(uid)
	t := int(time.Now().Unix())
	if t-meta["comm_time"] < 30 {
		return false
	}
	meta["comm_time"] = t
	data, _ := json.Marshal(meta)
	protect.DB.ShouldExec("UPDATE #DB#.users SET protect_meta=? WHERE uid=?", string(data), uid)
	return true
}

func CheckIPBan(IPAddr string, config ConfigBlob) bool {
	return InArray(config.SecurityConfig.BannedIPs, IPAddr)
}
