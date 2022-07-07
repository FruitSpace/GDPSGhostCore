package core

import (
	"strconv"
	"time"
)

type CComment struct {
	Id int
	Uid int
	Likes int
	PostedTime string
	Comment string

	LvlId int
	Percent int
	IsSpam bool

	DB MySQLConn
}

func (cc *CComment) ExistsLevelComment(id int) bool {
	var cnt int
	cc.DB.ShouldQueryRow("SELECT count(*) as cnt FROM comments WHERE id=?",id).Scan(&cnt)
	return cnt>0
}

func (cc *CComment) ExistsAccComment(id int) bool {
	var cnt int
	cc.DB.ShouldQueryRow("SELECT count(*) as cnt FROM acccomments WHERE id=?",id).Scan(&cnt)
	return cnt>0
}

func (cc *CComment) CountAccComments(uid int) int {
	var cnt int
	pf:=""
	if uid>0 {pf=" WHERE uid="+strconv.Itoa(uid)}
	cc.DB.ShouldQueryRow("SELECT count(*) as cnt FROM acccomments"+pf).Scan(&cnt)
	return cnt
}

func (cc *CComment) CountLevelComments(lvlId int) int {
	var cnt int
	pf:=""
	if lvlId>0 {pf=" WHERE lvl_id="+strconv.Itoa(lvlId)}
	cc.DB.ShouldQueryRow("SELECT count(*) as cnt FROM comments"+pf).Scan(&cnt)
	return cnt
}

func (cc *CComment) CountCommentHistory(uid int) int {
	var cnt int
	cc.DB.ShouldQueryRow("SELECT count(*) as cnt FROM comments WHERE uid=?",uid).Scan(&cnt)
	return cnt
}

func (cc *CComment) LoadAccComment() {
	cc.DB.ShouldQueryRow("SELECT uid,comment,postedTime,likes,isSpam FROM acccomments WHERE id=?",cc.Id).Scan(
		&cc.Uid,&cc.Comment,&cc.PostedTime,&cc.Likes,&cc.IsSpam)
}

func (cc *CComment) LoadLevelComment() {
	cc.DB.ShouldQueryRow("SELECT uid,lvl_id,comment,postedTime,likes,isSpam,percent FROM comments WHERE id=?",cc.Id).Scan(
		&cc.Uid,&cc.LvlId,&cc.Comment,&cc.PostedTime,&cc.Likes,&cc.IsSpam,&cc.Percent)
}

func (cc *CComment) GetAllAccComments(uid int, page int) []CComment {
	page*=10
	rows:=cc.DB.ShouldQuery("SELECT id,comment,postedTime,likes,isSpam FROM acccomments WHERE uid=? ORDER BY postedTime DESC LIMIT 10 OFFSET "+strconv.Itoa(page),uid)
	var out []CComment
	for rows.Next() {
		comm:=CComment{Uid: uid}
		rows.Scan(&comm.Id,&comm.Comment,&comm.PostedTime,&comm.Likes,&comm.IsSpam)
		out=append(out,comm)
	}
	return out
}

func (cc *CComment) GetAllLevelComments(lvlId int, page int, sortMode bool) []CComment {
	filter:="postedTime"
	page*=10
	if sortMode {filter="likes"}
	rows:=cc.DB.ShouldQuery("SELECT id,uid,comment,postedTime,likes,isSpam,percent FROM comments WHERE lvl_id=? ORDER BY "+filter+" DESC LIMIT 10 OFFSET "+strconv.Itoa(page),lvlId)
	var out []CComment
	for rows.Next() {
		comm:=CComment{LvlId: lvlId}
		rows.Scan(&comm.Id,&comm.Uid,&comm.Comment,&comm.PostedTime,&comm.Likes,&comm.IsSpam,&comm.Percent)
		out=append(out,comm)
	}
	return out
}

// GetAllCommentsHistory Warns that page is not multiplied
func (cc *CComment) GetAllCommentsHistory(uid int, page int, sortMode bool) []CComment {
	filter:="postedTime"
	if sortMode {filter="likes"}
	rows:=cc.DB.ShouldQuery("SELECT id,lvl_id,comment,postedTime,likes,isSpam,percent FROM comments WHERE uid=? ORDER BY "+filter+" DESC LIMIT 10 OFFSET "+strconv.Itoa(page),uid)
	var out []CComment
	for rows.Next() {
		comm:=CComment{Uid: uid}
		rows.Scan(&comm.Id,&comm.LvlId,&comm.Comment,&comm.PostedTime,&comm.Likes,&comm.IsSpam,&comm.Percent)
		out=append(out,comm)
	}
	return out
}

func (cc *CComment) PostAccComment() bool{
	if len(cc.Comment)>128 {return false}
	cc.DB.ShouldQuery("INSERT INTO acccomments (uid,comment,postedTime) VALUES (?,?,?)",cc.Uid,cc.Comment,
		time.Now().Format("2006-01-02 15:04:05"))
	return true
}

func (cc *CComment) PostLevelComment() bool{
	if len(cc.Comment)>128 {return false}
	cc.DB.ShouldQuery("INSERT INTO comments (uid,lvl_id,comment,postedTime,percent) VALUES (?,?,?,?,?)",cc.Uid,
		cc.LvlId,cc.Comment,time.Now().Format("2006-01-02 15:04:05"),cc.Percent)
	return true
}

func (cc *CComment) DeleteAccComment(id int, uid int) {
	cc.DB.ShouldQuery("DELETE FROM acccomments WHERE id=? AND uid=?",id,uid)
}

func (cc *CComment) DeleteLevelComment(id int, uid int) {
	cc.DB.ShouldQuery("DELETE FROM comments WHERE id=? AND uid=?",id,uid)
}

func (cc *CComment) DeleteOwnerLevelComment(id int, lvlId int) {
	cc.DB.ShouldQuery("DELETE FROM comments WHERE id=? AND lvl_id=?",id,lvlId)
}

func (cc *CComment) LikeAccComment(id int, uid int, actionLike bool) bool {
	if IsLiked(ITEMTYPE_ACCCOMMENT,uid,id,cc.DB) {return false}
	operator:="-"
	actionc:="Dislike"
	if actionLike {
		operator="+"
		actionc="Like"
	}
	cc.DB.ShouldQuery("UPDATE acccomments SET likes=likes"+operator+"1 WHERE id=?",id)
	RegisterAction(ACTION_ACCCOMMENT_LIKE,uid,id, map[string]string{"type":actionc},cc.DB)
	return true
}

func (cc *CComment) LikeLevelComment(id int, uid int, actionLike bool) bool {
	if IsLiked(ITEMTYPE_COMMENT,uid,id,cc.DB) {return false}
	operator:="-"
	actionc:="Dislike"
	if actionLike {
		operator="+"
		actionc="Like"
	}
	cc.DB.ShouldQuery("UPDATE comments SET likes=likes"+operator+"1 WHERE id=?",id)
	RegisterAction(ACTION_COMMENT_LIKE,uid,id, map[string]string{"type":actionc},cc.DB)
	return true
}