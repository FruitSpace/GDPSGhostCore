package core

import (
	"slices"
	"strconv"
	"strings"
	"time"
)

type CMessage struct {
	Id         int    `json:"id"`
	UidSrc     int    `json:"uid_src"`
	UidDest    int    `json:"uid_dest"`
	Subject    string `json:"subject"`
	Message    string `json:"message"`
	PostedTime string `json:"posted_time"`
	IsNew      bool   `json:"is_new"`

	DB *MySQLConn `json:"-"`
}

func (cm *CMessage) Exists(id int) bool {
	var cnt int
	cm.DB.MustQueryRow("SELECT count(*) as cnt FROM #DB#.messages WHERE id=?", id).Scan(&cnt)
	return cnt > 0
}

func (cm *CMessage) CountMessages(uid int, isNew bool) int {
	var cnt int
	var postfix string
	if isNew {
		postfix = " AND isNew=1"
	}
	cm.DB.MustQueryRow("SELECT count(*) as cnt FROM #DB#.messages WHERE uid_dest=?"+postfix, uid).Scan(&cnt)
	return cnt
}

func (cm *CMessage) LoadMessageById(id int) {
	if id > 0 {
		cm.Id = id
	}
	if cm.DB.ShouldQueryRow("SELECT uid_src,uid_dest,subject,body,postedTime,isNew FROM #DB#.messages WHERE id=?", cm.Id).Scan(
		&cm.UidSrc, &cm.UidDest, &cm.Subject, &cm.Message, &cm.PostedTime, &cm.IsNew) == nil {
		cm.DB.ShouldExec("UPDATE #DB#.messages SET isNew=0 WHERE id=?", cm.Id)
	}
}

func (cm *CMessage) DeleteMessage(uid int) {
	cm.DB.ShouldExec("DELETE FROM #DB#.messages WHERE id=? AND (uid_src=? OR uid_dest=?)", cm.Id, uid, uid)
}

func (cm *CMessage) SendMessageObj() bool {
	if len(cm.Subject) > 256 || len(cm.Message) > 1024 {
		return false
	}
	acc := CAccount{DB: cm.DB}
	acc.Uid = cm.UidDest
	acc.LoadSettings()
	if acc.MS == 2 {
		return false
	}
	acc.LoadSocial()
	blacklist := strings.Split(acc.Blacklist, ",")
	if slices.Contains(blacklist, strconv.Itoa(cm.UidSrc)) {
		return false
	}
	if acc.MS == 1 {
		cf := CFriendship{DB: cm.DB}
		if !cf.IsAlreadyFriend(cm.UidSrc, cm.UidDest) {
			return false
		}
	}
	cm.DB.ShouldExec("INSERT INTO #DB#.messages (uid_src,uid_dest,subject,body,postedTime) VALUES(?,?,?,?,?)",
		cm.UidSrc, cm.UidDest, cm.Subject, cm.Message, time.Now().Format("2006-01-02 15:04:05"))
	return true
}

func (cm *CMessage) GetMessageForUid(uid int, page int, sent bool) (int, []map[string]string) {
	page *= 10
	var cnt int
	pf := "uid_dest"
	if sent {
		pf = "uid_src"
	}
	cm.DB.MustQueryRow("SELECT count(*) as cnt FROM #DB#.messages WHERE "+pf+"=?", uid).Scan(&cnt)
	if cnt == 0 {
		return 0, []map[string]string{}
	}
	rows := cm.DB.ShouldQuery("SELECT id,uid_src,uid_dest,subject,body,postedTime,isNew FROM #DB#.messages WHERE "+pf+"=? ORDER BY id limit 10 OFFSET "+strconv.Itoa(page), uid)
	defer rows.Close()
	var out []map[string]string
	for rows.Next() {
		msg := CMessage{}
		rows.Scan(&msg.Id, &msg.UidSrc, &msg.UidDest, &msg.Subject, &msg.Message, &msg.PostedTime, &msg.IsNew)
		blk := map[string]string{
			"id":      strconv.Itoa(msg.Id),
			"subject": msg.Subject,
			"message": msg.Message,
			"isOld":   strconv.Itoa(ToInt(!msg.IsNew)),
			"date":    msg.PostedTime,
		}

		uid := msg.UidSrc
		if sent {
			uid = msg.UidDest
		}
		acc := CAccount{DB: cm.DB, Uid: uid}
		if acc.Exists(uid) {
			acc.LoadAuth(CAUTH_UID)
			blk["uname"] = acc.Uname
		} else {
			blk["uname"] = "[DELETED]"
		}
		blk["uid"] = strconv.Itoa(uid)
		if msg.IsNew {
			cm.DB.ShouldExec("UPDATE #DB#.messages SET isNew=0 WHERE id=?", msg.Id)
		}
		out = append(out, blk)
	}
	return cnt, out
}
