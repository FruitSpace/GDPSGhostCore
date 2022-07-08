package core

import (
	"database/sql"
	"encoding/base64"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
	"time"
)

func InvokeCommands(db MySQLConn, cl CLevel, acc CAccount, comment string, isOwner bool, role Role) bool {
	command:=strings.Split(comment," ")
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
			cl.VerifyCoins(true)
			RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Coins:Verify"},db)
		}else if command[1]=="reset" {
			cl.VerifyCoins(false)
			RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Coins:Reset"},db)
		} else{return false}
		return true
	case "!daily":
		if role.RoleName=="" || role.Privs["cDaily"]!=1 {return false}
		if len(command)>1 && command[1]=="reset" {
			db.ShouldQuery("DELETE FROM quests WHERE lvl_id=? and type=0",cl.Id)
			RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Daily:Reset"},db)
		}else{
			var date string
			err:=db.DB.QueryRow("SELECT timeExpire FROM quests WHERE type=0 ORDER BY timeExpire DESC LIMIT 1").Scan(&date)
			if err==sql.ErrNoRows {
				date=strings.Split(time.Now().Format("2006-01-02 15:04:05")," ")[0]+" 00:00:00"
			}else{
				tme,_:=time.Parse("2006-01-02 15:04:05",strings.Split(date," ")[0]+" 00:00:00")
				tme.AddDate(0,0,1)
				date=tme.Format("2006-01-02 15:04:05")
			}
			db.ShouldQuery("INSERT INTO quests (type,lvl_id,timeExpire) VALUES (0,?,?)",cl.Id,date)
			RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Daily:Publish"},db)
		}
		return true
	case "!weekly":
		if role.RoleName=="" || role.Privs["cWeekly"]!=1 {return false}
		if len(command)>1 && command[1]=="reset" {
			db.ShouldQuery("DELETE FROM quests WHERE lvl_id=? and type=1",cl.Id)
			RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Weekly:Reset"},db)
		}else{
			var date string
			err:=db.DB.QueryRow("SELECT timeExpire FROM quests WHERE type=0 ORDER BY timeExpire DESC LIMIT 1").Scan(&date)
			if err==sql.ErrNoRows {
				date=strings.Split(time.Now().Format("2006-01-02 15:04:05")," ")[0]+" 00:00:00"
			}else{
				tme,_:=time.Parse("2006-01-02 15:04:05",strings.Split(date," ")[0]+" 00:00:00")
				tme.AddDate(0,0,7)
				date=tme.Format("2006-01-02 15:04:05")
			}
			db.ShouldQuery("INSERT INTO quests (type,lvl_id,timeExpire) VALUES (1,?,?)",cl.Id,date)
			RegisterAction(ACTION_LEVEL_RATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Weekly:Publish"},db)
		}
		return true
	case "!rate":
		if role.RoleName=="" || role.Privs["cRate"]!=1 {return false}
		if len(command)<2 {return false}
		diff:="0"
		switch strings.ToLower(command[1]) {
		case "auto":
			diff="-1"
			break
		case "easy":
			diff="10"
			break
		case "normal":
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
		db.ShouldQuery("UPDATE levels SET difficulty="+diff+" WHERE id=?",cl.Id)
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
			text:=strings.Replace(comment,"!lvl rename ","",1)
			db.ShouldQuery("UPDATE levels SET name=? WHERE id=?",text,cl.Id)
			RegisterAction(ACTION_LEVEL_UPDATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Rename:"+m},db)
			return true
		case "copy":
			if !isOwner {return false}
			if len(command)<3 {return false}
			switch command[2] {
			case "on":
				db.ShouldQuery("UPDATE levels SET password=1 WHERE id=?",cl.Id)
				RegisterAction(ACTION_LEVEL_UPDATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Copy:Enable"},db)
				return true
			case "off":
				db.ShouldQuery("UPDATE levels SET password=0 WHERE id=?",cl.Id)
				RegisterAction(ACTION_LEVEL_UPDATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Copy:Disable"},db)
				return true
			case "pass":
				if len(command)<4 || len(command[3])!=6 {return false}
				if c,err:=strconv.Atoi(command[3]); err!=nil{
					if c<0 {return false}
					db.ShouldQuery("UPDATE levels SET password=? WHERE id=?",c,cl.Id)
					RegisterAction(ACTION_LEVEL_UPDATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Copy:Password"},db)
					return true
				}
				return false
			case "chown":
				if !isOwner && (role.RoleName=="" || role.Privs["cLvlAccess"]!=1) {return false}
				if len(command)<4 {return false}
				if c,err:=strconv.Atoi(command[2]); err!=nil{
					if c!=cl.Id {return false}
					xacc:=CAccount{DB: db}
					uid:=xacc.GetUIDByUname(command[3],false)
					if uid<1 {return false}
					db.ShouldQuery("UPDATE levels SET iud=? WHERE id=?",uid,cl.Id)
					RegisterAction(ACTION_LEVEL_UPDATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Chown:"+command[3]},db)
					return true
				}
				return false
			case "desc":
				if !isOwner && (role.RoleName=="" || role.Privs["cLvlAccess"]!=1) {return false}
				if len(command)<3 || len(strings.Replace(comment,"!lvl desc ","",1))>256{return false}
				data:=base64.StdEncoding.EncodeToString([]byte(strings.Replace(comment,"!lvl desc ","",1)))
				db.ShouldQuery("UPDATE levels SET description=? WHERE id=?",data,cl.Id)
				RegisterAction(ACTION_LEVEL_UPDATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"UpdDescription:"+m},db)
				return true
			case "list":
				if !isOwner && (role.RoleName=="" || role.Privs["cLvlAccess"]!=1) {return false}
				db.ShouldQuery("UPDATE levels SET isUnlisted=0 WHERE id=?",cl.Id)
				RegisterAction(ACTION_LEVEL_UPDATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"List:"+m},db)
				return true
			case "unlist":
				if !isOwner && (role.RoleName=="" || role.Privs["cLvlAccess"]!=1) {return false}
				db.ShouldQuery("UPDATE levels SET isUnlisted=1 WHERE id=?",cl.Id)
				RegisterAction(ACTION_LEVEL_UPDATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Unist:"+m},db)
				return true
			case "friendlist":
				if !isOwner && (role.RoleName=="" || role.Privs["cLvlAccess"]!=1) {return false}
				db.ShouldQuery("UPDATE levels SET isUnlisted=2 WHERE id=?",cl.Id)
				RegisterAction(ACTION_LEVEL_UPDATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Friendlist:"+m},db)
				return true
			case "ldm":
				if !isOwner && (role.RoleName=="" || role.Privs["cLvlAccess"]!=1) {return false}
				if len(command)<3 {return false}
				switch command[2] {
				case "on":
					db.ShouldQuery("UPDATE levels SET isLDM=1 WHERE id=?",cl.Id)
					return true
				case "off":
					db.ShouldQuery("UPDATE levels SET isLDM=0 WHERE id=?",cl.Id)
					return true
				default:
					return false
				}
			default:
				return false
			}
		case "!song":
			if !isOwner && (role.RoleName=="" || role.Privs["cLvlAccess"]!=1) {return false}
			if len(command)<2 {return false}
			if c,err:=strconv.Atoi(command[1]); err!=nil{
				if c<0 {return false}
				db.ShouldQuery("UPDATE levels SET song_id=?,track_id=0 WHERE id=?",c,cl.Id)
				RegisterAction(ACTION_LEVEL_UPDATE,acc.Uid,cl.Id,map[string]string{"uname":acc.Uname,"type":"Song:"+command[1]},db)
				return true
			}
			return false
		}
	case "!collab":
		if !isOwner {return false}
		if len(command)<3 {return false}
		var req string
		db.MustQueryRow("SELECT collab FROM levels WHERE id=?",cl.Id).Scan(&req)
		collabMembers:=strings.Split(req,",")
		switch command[1] {
		case "add":
			xacc:=CAccount{DB: db}
			uid:=xacc.GetUIDByUname(command[2],false)
			if uid<0 {return false}
			if !slices.Contains(collabMembers,strconv.Itoa(uid)) {collabMembers=append(collabMembers,strconv.Itoa(uid))}
			break
		case "del":
			xacc:=CAccount{DB: db}
			uid:=xacc.GetUIDByUname(command[2],false)
			if uid<0 {return false}
			if slices.Contains(collabMembers,strconv.Itoa(uid)) {
				i:=slices.Index(collabMembers,strconv.Itoa(uid))
				collabMembers=sliceRemove(collabMembers,i)
			}
			break
		default:
			return false
		}
		db.ShouldQuery("UPDATE levles SET collab=? WHERE id=?",strings.Join(collabMembers,","),cl.Id)
		return true
	}
	return false
}
