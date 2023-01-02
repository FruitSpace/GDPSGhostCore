package core

import (
	"encoding/json"
	"time"
)

const (
	ACTION_USER_REGISTER         int = 110
	ACTION_USER_LOGIN            int = 111
	ACTION_USER_DELETE           int = 112
	ACTION_BAN_BAN               int = 113
	ACTION_BAN_UNBAN             int = 114
	ACTION_LEVEL_UPLOAD          int = 115
	ACTION_LEVEL_DELETE          int = 116
	ACTION_LEVEL_UPDATE          int = 117
	ACTION_LEVEL_RATE            int = 109
	ACTION_PANEL_GAUNTLET_ADD    int = 118
	ACTION_PANEL_GAUNTLET_DELETE int = 119
	ACTION_PANEL_GAUNTLET_EDIT   int = 120
	ACTION_PANEL_MAPPACK_ADD     int = 121
	ACTION_PANEL_MAPPACK_DELETE  int = 122
	ACTION_PANEL_MAPPACK_EDIT    int = 123
	ACTION_PANEL_QUEST_ADD       int = 124
	ACTION_PANEL_QUEST_DELETE    int = 125
	ACTION_PANEL_QUEST_EDIT      int = 126
	ACTION_LEVEL_LIKE            int = 127
	ACTION_ACCCOMMENT_LIKE       int = 128
	ACTION_COMMENT_LIKE          int = 129

	ITEMTYPE_LEVEL      int = 130
	ITEMTYPE_ACCCOMMENT int = 131
	ITEMTYPE_COMMENT    int = 132
)

func RegisterAction(action int, uid int, target_id int, data map[string]string, db *MySQLConn) {
	var types int
	switch action {
	case ACTION_USER_REGISTER:
		types = 0
		data["action"] = "Register"
	case ACTION_USER_LOGIN:
		types = 1
		data["action"] = "Login"
	case ACTION_USER_DELETE:
		types = 2
		data["action"] = "Delete"
	case ACTION_BAN_BAN:
		types = 3
		data["action"] = "Ban"
	case ACTION_BAN_UNBAN:
		types = 3
		data["action"] = "Unban"
	case ACTION_LEVEL_UPLOAD:
		types = 4
		data["action"] = "Upload"
	case ACTION_LEVEL_DELETE:
		types = 4
		data["action"] = "Delete"
	case ACTION_LEVEL_UPDATE:
		types = 4
		data["action"] = "Update"
	case ACTION_LEVEL_RATE:
		types = 4
		data["action"] = "Rate"
	case ACTION_PANEL_GAUNTLET_ADD:
		types = 5
		data["action"] = "GauntletAdd"
	case ACTION_PANEL_GAUNTLET_DELETE:
		types = 5
		data["action"] = "GauntletDelete"
	case ACTION_PANEL_GAUNTLET_EDIT:
		types = 5
		data["action"] = "GauntletEdit"
	case ACTION_PANEL_MAPPACK_ADD:
		types = 5
		data["action"] = "MapPackAdd"
	case ACTION_PANEL_MAPPACK_DELETE:
		types = 5
		data["action"] = "MapPackDelete"
	case ACTION_PANEL_MAPPACK_EDIT:
		types = 5
		data["action"] = "MapPackEdit"
	case ACTION_PANEL_QUEST_ADD:
		types = 5
		data["action"] = "QuestAdd"
	case ACTION_PANEL_QUEST_DELETE:
		types = 5
		data["action"] = "QuestDelete"
	case ACTION_PANEL_QUEST_EDIT:
		types = 5
		data["action"] = "QuestEdit"
	case ACTION_LEVEL_LIKE:
		types = 6
		data["action"] = "LikeLevel"
	case ACTION_ACCCOMMENT_LIKE:
		types = 7
		data["action"] = "LikeAcccomment"
	case ACTION_COMMENT_LIKE:
		types = 8
		data["action"] = "LikeComment"
	default:
		return
	}
	isMod := 0
	if uid > 0 {
		ret := 0
		db.ShouldQueryRow("SELECT role_id FROM #DB#.users WHERE uid=?", uid).Scan(&ret)
		if ret > 0 {
			isMod = 1
		}
	}
	datac, _ := json.Marshal(data)
	date := time.Now().Format("2006-01-02 15:04:05")
	db.ShouldExec("INSERT INTO #DB#.actions (date, uid, type, target_id, isMod, data) VALUES (?,?,?,?,?,?)",
		date,
		uid,
		types,
		target_id,
		isMod,
		string(datac),
	)
}

func IsLiked(itemType int, uid int, dest_id int, db *MySQLConn) bool {
	event_id := 0
	switch itemType {
	case ITEMTYPE_LEVEL:
		event_id = 6
	case ITEMTYPE_ACCCOMMENT:
		event_id = 7
	case ITEMTYPE_COMMENT:
		event_id = 8
	default:
		return true
	}
	var q int
	db.ShouldQueryRow("SELECT count(*) as cnt FROM #DB#.actions WHERE type=? AND uid=? AND target_id=?", event_id, uid, dest_id).Scan(&q)
	return q > 0
}
