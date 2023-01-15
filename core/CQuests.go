package core

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	QUEST_TYPE_DAILY     = 0
	QUEST_TYPE_WEEKLY    = 1
	QUEST_TYPE_EVENT     = -1
	QUEST_TYPE_CHALLENGE = 2
)

type CQuests struct {
	DB *MySQLConn
}

func (cq *CQuests) Exists(cType int) bool {
	switch cType {
	case -1:
		cType = 0
	case -2:
		cType = 1
	case -3:
		cType = -1
	}
	xType := "=" + strconv.Itoa(cType)
	if cType == 2 {
		xType = ">1"
	}
	var cnt int
	cq.DB.MustQueryRow("SELECT count(*) as cnt FROM #DB#.quests WHERE type" + xType).Scan(&cnt)
	return cnt > 0
}

func (cq *CQuests) GetDaily() (int, int) {
	var id, lvlId int
	cq.DB.ShouldQueryRow("SELECT id, lvl_id FROM #DB#.quests WHERE type=0 AND timeExpire<now() ORDER BY timeExpire DESC LIMIT 1").Scan(&id, &lvlId)
	return id, lvlId
}

func (cq *CQuests) GetWeekly() (int, int) {
	var id, lvlId int
	cq.DB.ShouldQueryRow("SELECT id, lvl_id FROM #DB#.quests WHERE type=1 AND timeExpire<now() ORDER BY timeExpire DESC LIMIT 1").Scan(&id, &lvlId)
	return id, lvlId + 100001
}

func (cq *CQuests) GetEvent() (int, int) {
	var id, lvlId int
	cq.DB.ShouldQueryRow("SELECT id, lvl_id FROM #DB#.quests WHERE type=-1 AND timeExpire<now() ORDER BY timeExpire DESC LIMIT 1").Scan(&id, &lvlId)
	return id, lvlId
}

func (cq *CQuests) PushLevel(lvlId, cType int) int {
	res := cq.DB.ShouldPrepareExec("INSERT INTO #DB#.quests (type,lvl_id) VALUES (?,?)", cType, lvlId)
	id, _ := res.LastInsertId()
	return int(id)
}

func (cq *CQuests) GetQuests(uid int) string {
	rand.Seed(int64(time.Now().YearDay() * uid))
	req := cq.DB.ShouldQuery("SELECT r1.id,type,needed,reward,name,timeExpire FROM #DB#.quests AS r1 " +
		"JOIN (SELECT CEIL(" + fmt.Sprintf("%f", rand.Float64()) + " * (SELECT MAX(id) FROM #DB#.quests WHERE type>1)) AS id) AS r2 " +
		"WHERE r1.id >= r2.id AND r1.timeExpire<now() AND r1.type>1 ORDER BY r1.id ASC LIMIT 3")
	defer req.Close()
	out := ""
	var cnt int
	for req.Next() {
		var id, xType, needed, reward int
		var name, timeExpire string
		req.Scan(&id, &xType, &needed, &reward, &name, &timeExpire)
		out += strconv.Itoa(id) + "," + strconv.Itoa(xType-1) + "," + strconv.Itoa(needed) + "," + strconv.Itoa(reward) + "," + name + ":"
		cnt++
	}
	for cnt < 3 {
		out += "1337,2,5,50,Random issues:"
		cnt++
	}
	return out[:len(out)-1]
}

func (cq *CQuests) GetSpecialLevel(xType int) string {
	timeLeft := 0
	var lvlId, xLvlid int
	tme, _ := time.ParseInLocation("2006-01-02 15:04:05", strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")[0]+" 00:00:00", loc)
	switch xType {
	case -1:
		fallthrough
	case 0:
		//!Additional 10800 Review is needed
		timeLeft = int(tme.AddDate(0, 0, 1).Unix() - (time.Now().Unix()))
	case 1:
		timeLeft = int(tme.AddDate(0, 0, 7).Unix() - (time.Now().Unix()))
		lvlId = 100001
	}
	cq.DB.ShouldQueryRow("SELECT lvl_id FROM #DB#.quests WHERE type=" + strconv.Itoa(xType) + " AND timeExpire<now() ORDER BY timeExpire DESC LIMIT 1").Scan(&xLvlid)
	return strconv.Itoa(xLvlid+lvlId) + "|" + strconv.Itoa(timeLeft)
}
