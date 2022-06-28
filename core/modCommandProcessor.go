package core

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

func InvokeCommands(db MySQLConn, cl CLevel, acc CAccount, comment string) bool {
	command:=strings.Split(comment," ")
	isOwner:=cl.IsOwnedBy(acc.Uid)
	role:=acc.GetRoleObj(true)
	switch command[0] {
	case "!feature":
		if role.RoleName=="" || role.Privs["cFeature"]!=1 {return false}
		cl.FeatureLevel(true)
		RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Feature"},db)
		return true
	case "!unfeature":
		if role.RoleName=="" || role.Privs["cFeature"]!=1 {return false}
		cl.FeatureLevel(false)
		RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Uneature"},db)
		return true
	case "!epic":
		if role.RoleName=="" || role.Privs["cEpic"]!=1 {return false}
		cl.FeatureLevel(true)
		RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Epic"},db)
		return true
	case "!unepic":
		if role.RoleName=="" || role.Privs["cEpic"]!=1 {return false}
		cl.FeatureLevel(false)
		RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Unepic"},db)
		return true
	case "!coins":
		if role.RoleName=="" || role.Privs["cVerCoins"]!=1 {return false}
		if len(command)<2 {return false}
		if command[1]=="verify" {
			cl.VerifyCoins(true);
			RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Coins:Verify"},db)
		}else if command[1]=="reset" {
			cl.VerifyCoins(false);
			RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Coins:Reset"},db)
		} else{return false}
		return true
	case "!daily":
		if role.RoleName=="" || role.Privs["cDaily"]!=1 {return false}
		if len(command)>1 && command[1]=="reset" {
			db.DB.Query("DELETE FROM quests WHERE lvl_id=? and type=0",cl.Id)
			RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Daily:Reset"},db)
		}else{
			var date string
			err:=db.DB.QueryRow("SELECT timeExpire FROM quests WHERE type=0 ORDER BY timeExpire DESC LIMIT 1").Scan(&date)
			if err==sql.ErrNoRows {
				date=strings.Split(time.Now().Format("2006-01-02 15:04:05")," ")[0]+" 00:00:00"
			}else{
				tme,_:=time.Parse("2006-01-02 00:00:00",strings.Split(date," ")[0]+" 00:00:00")
				tme.AddDate(0,0,1)
				date=tme.Format("2006-01-02 15:04:05")
			}
			db.DB.Query("INSERT INTO quests (type,lvl_id,timeExpire) VALUES (0,?,?)",cl.Id,date)
			RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Daily:Publish"},db)
		}
		return true
	case "!weekly":
		if role.RoleName=="" || role.Privs["cWeekly"]!=1 {return false}
		if len(command)>1 && command[1]=="reset" {
			db.DB.Query("DELETE FROM quests WHERE lvl_id=? and type=1",cl.Id)
			RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Weekly:Reset"},db)
		}else{
			var date string
			err:=db.DB.QueryRow("SELECT timeExpire FROM quests WHERE type=0 ORDER BY timeExpire DESC LIMIT 1").Scan(&date)
			if err==sql.ErrNoRows {
				date=strings.Split(time.Now().Format("2006-01-02 15:04:05")," ")[0]+" 00:00:00"
			}else{
				tme,_:=time.Parse("2006-01-02 00:00:00",strings.Split(date," ")[0]+" 00:00:00")
				tme.AddDate(0,0,7)
				date=tme.Format("2006-01-02 15:04:05")
			}
			db.DB.Query("INSERT INTO quests (type,lvl_id,timeExpire) VALUES (1,?,?)",cl.Id,date)
			RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Weekly:Publish"},db)
		}
		return true
	case "!rate":
		if role.RoleName=="" || role.Privs["cRate"]!=1 {return false}
		if len(command)<2 {return false}
		diff:="0"
		switch(strings.ToLower(command[1])){
		case "auto":
			diff="-1"
			break
		case "easy":
			diff="10"
			break
		case "normal";
			diff="20"
			break
		case "hard":
			diff="30"
			break
		case "harder":
			diff="40"
			break
		case "insane":
			diff="50"
			break
		case "reset":
			diff="0,starsGot=0"
			break
		default:
			return false
		}
		db.DB.Query("UPDATE levels SET difficulty="+diff+" WHERE id=?",cl.Id)
		RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Rate:"+strings.Title(strings.ToLower(command[1]))},db)
		return true
	case "!lvl":
		if len(command)<2 {return false}
		m:="Mod"
		if isOwner{m="Owner"}
		switch command[1] {
		case "delete":
			if role.RoleName=="" || role.Privs["cDelete"]!=1 {return false}
			if len(command)<3 || command[2]!=strconv.Itoa(cl.Id) {return false}
			cl.DeleteLevel()
			RegisterAction(ACTION_LEVEL_DELETE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Delete:"+m},db)
			return true
		case "rename":
			if !isOwner && (role.RoleName=="" || role.Privs["cLvlAccess"]!=1) {return false}
			if len(command)<3 {return false}
			text:=strings.Replace(comment,"!lvl rename ","")
			db.DB.Query()
			RegisterAction(ACTION_LEVEL_UPDATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Rename:"+m},db)
		}
	}
}
