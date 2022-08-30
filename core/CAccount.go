package core

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"golang.org/x/exp/slices"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	CAUTH_UID int = 17
	CAUTH_UNAME int = 26
	CAUTH_EMAIL int = 35
	CBAN_BAN int =  44
	CBAN_UNBAN int =  53
	CBLACKLIST_BLOCK int =  62
	CBLACKLIST_UNBLOCK int = 71
	CFRIENDSHIP_ADD int =  37
	CFRIENDSHIP_REMOVE int =  38
	CLEADERBOARD_BY_CPOINTS int = 14
	CLEADERBOARD_BY_STARS int = 15
	CLEADERBOARD_GLOBAL int = 21
	CLEADERBOARD_FRIENDS int = 22
	CREWARD_CHEST_BIG int = 500
	CREWARD_CHEST_SMALL int = 501
)


type CAccount struct {
	//Main/Auth
	Uid int
	Uname string
	Passhash string
	GjpHash string
	Email string
	Role_id int
	IsBanned int

	//Stats
	Stars int
	Diamonds int
	Coins int
	UCoins int
	Demons int
	CPoints int
	Orbs int
	Moons int
	Special int
	LvlsCompleted int

	//Technical
	RegDate string
	AccessDate string
	LastIP string
	GameVer string

	//Social
	Blacklist string
	FriendsCount int
	FriendshipIds string

	//Vessels
	IconType int
	ColorPrimary int
	ColorSecondary int
	Cube int
	Ship int
	Ball int
	Ufo int
	Wave int
	Robot int
	Spider int
	Swing int
	Jetpack int
	Trace int
	Death int

	//Chests
	ChestSmallCount int
	ChestSmallTime int
	ChestBigCount int
	ChestBigTime int

	//Settings
	FrS int
	CS int
	MS int
	Youtube string
	Twitch string
	Twitter string

	DB MySQLConn
}


func (acc *CAccount) CountUsers() int {
	var cnt int
	acc.DB.MustQueryRow("SELECT count(*) as cnt FROM users").Scan(&cnt)
	return cnt
}

func (acc *CAccount) Exists(uid int) bool {
	var cnt int
	acc.DB.MustQueryRow("SELECT count(*) as cnt FROM users WHERE uid=?",uid).Scan(&cnt)
	return cnt>0
}

func (acc *CAccount) SearchUsers(sterm string) int {
	var uid int
	acc.DB.DB.QueryRow("SELECT uid FROM users WHERE uid=? OR uname=? ORDER BY stars LIMIT 1",sterm,sterm).Scan(&uid)
	return uid
}

func (acc *CAccount) LoadSettings() {
	var settings string
	acc.DB.MustQueryRow("SELECT settings FROM users WHERE uid=?",acc.Uid).Scan(&settings)
	json.Unmarshal([]byte(settings),acc)
}

func (acc *CAccount) PushSettings() {
	data:= map[string]interface{} {
		"frS": acc.FrS, "cS": acc.CS, "mS": acc.MS,
		"youtube": acc.Youtube, "twitch": acc.Twitch, "twitter": acc.Twitter}
	js,_:=json.Marshal(data)
	acc.DB.ShouldQuery("UPDATE users SET settings=? WHERE uid=?",string(js),acc.Uid)
}

func (acc *CAccount) LoadChests() {
	var chests string
	acc.DB.MustQueryRow("SELECT chests FROM users WHERE uid=?",acc.Uid).Scan(&chests)
	var chst map[string]int
	json.Unmarshal([]byte(chests),&chst)
	acc.ChestSmallCount=chst["small_count"]
	acc.ChestBigCount=chst["big_count"]
	acc.ChestSmallTime=chst["small_time"]
	acc.ChestBigTime=chst["big_time"]
}

func (acc *CAccount) PushChests() {
	data:= map[string]int {"small_count": acc.ChestSmallCount, "big_count": acc.ChestBigCount,
		"small_time": acc.ChestSmallTime, "big_time": acc.ChestBigTime}
	js,_:=json.Marshal(data)
	acc.DB.ShouldQuery("UPDATE users SET chests=? WHERE uid=?",string(js),acc.Uid)
}

func (acc *CAccount) LoadVessels() {
	var vessels string
	acc.DB.MustQueryRow("SELECT iconType,vessels FROM users WHERE uid=?",acc.Uid).Scan(&acc.IconType,&vessels)
	json.Unmarshal([]byte(vessels),acc)
	var clrs map[string]int
	json.Unmarshal([]byte(vessels),&clrs)
	acc.ColorPrimary=clrs["clr_primary"]
	acc.ColorSecondary=clrs["clr_secondary"]
}

func (acc *CAccount) PushVessels() {
	data:= map[string]int {"clr_primary": acc.ColorPrimary, "clr_secondary": acc.ColorSecondary, "cube": acc.Cube, "ship": acc.Ship,
		"ball": acc.Ball, "ufo": acc.Ufo, "wave": acc.Wave, "robot": acc.Robot, "spider": acc.Spider, "swing": acc.Swing,
		"jetpack": acc.Jetpack, "trace": acc.Trace, "death": acc.Death}
	js,_:=json.Marshal(data)
	acc.DB.DB.Query("UPDATE users SET vessels=?, iconType=? WHERE uid=?",string(js),acc.IconType,acc.Uid)
}

func (acc *CAccount) LoadStats() {
	acc.DB.MustQueryRow("SELECT stars,diamonds,coins,ucoins,demons,cpoints,orbs,moons,special,lvlsCompleted FROM users WHERE uid=?",acc.Uid).Scan(
		&acc.Stars,&acc.Diamonds,&acc.Coins,&acc.UCoins,&acc.Demons,&acc.CPoints,&acc.Orbs,&acc.Moons,&acc.Special,&acc.LvlsCompleted)
}

func (acc *CAccount) PushStats() {
	acc.DB.ShouldQuery("UPDATE users SET stars=?,diamonds=?,coins=?,ucoins=?,demons=?,cpoints=?,orbs=?,moons=?,special=?,lvlsCompleted=? WHERE uid=?",
		acc.Stars,acc.Diamonds,acc.Coins,acc.UCoins,acc.Demons,acc.CPoints,acc.Orbs,acc.Moons,acc.Special,acc.LvlsCompleted,acc.Uid)
}

func (acc *CAccount) LoadAuth(method int) {
	var req *sql.Row
	switch method {
	case CAUTH_UID:
		req=acc.DB.MustQueryRow("SELECT uid,uname,passhash,gjphash,email,role_id,isBanned FROM users WHERE uid=?",acc.Uid)
	case CAUTH_UNAME:
		req=acc.DB.MustQueryRow("SELECT uid,uname,passhash,gjphash,email,role_id,isBanned FROM users WHERE uname=?",acc.Uname)
	case CAUTH_EMAIL:
		req=acc.DB.MustQueryRow("SELECT uid,uname,passhash,gjphash,email,role_id,isBanned FROM users WHERE email=?",acc.Email)
	default:
		return
	}
	req.Scan(&acc.Uid,&acc.Uname,&acc.Passhash,&acc.GjpHash,&acc.Email,&acc.Role_id,&acc.IsBanned)
}

func (acc *CAccount) LoadTechnical() {
	acc.DB.MustQueryRow("SELECT regDate,accessDate,lastIP,gameVer FROM users WHERE uid=?",acc.Uid).Scan(
		&acc.RegDate,&acc.AccessDate,&acc.LastIP,&acc.GameVer)
}

func (acc *CAccount) LoadSocial() {
	acc.DB.MustQueryRow("SELECT blacklist,friends_cnt,friendship_ids FROM users WHERE uid=?",acc.Uid).Scan(
		&acc.Blacklist,&acc.FriendsCount,&acc.FriendshipIds)
	acc.Blacklist=strings.ReplaceAll(acc.Blacklist,",,",",")
	acc.FriendshipIds=strings.ReplaceAll(acc.FriendshipIds,",,",",")
}

func (acc *CAccount) LoadAll() {
	acc.LoadAuth(CAUTH_UID)
	acc.LoadStats()
	acc.LoadTechnical()
	acc.LoadSocial()
	acc.LoadVessels()
	acc.LoadSettings()
}

func (acc *CAccount) GetUIDByUname(uname string, autoSave bool) int {
	var uid int
	acc.DB.DB.QueryRow("SELECT uid FROM users WHERE uname=?",uname).Scan(&uid)
	if uid==0 {return -1}
	if autoSave {acc.Uid=uid}
	return uid
}

func (acc *CAccount) GetUnameByUID(uid int) string {
	var uname string
	acc.DB.DB.QueryRow("SELECT uname FROM users WHERE uid=?",uid).Scan(&uname)
	if uname=="" {return "-1"}
	return uname
}

func (acc *CAccount) UpdateIP(ip string) {
	acc.LastIP=ip
	acc.DB.ShouldQuery("UPDATE users SET lastIP=?, accessDate=? WHERE uid=?",ip,time.Now().Format("2006-01-02 15:04:05"),acc.Uid)
}

func (acc *CAccount) UpdateGJP2(gjp2 string) {
	acc.GjpHash=gjp2
	acc.DB.ShouldQuery("UPDATE users SET gjphash=? WHERE uid=?",acc.GjpHash,acc.Uid)
}

func (acc *CAccount) CountIPs(ip string) int {
	var cnt int
	acc.DB.MustQueryRow("SELECT count(*) as cnt FROM users WHERE lastIP=?",ip).Scan(&cnt)
	return cnt
}

func (acc *CAccount) UpdateBlacklist(action int, uid int) {
	acc.LoadSocial()
	blacklist:=strings.Split(acc.Blacklist,",")
	if action==CBLACKLIST_BLOCK && !slices.Contains(blacklist,strconv.Itoa(uid)) {blacklist=append(blacklist, strconv.Itoa(uid))}
	if action==CBLACKLIST_UNBLOCK && slices.Contains(blacklist,strconv.Itoa(uid)) {
		i:=slices.Index(blacklist,strconv.Itoa(uid))
		blacklist=sliceRemove(blacklist,i)
	}
	acc.Blacklist=strings.Join(blacklist,",")
	acc.DB.ShouldQuery("UPDATE users SET blacklist=? WHERE uid=?",acc.Blacklist,acc.Uid)
}

func (acc *CAccount) UpdateFriendships(action int, uid int) int{
	acc.LoadSocial()
	friendships:=strings.Split(acc.FriendshipIds,",")
	if action==CFRIENDSHIP_ADD && !slices.Contains(friendships,strconv.Itoa(uid)) {
		acc.FriendsCount++
		friendships=append(friendships,strconv.Itoa(uid))
	} else if action==CFRIENDSHIP_REMOVE && slices.Contains(friendships,strconv.Itoa(uid)) {
		acc.FriendsCount--
		i:=slices.Index(friendships,strconv.Itoa(uid))
		friendships=sliceRemove(friendships,i)
	} else {return -1}
	acc.FriendshipIds=strings.Join(friendships,",")
	acc.DB.ShouldQuery("UPDATE users SET friends_cnt=?, friendship_ids=? WHERE uid=?",acc.FriendsCount,acc.FriendshipIds,acc.Uid)
	return 1
}

func (acc *CAccount) GetShownIcon() int {
	switch acc.IconType {
	case 1:
		return acc.Ship
	case 2:
		return acc.Ball
	case 3:
		return acc.Ufo
	case 4:
		return acc.Wave
	case 5:
		return acc.Robot
	case 6:
		return acc.Spider
	case 0:
	default:
		return acc.Cube
	}
	return acc.Cube
}

func (acc *CAccount) GetLeaderboardRank() int {
	var cnt int
	acc.DB.ShouldQueryRow("SELECT count(*) as cnt FROM users WHERE stars>=? AND isBanned=0",acc.Stars).Scan(&cnt)
	return cnt
}

func (acc *CAccount) GetLeaderboard(atype int, grep []string, globalStars int, limit int) []int {
	var query string
	switch atype {
	case CLEADERBOARD_BY_STARS:
		query="SELECT uid FROM users WHERE stars>0 AND isBanned=0 ORDER BY stars DESC, uname ASC LIMIT "+strconv.Itoa(limit)
	case CLEADERBOARD_BY_CPOINTS:
		query="SELECT uid FROM users WHERE cpoints>0 AND isBanned=0 ORDER BY cpoints DESC, uname ASC LIMIT "+strconv.Itoa(limit)
	case CLEADERBOARD_GLOBAL:
		query="SELECT X.uid as uid,X.stars FROM ((SELECT uid,stars,uname FROM users WHERE stars>"+ strconv.Itoa(globalStars) +" AND isBanned=0 ORDER BY stars ASC LIMIT 50)"
		query+=" UNION (SELECT uid,stars,uname FROM users WHERE stars<="+ strconv.Itoa(globalStars) +" AND stars>0 AND isBanned=0 ORDER BY stars DESC LIMIT 50)) as X ORDER BY X.stars DESC, X.uname ASC"
	case CLEADERBOARD_FRIENDS:
		friends:=strings.Join(grep,",")
		query="SELECT uid FROM users WHERE stars>0 AND isBanned=0 and uid IN ("+friends+") ORDER BY stars DESC, uname ASC";
	default:
		query="SELECT uid FROM users WHERE 1=0" //IDK WHY I DID THIS
	}
	rows:=acc.DB.ShouldQuery(query)

	var users []int
	for rows.Next() {
		var uid int
		if atype==CLEADERBOARD_GLOBAL {
			var stars int
			rows.Scan(&uid,&stars) //Workaround for Globals
		}else {
			rows.Scan(&uid)
		}
		users=append(users,uid)
	}
	return users
}

func (acc *CAccount) UpdateRole(role_id int) {
	acc.Role_id = role_id
	acc.DB.ShouldQuery("UPDATE users SET role_id=? WHERE uid=?",role_id,acc.Uid)
}

func (acc *CAccount) GetRoleObj(fetchPrivs bool) Role {
	role:=Role{}
	if acc.Role_id==0 {return role}
	if fetchPrivs {
		var privs string
		acc.DB.MustQueryRow("SELECT roleName,commentColor,modLevel,privs FROM roles WHERE id=?",acc.Role_id).Scan(
			&role.RoleName,&role.CommentColor,&role.ModLevel,&privs)
		json.Unmarshal([]byte(privs),&role.Privs)
	}else{
		acc.DB.MustQueryRow("SELECT roleName,commentColor,modLevel FROM roles WHERE id=?",acc.Role_id).Scan(
			&role.RoleName,&role.CommentColor,&role.ModLevel)
	}
	return role
}

func (acc *CAccount) UpdateAccessTime() {
	acc.DB.ShouldQuery("UPDATE users SET accessDate=?  WHERE uid=?",time.Now().Format("2006-01-02 15:04:05"),acc.Uid)
}

func (acc *CAccount) BanUser(action int) {
	var ban int
	switch action{
	case CBAN_BAN:
		ban=2
	case CBAN_UNBAN:
		ban=0
	default:
		ban=1
	}
	acc.IsBanned=ban
	acc.DB.ShouldQuery("UPDATE users SET isBanned=? WHERE uid=?",ban,acc.Uid)
}

func (acc *CAccount) LogIn(uname string, pass string, ip string, uid int) int {
	if uid==0 {
		uid = acc.GetUIDByUname(uname, false)
	}
	if uid>0 {
		acc.Uid=uid
		acc.LoadAuth(CAUTH_UID)
		if acc.IsBanned>0 {return -12}
		passx:=MD5(MD5(pass+"HalogenCore1704")+"ae07")+MD5(pass)[:4]
		if acc.Passhash==passx {
			acc.UpdateIP(ip)
			return uid
		}
	}
	return -1
}

func (acc *CAccount) Register(uname string, pass string, email string, ip string) int {
	if len(uname)>16 || !FilterEmail(email) {return -1}
	if acc.GetUIDByUname(uname,false)!=-1 {return -2}
	var uid int
	acc.DB.DB.QueryRow("SELECT uid FROM users WHERE email=?",email).Scan(&uid)
	if uid!=0 {return -3}
	passx:=MD5(MD5(pass+"HalogenCore1704")+"ae07")+MD5(pass)[:4]
	rdate:=time.Now().Format("2006-01-02 15:04:05")
	sreq:=acc.DB.MustPrepareExec("INSERT INTO users (uname,passhash,gjphash,email,regDate,accessDate,isBanned) VALUES (?,?,?,?,?,?,1)",
		uname,passx,DoGjp2(pass),email,rdate,rdate)
	vuid,_:=sreq.LastInsertId()
	acc.Uid=int(vuid)
	acc.UpdateIP(ip)
	return 1
}

func (acc *CAccount) VerifySession(uid int, ip string, gjp string, is22 bool) bool {
	var (
		aDate string
		lastIP string
		isBanned int
	)
	acc.DB.ShouldQueryRow("SELECT accessDate, lastIP, isBanned FROM users WHERE uid=?",uid).Scan(&aDate,&lastIP,&isBanned)
	if aDate=="" || isBanned>0 {
		acc.IsBanned=isBanned
		return false
	}
	ptime,_:=time.Parse("2006-01-02 15:04:05",aDate)
	if ip==lastIP && (time.Now().Unix()-ptime.Unix())<3600 {
		acc.Uid=uid
		acc.LoadAuth(CAUTH_UID)
		return true
	}
	if is22 {
		acc.Uid=uid
		acc.LoadAuth(CAUTH_UID)
		if acc.IsBanned>0 {return false}
		if acc.GjpHash==gjp {
			acc.UpdateIP(ip)
			return true
		}
	}else{
		gjp=strings.ReplaceAll(strings.ReplaceAll(gjp,"_","/"),"-","+")
		vgjp,_:=base64.StdEncoding.DecodeString(gjp)
		gjp=DoXOR(string(vgjp),"37526")
		if acc.LogIn("",gjp,ip,uid)>0 {return true}
	}
	return false
}

func (acc *CAccount) PerformGJPAuth(Post url.Values, IPAddr string) bool {
	var uid int
	TryInt(&uid,Post.Get("accountID"))
	if GetGDVersion(Post)==22{
		gjp:=ClearGDRequest(Post.Get("gjp2"))
		if !acc.VerifySession(uid,IPAddr,gjp,true) {return false}
	}else{
		gjp:=ClearGDRequest(Post.Get("gjp"))
		if !acc.VerifySession(uid,IPAddr,gjp,false) {return false}
	}
	return true
}

//Role
type Role struct {
	RoleName string
	CommentColor string
	ModLevel int
	Privs map[string]int
}