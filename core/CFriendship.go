package core

import (
	"slices"
	"strconv"
	"strings"
	"time"
)

type CFriendship struct {
	DB *MySQLConn
}

func (cf *CFriendship) IsAlreadyFriend(uid_dest int, uid int) bool {
	var cnt int
	cf.DB.MustQueryRow("SELECT count(*) as cnt FROM #DB#.friendships WHERE (uid1=? AND uid2=?) OR (uid2=? AND uid1=?)",
		uid, uid_dest, uid, uid_dest).Scan(&cnt)
	return cnt > 0
}

func (cf *CFriendship) IsAlreadySentFriend(uid_dest int, uid int) bool {
	var cnt int
	cf.DB.MustQueryRow("SELECT count(*) as cnt FROM #DB#.friendreqs WHERE uid_src=? AND uid_dest=?", uid, uid_dest).Scan(&cnt)
	return cnt > 0
}

func (cf *CFriendship) CountFriendRequests(uid int, new bool) int {
	var cnt int
	q := "SELECT count(*) as cnt FROM #DB#.friendreqs WHERE uid_dest=?"
	if new {
		q += " AND isNew=1"
	}
	cf.DB.MustQueryRow(q, uid).Scan(&cnt)
	return cnt
}

func (cf *CFriendship) GetFriendRequests(uid int, page int, sent bool) (int, []map[string]string) {
	q := "SELECT id,uid_src,uid_dest,uploadDate,comment,isNew FROM #DB#.friendreqs WHERE "
	if sent {
		q += "uid_src=?"
	} else {
		q += "uid_dest=?"
	}
	q += " LIMIT 10 OFFSET " + strconv.Itoa(page*10)
	rows := cf.DB.MustQuery(q, uid)
	defer rows.Close()
	var users []map[string]string
	var cnt int
	for rows.Next() {
		var (
			id      int
			src     int
			dest    int
			date    string
			comment string
			isNew   int
		)
		rows.Scan(&id, &src, &dest, &date, &comment, &isNew)
		cnt++
		user := make(map[string]string)
		user["id"] = strconv.Itoa(id)
		user["comment"] = comment
		acc := CAccount{DB: cf.DB}
		if sent {
			acc.Uid = dest
		} else {
			acc.Uid = src
		}
		user["uid"] = strconv.Itoa(acc.Uid)
		acc.LoadAuth(CAUTH_UID)
		acc.LoadStats()
		acc.LoadVessels()
		user["uname"] = acc.Uname
		user["isNew"] = strconv.Itoa(isNew)
		user["special"] = strconv.Itoa(acc.Special)
		user["iconType"] = strconv.Itoa(acc.IconType)
		user["clr_primary"] = strconv.Itoa(acc.ColorPrimary)
		user["clr_secondary"] = strconv.Itoa(acc.ColorSecondary)
		user["iconId"] = strconv.Itoa(acc.GetShownIcon())
		user["date"] = date
		users = append(users, user)
	}
	return cnt, users
}

func (cf *CFriendship) GetFriendRequestsCount(uid int, sent bool) int {
	var cnt int
	q := "SELECT count(*) as cnt FROM #DB#.friendreqs WHERE "
	if sent {
		q += "uid_src=?"
	} else {
		q += "uid_dest=?"
	}
	cf.DB.MustQueryRow(q, uid).Scan(&cnt)
	return cnt
}

func (cf *CFriendship) DeleteFriendship(uid int, uid_dest int) {
	id := cf.GetFriendshipId(uid, uid_dest)
	if id == 0 {
		return
	}
	cf.DB.ShouldExec("DELETE FROM #DB#.friendships WHERE (uid1=? AND uid2=?) OR (uid2=? AND uid1=?)", uid, uid_dest, uid, uid_dest)
	u1 := CAccount{DB: cf.DB}
	u2 := CAccount{DB: cf.DB}
	u1.Uid = uid
	u2.Uid = uid_dest
	u1.UpdateFriendships(CFRIENDSHIP_REMOVE, id)
	u2.UpdateFriendships(CFRIENDSHIP_REMOVE, id)
}

func (cf *CFriendship) GetFriendshipId(uid int, uid_dest int) int {
	var id int
	cf.DB.ShouldQueryRow("SELECT id FROM #DB#.friendships WHERE (uid1=? AND uid2=?) OR (uid2=? AND uid1=?)", uid, uid_dest, uid, uid_dest).Scan(&id)
	return id
}

func (cf *CFriendship) GetFriendByFID(id int) (int, int) {
	var (
		uid1 int
		uid2 int
	)
	cf.DB.ShouldQueryRow("SELECT uid1,uid2 FROM #DB#.friendships WHERE id=?", id).Scan(&uid1, &uid2)
	return uid1, uid2
}

func (cf *CFriendship) GetAccFriends(acc CAccount) []int {
	fr := strings.Split(acc.FriendshipIds, ",")
	var frlist []int
	for _, sfr := range fr {
		id, err := strconv.Atoi(sfr)
		if err != nil {
			continue
		}
		uid1, uid2 := cf.GetFriendByFID(id)
		targetUid := uid1
		if uid1 == acc.Uid {
			targetUid = uid2
		}
		frlist = append(frlist, targetUid)
	}
	return frlist
}

func (cf *CFriendship) ReadFriendRequest(id int) {
	cf.DB.ShouldExec("UPDATE #DB#.friendreqs SET isNew=0 WHERE id=?", id)
}

func (cf *CFriendship) RequestFriend(uid int, uidDest int, comment string) int {
	if uid == uidDest || cf.IsAlreadyFriend(uid, uidDest) || cf.IsAlreadySentFriend(uid, uidDest) || len(comment) > 512 {
		return -1
	}
	acc := CAccount{DB: cf.DB}
	acc.Uid = uidDest
	acc.LoadSettings()
	if acc.FrS > 0 {
		return -1
	}
	acc.LoadSocial()
	blacklist := strings.Split(acc.Blacklist, ",")
	if slices.Contains(blacklist, strconv.Itoa(uid)) {
		return -1
	}
	acc.DB.ShouldExec("INSERT INTO #DB#.friendreqs (uid_src, uid_dest, uploadDate, comment) VALUES (?,?,?,?)",
		uid, uidDest, time.Now().Format("2006-01-02 15:04:05"), comment)
	return 1
}

func (cf *CFriendship) AcceptFriendRequest(id int, uid int) int {
	var (
		src  int
		dest int
	)
	cf.DB.ShouldQueryRow("SELECT uid_src,uid_dest FROM #DB#.friendreqs WHERE id=?", id).Scan(&src, &dest)
	if src == dest || uid != dest {
		return -1
	}
	req, _ := cf.DB.PrepareExec("INSERT INTO #DB#.friendships (uid1, uid2) VALUES (?,?)", src, dest)
	iid, _ := req.LastInsertId()
	cf.DB.ShouldExec("DELETE FROM #DB#.friendreqs WHERE id=?", id)
	u1 := CAccount{DB: cf.DB}
	u2 := CAccount{DB: cf.DB}
	u1.Uid = src
	u2.Uid = dest
	res := u1.UpdateFriendships(CFRIENDSHIP_ADD, int(iid))
	res += u2.UpdateFriendships(CFRIENDSHIP_ADD, int(iid))
	if res != 2 {
		return -1
	}
	return 1
}

func (cf *CFriendship) RejectFriendRequestById(id int, uid int) int {
	var (
		uid1 int
		uid2 int
	)
	cf.DB.ShouldQueryRow("SELECT uid_src, uid_dest FROM #DB#.friendreqs WHERE id=?", id).Scan(&uid1, &uid2)
	if uid1 == uid2 || uid != uid2 {
		return -1
	}
	cf.DB.ShouldExec("DELETE FROM #DB#.friendreqs WHERE id=?", id)
	return 1
}

func (cf *CFriendship) RejectFriendRequestByUid(uid int, uid_dest int, isSender bool) {
	var (
		uid1 int
		uid2 int
	)
	if isSender {
		uid1 = uid
		uid2 = uid_dest
	} else {
		uid1 = uid_dest
		uid2 = uid
	}
	cf.DB.ShouldExec("DELETE FROM #DB#.friendreqs WHERE uid_src=? AND uid_dest=?", uid1, uid2)
}
