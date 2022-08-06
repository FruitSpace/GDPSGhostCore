package core

import (
	"fmt"
	"strings"
	"time"
)

const (
	CSCORE_TYPE_FRIENDS int = 400
	CSCORE_TYPE_TOP int = 401
	CSCORE_TYPE_WEEK int = 402
)

type CScores struct {
	Id int
	Uid int
	LvlId int
	PostedTime string
	Percent int
	Ranking int
	Attempts int
	Coins int

	DB MySQLConn
}

func (cs *CScores) ScoreExistsByUid(uid int, lvlId int) bool {
	var cnt int
	cs.DB.MustQueryRow("SELECT count(*) as cnt FROM scores WHERE uid=? AND lvl_id=?",uid,lvlId).Scan(&cnt)
	return cnt>0
}

func (cs *CScores) LoadScoreById() {
	cs.DB.ShouldQueryRow("SELECT uid,lvl_id,postedTime,percent,attempts,coins FROM scores WHERE id=?",cs.Id).Scan(
		&cs.Uid,&cs.LvlId,&cs.PostedTime,&cs.Percent,&cs.Attempts,&cs.Coins)
}

func (cs *CScores) GetScoresForLevelId(lvlId int, types int, acc CAccount) []CScores {
	var suffix string
	switch types {
	case CSCORE_TYPE_WEEK:
		date:=strings.Split(time.Now().AddDate(0,0,8-int(time.Now().Weekday())).Format("2006-01-02 15:04:05")," ")[0]+" 00:00:00"
		suffix="AND postedTime>='"+date+"'"
		break
	case CSCORE_TYPE_FRIENDS:
		acc.LoadSocial()
		cf:=CFriendship{DB: cs.DB}
		frs:=cf.GetAccFriends(acc)
		frs=append(frs,acc.Uid)
		xfrs:=strings.Trim(strings.Join(strings.Fields(fmt.Sprint(frs)), ","), "[]")
		suffix="AND uid IN("+strings.ReplaceAll(xfrs,",,",",")+")"
		break
	}
	req:=cs.DB.ShouldQuery("SELECT uid,lvl_id,postedTime,percent,attempts,coins FROM scores WHERE lvl_id=? "+suffix+" ORDER BY percent DESC")
	var scores []CScores
	for req.Next() {
		xcs:=CScores{}
		req.Scan(&xcs.Uid,&xcs.LvlId,&xcs.PostedTime,&xcs.Percent,&xcs.Attempts,&xcs.Coins)
		if xcs.Percent==100 { xcs.Ranking=1 }else if xcs.Percent>=75 { xcs.Ranking=2 }else{ xcs.Ranking=3 }
		scores=append(scores,xcs)
	}
	return scores
}

func (cs *CScores) UpdateLevelScore() {
	cs.DB.ShouldQuery("UPDATE scores SET postedTime=?,percent=?,attempts=?,coins=? WHERE lvl-id=? AND uid=?",
		time.Now().Format("2006-01-02 15:04:05"),cs.Percent,cs.Attempts,cs.Coins)
}

func (cs *CScores) UploadLevelScore() {
	cs.DB.ShouldQuery("INSERT INTO scores (uid,lvl_id,postedTime,percent,attempts,coins) VALUES(?,?,?,?,?,?)",
		cs.Uid,cs.LvlId,time.Now().Format("2006-01-02 15:04:05"),cs.Percent,cs.Attempts,cs.Coins)
}