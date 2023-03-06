package core

import (
	"database/sql"
	"encoding/base64"
	"strconv"
	"strings"
	"time"
)

func privErr(s string) string {
	return "You need '" + s + "' privilege to execute this command"
}

func InvokeCommands(db *MySQLConn, cl CLevel, acc CAccount, comment string, isOwner bool, role Role) string {
	command := strings.Split(comment, " ")
	switch command[0] {
	case "!feature":
		if role.RoleName == "" || role.Privs["cFeature"] != 1 {
			return privErr("cFeature")
		}
		cl.FeatureLevel(1)
		RegisterAction(ACTION_LEVEL_RATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Feature"}, db)
		return "ok"
	case "!legendary":
		if role.RoleName == "" || role.Privs["cEpic"] != 1 {
			return privErr("cEpic")
		}
		cl.FeatureLevel(3)
		cl.LegendaryLevel(true)
		RegisterAction(ACTION_LEVEL_RATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Legendary"}, db)
		return "ok"
	case "!godlike":
		if role.RoleName == "" || role.Privs["cEpic"] != 1 {
			return privErr("cEpic")
		}
		cl.FeatureLevel(4)
		RegisterAction(ACTION_LEVEL_RATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Feature"}, db)
		return "ok"
	case "!unfeature":
		if role.RoleName == "" || role.Privs["cFeature"] != 1 {
			return privErr("cFeature")
		}
		if !cl.FeatureLevel(0) {
			return "You need to unepic level first"
		}
		RegisterAction(ACTION_LEVEL_RATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Uneature"}, db)
		return "ok"
	case "!epic":
		if role.RoleName == "" || role.Privs["cEpic"] != 1 {
			return privErr("cEpic")
		}
		cl.FeatureLevel(2)
		cl.EpicLevel(true)
		RegisterAction(ACTION_LEVEL_RATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Epic"}, db)
		return "ok"
	case "!unepic":
		if role.RoleName == "" || role.Privs["cEpic"] != 1 {
			return privErr("cEpic")
		}
		cl.FeatureLevel(0)
		cl.EpicLevel(false)
		RegisterAction(ACTION_LEVEL_RATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Unepic"}, db)
		return "ok"
	case "!coins":
		if role.RoleName == "" || role.Privs["cVerCoins"] != 1 {
			return privErr("cVerCoins")
		}
		if len(command) < 2 {
			return "Specify 'verify' or 'reset' argument"
		}
		if command[1] == "verify" {
			cl.VerifyCoins(true)
			RegisterAction(ACTION_LEVEL_RATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Coins:Verify"}, db)
		} else if command[1] == "reset" {
			cl.VerifyCoins(false)
			RegisterAction(ACTION_LEVEL_RATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Coins:Reset"}, db)
		} else {
			return "Invalid argument. Specify 'verify' or 'reset' argument"
		}
		return "ok"
	case "!daily":
		if role.RoleName == "" || role.Privs["cDaily"] != 1 {
			return privErr("cDaily")
		}
		if len(command) > 1 && command[1] == "reset" {
			db.ShouldExec("DELETE FROM #DB#.quests WHERE lvl_id=? and type=0", cl.Id)
			RegisterAction(ACTION_LEVEL_RATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Daily:Reset"}, db)
		} else {
			var date string
			err := db.ShouldQueryRow("SELECT timeExpire FROM #DB#.quests WHERE type=0 ORDER BY timeExpire DESC LIMIT 1").Scan(&date)
			if err == sql.ErrNoRows {
				date = strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")[0] + " 00:00:00"
			} else {
				tme, _ := time.ParseInLocation("2006-01-02 15:04:05", strings.Split(date, " ")[0]+" 00:00:00", loc)
				tme.AddDate(0, 0, 1)
				date = tme.Format("2006-01-02 15:04:05")
			}
			db.ShouldExec("INSERT INTO #DB#.quests (type,lvl_id,timeExpire) VALUES (0,?,?)", cl.Id, date)
			RegisterAction(ACTION_LEVEL_RATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Daily:Publish"}, db)
		}
		return "ok"
	case "!weekly":
		if role.RoleName == "" || role.Privs["cWeekly"] != 1 {
			return privErr("cWeekly")
		}
		if len(command) > 1 && command[1] == "reset" {
			db.ShouldExec("DELETE FROM #DB#.quests WHERE lvl_id=? and type=1", cl.Id)
			RegisterAction(ACTION_LEVEL_RATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Weekly:Reset"}, db)
		} else {
			var date string
			err := db.ShouldQueryRow("SELECT timeExpire FROM #DB#.quests WHERE type=0 ORDER BY timeExpire DESC LIMIT 1").Scan(&date)
			if err == sql.ErrNoRows {
				date = strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")[0] + " 00:00:00"
			} else {
				tme, _ := time.ParseInLocation("2006-01-02 15:04:05", strings.Split(date, " ")[0]+" 00:00:00", loc)
				tme.AddDate(0, 0, 7)
				date = tme.Format("2006-01-02 15:04:05")
			}
			db.ShouldExec("INSERT INTO #DB#.quests (type,lvl_id,timeExpire) VALUES (1,?,?)", cl.Id, date)
			RegisterAction(ACTION_LEVEL_RATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Weekly:Publish"}, db)
		}
		return "ok"
	case "!rate":
		if role.RoleName == "" || role.Privs["cRate"] != 1 {
			return privErr("cRate")
		}
		if len(command) < 2 {
			return "Specify difficulty argument (easy, normal, hard, etc.) or 'reset'"
		}
		diff := "0"
		switch strings.ToLower(command[1]) {
		case "auto":
			diff = "-1"
		case "easy":
			diff = "10"
		case "normal":
			diff = "20"
		case "hard":
			diff = "30"
		case "harder":
			diff = "40"
		case "insane":
			diff = "50"
		case "reset":
			diff = "0,starsGot=0"
		default:
			return "Invalid difficulty argument"
		}
		db.ShouldExec("UPDATE #DB#.levels SET difficulty="+diff+" WHERE id=?", cl.Id)
		RegisterAction(ACTION_LEVEL_RATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Rate:" + strings.Title(strings.ToLower(command[1]))}, db)
		return "ok"
	case "!lvl":
		if len(command) < 2 {
			return "Specify subcommand (refer to docs)"
		}
		m := "Mod"
		if isOwner {
			m = "Owner"
		}
		switch command[1] {
		case "delete":
			if role.RoleName == "" || role.Privs["cDelete"] != 1 {
				return privErr("cDelete")
			}
			if len(command) < 3 || command[2] != strconv.Itoa(cl.Id) {
				return "Usage: !lvl delete <Level ID>"
			}
			cl.DeleteLevel()
			RegisterAction(ACTION_LEVEL_DELETE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Delete:" + m}, db)
			return "ok"
		case "rename":
			if !isOwner && (role.RoleName == "" || role.Privs["cLvlAccess"] != 1) {
				return privErr("cLvlAccess (or owner)")
			}
			if len(command) < 3 {
				return "Usage: !lvl rename <New name>"
			}
			text := strings.Replace(comment, "!lvl rename ", "", 1)
			db.ShouldExec("UPDATE #DB#.levels SET name=? WHERE id=?", text, cl.Id)
			RegisterAction(ACTION_LEVEL_UPDATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Rename:" + m}, db)
			return "ok"
		case "copy":
			if !isOwner {
				return "You are not level owner"
			}
			if len(command) < 3 {
				return "Usage: !lvl copy on/off/pass [password]"
			}
			switch command[2] {
			case "on":
				db.ShouldExec("UPDATE #DB#.levels SET password=1 WHERE id=?", cl.Id)
				RegisterAction(ACTION_LEVEL_UPDATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Copy:Enable"}, db)
				return "ok"
			case "off":
				db.ShouldExec("UPDATE #DB#.levels SET password=0 WHERE id=?", cl.Id)
				RegisterAction(ACTION_LEVEL_UPDATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Copy:Disable"}, db)
				return "ok"
			case "pass":
				if len(command) < 4 || len(command[3]) != 6 {
					return "Please specify valid password"
				}
				if c, err := strconv.Atoi(command[3]); err == nil {
					if c < 0 {
						return "Password should be positive number"
					}
					db.ShouldExec("UPDATE #DB#.levels SET password=? WHERE id=?", c, cl.Id)
					RegisterAction(ACTION_LEVEL_UPDATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Copy:Password"}, db)
					return "ok"
				}
				return "Password should be numeric"
			default:
				return "Usage: !lvl copy on/off/pass [password]"
			}
		case "chown":
			if !isOwner && (role.RoleName == "" || role.Privs["cLvlAccess"] != 1) {
				return privErr("cLvlAccess (or owner)")
			}
			if len(command) < 4 {
				return "Usage: !lvl chown <Level ID> <NewOwner username>"
			}
			if c, err := strconv.Atoi(command[2]); err != nil {
				if c != cl.Id {
					return "Level ID doesn't match"
				}
				xacc := CAccount{DB: db}
				uid := xacc.GetUIDByUname(command[3], false)
				if uid < 1 {
					return "New owner username not found"
				}
				db.ShouldExec("UPDATE #DB#.levels SET iud=? WHERE id=?", uid, cl.Id)
				RegisterAction(ACTION_LEVEL_UPDATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Chown:" + command[3]}, db)
				return "ok"
			}
			return "Level ID should be numeric"
		case "desc":
			if !isOwner && (role.RoleName == "" || role.Privs["cLvlAccess"] != 1) {
				return privErr("cLvlAccess (or owner)")
			}
			if len(command) < 3 || len(strings.Replace(comment, "!lvl desc ", "", 1)) > 256 {
				return "New description is too long (>256 symbols)"
			}
			data := base64.StdEncoding.EncodeToString([]byte(strings.Replace(comment, "!lvl desc ", "", 1)))
			db.ShouldExec("UPDATE #DB#.levels SET description=? WHERE id=?", data, cl.Id)
			RegisterAction(ACTION_LEVEL_UPDATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "UpdDescription:" + m}, db)
			return "ok"
		case "list":
			if !isOwner && (role.RoleName == "" || role.Privs["cLvlAccess"] != 1) {
				return privErr("cLvlAccess (or owner)")
			}
			db.ShouldExec("UPDATE #DB#.levels SET isUnlisted=0 WHERE id=?", cl.Id)
			RegisterAction(ACTION_LEVEL_UPDATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "List:" + m}, db)
			return "ok"
		case "unlist":
			if !isOwner && (role.RoleName == "" || role.Privs["cLvlAccess"] != 1) {
				return privErr("cLvlAccess (or owner)")
			}
			db.ShouldExec("UPDATE #DB#.levels SET isUnlisted=1 WHERE id=?", cl.Id)
			RegisterAction(ACTION_LEVEL_UPDATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Unist:" + m}, db)
			return "ok"
		case "friendlist":
			if !isOwner && (role.RoleName == "" || role.Privs["cLvlAccess"] != 1) {
				return privErr("cLvlAccess (or owner)")
			}
			db.ShouldExec("UPDATE #DB#.levels SET isUnlisted=2 WHERE id=?", cl.Id)
			RegisterAction(ACTION_LEVEL_UPDATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Friendlist:" + m}, db)
			return "ok"
		case "ldm":
			if !isOwner && (role.RoleName == "" || role.Privs["cLvlAccess"] != 1) {
				return privErr("cLvlAccess (or owner)")
			}
			if len(command) < 3 {
				return "Usage: !lvl ldm on/off"
			}
			switch command[2] {
			case "on":
				db.ShouldExec("UPDATE #DB#.levels SET isLDM=1 WHERE id=?", cl.Id)
				return "ok"
			case "off":
				db.ShouldExec("UPDATE #DB#.levels SET isLDM=0 WHERE id=?", cl.Id)
				return "ok"
			default:
				return "Usage: !lvl ldm on/off"
			}
		default:
			return "Invalid subcommand (refer to docs)"
		}
	case "!song":
		if !isOwner && (role.RoleName == "" || role.Privs["cLvlAccess"] != 1) {
			return privErr("cLvlAccess (or owner)")
		}
		if len(command) < 2 {
			return "Usage: !song <id>"
		}
		if c, err := strconv.Atoi(command[1]); err == nil {
			if c < 0 {
				return "Song ID should be positive number"
			}
			db.ShouldExec("UPDATE #DB#.levels SET song_id=?,track_id=0 WHERE id=?", c, cl.Id)
			RegisterAction(ACTION_LEVEL_UPDATE, acc.Uid, cl.Id, map[string]string{"uname": acc.Uname, "type": "Song:" + command[1]}, db)
			return "ok"
		}
		return "Song ID should be numeric"

		//case "!collab":
		//	if !isOwner {
		//		return false
		//	}
		//	if len(command) < 3 {
		//		return false
		//	}
		//	var req string
		//	db.MustQueryRow("SELECT collab FROM #DB#.levels WHERE id=?", cl.Id).Scan(&req)
		//	collabMembers := strings.Split(req, ",")
		//	switch command[1] {
		//	case "add":
		//		xacc := CAccount{DB: db}
		//		uid := xacc.GetUIDByUname(command[2], false)
		//		if uid < 0 {
		//			return false
		//		}
		//		if !slices.Contains(collabMembers, strconv.Itoa(uid)) {
		//			collabMembers = append(collabMembers, strconv.Itoa(uid))
		//		}
		//		break
		//	case "del":
		//		xacc := CAccount{DB: db}
		//		uid := xacc.GetUIDByUname(command[2], false)
		//		if uid < 0 {
		//			return false
		//		}
		//		if slices.Contains(collabMembers, strconv.Itoa(uid)) {
		//			i := slices.Index(collabMembers, strconv.Itoa(uid))
		//			collabMembers = sliceRemove(collabMembers, i)
		//		}
		//		break
		//	default:
		//		return false
		//	}
		//	db.ShouldExec("UPDATE #DB#.levles SET collab=? WHERE id=?", strings.Join(collabMembers, ","), cl.Id)
		//	return true
	}
	return "Invalid command (refer to docs)"
}
