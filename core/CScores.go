package core

import (
	"fmt"
	"strings"
	"time"
)

const (
	CSCORE_TYPE_FRIENDS int = 400
	CSCORE_TYPE_TOP     int = 401
	CSCORE_TYPE_WEEK    int = 402
	CSCORE_PLAT_FRIENDS     = 500
	CSCORE_PLAT_TOP         = 501
	CSCORE_PLAT_WEEK        = 502
)

type CScores struct {
	Id         int
	Uid        int
	LvlId      int
	PostedTime string
	Percent    int
	Ranking    int
	Attempts   int
	Coins      int

	DB *MySQLConn
}

func (cs *CScores) ScoreExistsByUid(uid int, lvlId int) bool {
	var cnt int
	cs.DB.MustQueryRow("SELECT count(*) as cnt FROM #DB#.scores WHERE uid=? AND lvl_id=?", uid, lvlId).Scan(&cnt)
	return cnt > 0
}

func (cs *CScores) LoadScoreById() {
	cs.DB.ShouldQueryRow("SELECT uid,lvl_id,postedTime,percent,attempts,coins FROM #DB#.scores WHERE id=?", cs.Id).Scan(
		&cs.Uid, &cs.LvlId, &cs.PostedTime, &cs.Percent, &cs.Attempts, &cs.Coins)
}

func (cs *CScores) GetScoresForLevelId(lvlId int, types int, acc CAccount) []CScores {
	var suffix string
	switch types {
	case CSCORE_TYPE_WEEK:
		date := strings.Split(time.Now().AddDate(0, 0, 8-int(time.Now().Weekday())).Format("2006-01-02 15:04:05"), " ")[0] + " 00:00:00"
		suffix = "AND postedTime>='" + date + "'"
	case CSCORE_TYPE_FRIENDS:
		acc.LoadSocial()
		cf := CFriendship{DB: cs.DB}
		frs := cf.GetAccFriends(acc)
		frs = append(frs, acc.Uid)
		xfrs := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(frs)), ","), "[]")
		suffix = "AND uid IN(" + strings.ReplaceAll(xfrs, ",,", ",") + ")"
	}
	req := cs.DB.ShouldQuery("SELECT uid,lvl_id,postedTime,percent,attempts,coins FROM #DB#.scores WHERE lvl_id=? "+suffix+" ORDER BY percent DESC", lvlId)
	var scores []CScores
	defer req.Close()
	for req.Next() {
		xcs := CScores{DB: cs.DB}
		req.Scan(&xcs.Uid, &xcs.LvlId, &xcs.PostedTime, &xcs.Percent, &xcs.Attempts, &xcs.Coins)
		if xcs.Percent == 100 {
			xcs.Ranking = 1
		} else if xcs.Percent >= 75 {
			xcs.Ranking = 2
		} else {
			xcs.Ranking = 3
		}
		scores = append(scores, xcs)
	}
	return scores
}

func (cs *CScores) GetScoresForPlatformerLevelId(lvlId int, types int, modeCoins bool, acc CAccount) []CScores {
	var suffix string
	switch types {
	case CSCORE_PLAT_WEEK:
		date := strings.Split(time.Now().AddDate(0, 0, 8-int(time.Now().Weekday())).Format("2006-01-02 15:04:05"), " ")[0] + " 00:00:00"
		suffix = "AND postedTime>='" + date + "'"
	case CSCORE_PLAT_FRIENDS:
		acc.LoadSocial()
		cf := CFriendship{DB: cs.DB}
		frs := cf.GetAccFriends(acc)
		frs = append(frs, acc.Uid)
		xfrs := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(frs)), ","), "[]")
		suffix = "AND uid IN(" + strings.ReplaceAll(xfrs, ",,", ",") + ")"
	}
	req := cs.DB.ShouldQuery("SELECT uid,lvl_id,postedTime,percent,attempts,coins FROM #DB#.scores WHERE lvl_id=? "+suffix+" ORDER BY percent DESC", lvlId)
	var scores []CScores
	defer req.Close()
	rankx := 1
	for req.Next() {
		xcs := CScores{DB: cs.DB}
		req.Scan(&xcs.Uid, &xcs.LvlId, &xcs.PostedTime, &xcs.Percent, &xcs.Attempts, &xcs.Coins)
		xcs.Ranking = rankx
		rankx++
		scores = append(scores, xcs)
	}
	return scores
}

func (cs *CScores) UpdateLevelScore() {
	cs.DB.ShouldExec("UPDATE #DB#.scores SET postedTime=?,percent=?,attempts=?,coins=? WHERE lvl_id=? AND uid=?",
		time.Now().Format("2006-01-02 15:04:05"), cs.Percent, cs.Attempts, cs.Coins, cs.LvlId, cs.Uid)
}

func (cs *CScores) UploadLevelScore() {
	cs.DB.ShouldExec("INSERT INTO #DB#.scores (uid,lvl_id,postedTime,percent,attempts,coins) VALUES(?,?,?,?,?,?)",
		cs.Uid, cs.LvlId, time.Now().Format("2006-01-02 15:04:05"), cs.Percent, cs.Attempts, cs.Coins)
}
