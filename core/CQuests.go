package core

import "strconv"

const (
	QUEST_TYPE_DAILY = 0
	QUEST_TYPE_WEEKLY = 1
	QUEST_TYPE_EVENT = -1
	QUEST_TYPE_CHALLENGE = 2
)

type CQuests struct {
	DB MySQLConn
}

func (cq *CQuests) Exists(cType int) bool {
	xType:="="+strconv.Itoa(cType)
	if cType==2 {xType=">1"}
	var cnt int
	cq.DB.MustQueryRow("SELECT count(*) as cnt FROM quests WHERE type"+xType).Scan(&cnt)
	return cnt>0
}
