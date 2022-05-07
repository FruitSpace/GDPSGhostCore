package core

import (
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
	"time"
)

type CFriendship struct {
	DB MySQLConn
	Logger Logger
}

func (cf *CFriendship) IsAlreadyFriend(uid_dest int, uid int) bool {
	var cnt int
	cf.Logger.Must(cf.DB,cf.DB.DB.QueryRow("SELECT count(*) as cnt FROM friendships WHERE (uid1=? AND uid2=?) OR (uid2=? AND uid1=?)",
		uid,uid_dest,uid,uid_dest).Scan(&cnt))
	return cnt>0
}

func (cf *CFriendship) IsAlreadySentFriend(uid_dest int ,uid int) bool {
	var cnt int
	cf.Logger.Must(cf.DB,cf.DB.DB.QueryRow("SELECT count(*) as cnt FROM friendreqs WHERE uid_src=? AND uid_dest=?",uid,uid_dest).Scan(&cnt))
	return cnt>0
}

func (cf *CFriendship) CountFriendRequests(uid int, new bool) int {
	var cnt int
	q:="SELECT count(*) as cnt FROM friendreqs WHERE uid_dest=?"
	if new {q+=" AND isNew=1"}
	cf.Logger.Must(cf.DB,cf.DB.DB.QueryRow(q,uid).Scan(&cnt))
	return cnt
}

func (cf *CFriendship) GetFriendRequests(uid int, page int, sent bool) (int,[]map[string]string) {
	q:="SELECT id,uid_src,uid_dest,uploadDate,comment,isNew FROM friendreqs WHERE "
	if sent {q+="uid_src=?"}else{q+="uid_dest=?"}
	q+=" LIMIT 10 OFFSET ?"
	rows,_:=cf.DB.DB.Query(q,uid,page)
	var users []map[string]string
	var cnt int
	for rows.Next() {
		var (
			id int
			src int
			dest int
			date string
			comment string
			isNew int
		)
		rows.Scan(&id,&src,&dest,&date,&comment,&isNew)
		cnt++
		var user map[string]string
		user["id"]=strconv.Itoa(id)
		user["comment"]=comment
		acc:=CAccount{DB:cf.DB, Logger: cf.Logger}
		if sent {acc.Uid=dest}else{acc.Uid=src}
		user["uid"]=strconv.Itoa(acc.Uid)
		acc.LoadAuth(CAUTH_UID)
		acc.LoadStats()
		acc.LoadVessels()
		user["uname"]=acc.Uname
		user["isNew"]=strconv.Itoa(isNew)
		user["special"]=strconv.Itoa(acc.Special)
		user["iconType"]=strconv.Itoa(acc.IconType)
		user["clr_primary"]=strconv.Itoa(acc.ColorPrimary)
		user["clr_secondary"]=strconv.Itoa(acc.ColorSecondary)
		user["iconId"]=strconv.Itoa(acc.GetShownIcon())
		user["date"]=date
		users=append(users,user)
	}
	return cnt,users
}

func (cf *CFriendship) GetFriendRequestsCount(uid int, sent bool) int {
	var cnt int
	q:="SELECT count(*) as cnt FROM friendreqs WHERE "
	if sent {q+="uid_src=?"}else{q+="uid_dest=?"}
	cf.Logger.Must(cf.DB,cf.DB.DB.QueryRow(q,uid).Scan(&cnt))
	return cnt
}

func (cf *CFriendship) DeleteFriendship(uid int, uid_dest int) {
	id:=cf.GetFriendshipId(uid,uid_dest)
	if id==0 {return}
	cf.DB.DB.Query("DELETE FROM friendships WHERE (uid1=? AND uid2=?) OR (uid2=? AND uid1=?)",uid,uid_dest,uid,uid_dest)
	u1:=CAccount{DB: cf.DB, Logger: cf.Logger}
	u2:=CAccount{DB: cf.DB, Logger: cf.Logger}
	u1.Uid=uid
	u2.Uid=uid_dest
	u1.UpdateFriendships(CFRIENDSHIP_REMOVE,id)
	u2.UpdateFriendships(CFRIENDSHIP_REMOVE,id)
}

func (cf *CFriendship) GetFriendshipId(uid int, uid_dest int) int {
	var id int
	cf.DB.DB.QueryRow("SELECT id FROM friendships WHERE (uid1=? AND uid2=?) OR (uid2=? AND uid1=?)",uid,uid_dest,uid,uid_dest).Scan(&id)
	return id
}

func (cf *CFriendship) GetFriendByFID(id int) (int,int) {
	var (
		uid1 int
		uid2 int
	)
	cf.DB.DB.QueryRow("SELECT uid1,uid2 FROM friendships WHERE id=?",id).Scan(&uid1,&uid2)
	return uid1,uid2
}

func (cf *CFriendship) GetAccFriends(acc CAccount) []map[string]int {
	fr:=strings.Split(acc.FriendshipIds,",")
	var frlist []map[string]int
	for _,sfr:=range fr {
		id,err:=strconv.Atoi(sfr)
		if err!=nil {continue}
		uid1,uid2:=cf.GetFriendByFID(id)
		frlist=append(frlist,map[string]int{"uid1":uid1,"uid2":uid2})
	}
	return frlist
}

func (cf *CFriendship) ReadFriendRequest(id int) {
	cf.DB.DB.QueryRow("UPDATE friendreqs SET isNew=0 WHERE id=?",id)
}

func (cf *CFriendship) RequestFriend(uid int, uid_dest int, comment string) int {
	if uid==uid_dest || cf.IsAlreadyFriend(uid,uid_dest) || cf.IsAlreadySentFriend(uid,uid_dest) || len(comment)>512 {return -1}
	acc:=CAccount{DB: cf.DB, Logger: cf.Logger}
	acc.Uid=uid_dest
	acc.LoadSettings()
	if acc.FrS>0 {return -1}
	acc.LoadSocial()
	blacklist:=strings.Split(acc.Blacklist,",")
	if slices.Contains(blacklist,uid) {return -1}
	acc.DB.DB.Query("INSERT INTO friendreqs (uid_src, uid_dest, uploadDate, comment) VALUES (?,?,?,?)",
		uid,uid_dest,time.Now().Format("2006-01-02 15:04:05"),comment)
	return 1
}

func (cf *CFriendship) AcceptFriendRequest(id int, uid int) int {
	var (
		src int
		dest int
	)
	cf.DB.DB.QueryRow("SELECT uid_src, uid_dest FROM friendreqs WHERE id=",id).Scan(&src,&dest)
	if src==dest || uid!=dest {return -1}
	req,_:=cf.DB.DB.Prepare("INSERT INTO friendships (uid1, uid2) VALUES (?,?)")
	rq,_:=req.Exec(uid,dest)
	iid,_:=rq.LastInsertId()
	cf.DB.DB.Query("DELETE FROM friendreqs WHERE id=?",id)
	u1:=CAccount{DB: cf.DB, Logger: cf.Logger}
	u2:=CAccount{DB: cf.DB, Logger: cf.Logger}
	u1.Uid=uid
	u2.Uid=dest
	res:=u1.UpdateFriendships(CFRIENDSHIP_ADD,int(iid))
	res+=u2.UpdateFriendships(CFRIENDSHIP_ADD,int(iid))
	if res!=2 {return -1}
	return 1
}

func (cf *CFriendship) RejectFriendRequestById(id int, uid int) int {
	var (
		uid1 int
		uid2 int
	)
	cf.DB.DB.QueryRow("SELECT uid_src, uid_dest FROM friendreqs WHERE id=?",id).Scan(&uid1,&uid2)
	if uid1==uid2 || uid!=uid2 {return -1}
	cf.DB.DB.Query("DELETE FROM friendreqs WHERE id=?",id)
	return 1
}

func (cf *CFriendship) RejectFriendRequestByUid(uid int, uid_dest int, isSender bool){
	var (
		uid1 int
		uid2 int
	)
	if isSender {
		uid1=uid
		uid2=uid_dest
	}else{
		uid1=uid_dest
		uid2=uid
	}
	cf.DB.DB.Query("DELETE FROM friendreqs WHERE uid_src=? AND uid_dest=?",uid1,uid2)
}