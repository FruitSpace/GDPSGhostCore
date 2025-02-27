package core

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	CAUTH_UID               int = 17
	CAUTH_UNAME             int = 26
	CAUTH_EMAIL             int = 35
	CBAN_BAN                int = 44
	CBAN_UNBAN              int = 53
	CBLACKLIST_BLOCK        int = 62
	CBLACKLIST_UNBLOCK      int = 71
	CFRIENDSHIP_ADD         int = 37
	CFRIENDSHIP_REMOVE      int = 38
	CLEADERBOARD_BY_CPOINTS int = 14
	CLEADERBOARD_BY_STARS   int = 15
	CLEADERBOARD_GLOBAL     int = 21
	CLEADERBOARD_FRIENDS    int = 22
	CREWARD_CHEST_BIG       int = 500
	CREWARD_CHEST_SMALL     int = 501
)

type CAccount struct {
	//Main/Auth
	Uid      int
	Uname    string
	Passhash string
	GjpHash  string
	Email    string
	Role_id  int
	IsBanned int

	//Stats
	Stars         int
	Diamonds      int
	Coins         int
	UCoins        int
	Demons        int
	CPoints       int
	Orbs          int
	Moons         int
	Special       int
	LvlsCompleted int

	//Technical
	RegDate         string
	AccessDate      string
	LastIP          string
	GameVer         string
	ExtraDataString string
	ExtraData       struct {
		DemonStats struct {
			Standard struct {
				Easy    int
				Medium  int
				Hard    int
				Insane  int
				Extreme int
			}
			Platformer struct {
				Easy    int
				Medium  int
				Hard    int
				Insane  int
				Extreme int
			}
			Weeklies  int
			Gauntlets int
		}
		StandardStats struct {
			Auto     int
			Easy     int
			Normal   int
			Hard     int
			Harder   int
			Insane   int
			Daily    int
			Gauntlet int
		}
		PlatformerStats struct {
			Auto   int
			Easy   int
			Normal int
			Hard   int
			Harder int
			Insane int
		}
	}

	//Social
	Blacklist     string
	FriendsCount  int
	FriendshipIds string

	//Vessels
	IconType       int
	ColorPrimary   int
	ColorSecondary int
	ColorGlow      int
	Cube           int
	Ship           int
	Ball           int
	Ufo            int
	Wave           int
	Robot          int
	Spider         int
	Swing          int
	Jetpack        int
	Trace          int
	Death          int

	//Chests
	ChestSmallCount int
	ChestSmallTime  int
	ChestBigCount   int
	ChestBigTime    int

	//Settings
	FrS     int
	CS      int
	MS      int
	Youtube string
	Twitch  string
	Twitter string

	DB *MySQLConn
}

func (acc *CAccount) CountUsers() int {
	var cnt int
	acc.DB.MustQueryRow("SELECT count(*) as cnt FROM #DB#.users").Scan(&cnt)
	return cnt
}

func (acc *CAccount) Exists(uid int) bool {
	var cnt int
	acc.DB.MustQueryRow("SELECT count(*) as cnt FROM #DB#.users WHERE uid=?", uid).Scan(&cnt)
	return cnt > 0
}

func (acc *CAccount) SearchUsers(sterm string) []int {
	var uids []int
	sterm = strings.ReplaceAll(sterm, "%", "")
	if _, err := strconv.Atoi(sterm); err != nil && len(sterm) < 3 {
		return uids
	}
	rows := acc.DB.ShouldQuery("SELECT uid FROM #DB#.users WHERE uid=? OR uname LIKE ? ORDER BY stars DESC LIMIT 10", sterm, "%"+sterm+"%")
	defer rows.Close()
	for rows.Next() {
		var uid int
		rows.Scan(&uid)
		uids = append(uids, uid)
	}
	return uids
}

func (acc *CAccount) LoadSettings() {
	var settings string
	acc.DB.MustQueryRow("SELECT settings FROM #DB#.users WHERE uid=?", acc.Uid).Scan(&settings)
	json.Unmarshal([]byte(settings), acc)
}

func (acc *CAccount) PushSettings() {
	data := map[string]interface{}{
		"frS": acc.FrS, "cS": acc.CS, "mS": acc.MS,
		"youtube": acc.Youtube, "twitch": acc.Twitch, "twitter": acc.Twitter}
	js, _ := json.Marshal(data)
	acc.DB.ShouldExec("UPDATE #DB#.users SET settings=? WHERE uid=?", string(js), acc.Uid)
}

func (acc *CAccount) LoadChests() {
	var chests string
	acc.DB.MustQueryRow("SELECT chests FROM #DB#.users WHERE uid=?", acc.Uid).Scan(&chests)
	var chst map[string]int
	json.Unmarshal([]byte(chests), &chst)
	acc.ChestSmallCount = chst["small_count"]
	acc.ChestBigCount = chst["big_count"]
	acc.ChestSmallTime = chst["small_time"]
	acc.ChestBigTime = chst["big_time"]
}

func (acc *CAccount) PushChests() {
	data := map[string]int{"small_count": acc.ChestSmallCount, "big_count": acc.ChestBigCount,
		"small_time": acc.ChestSmallTime, "big_time": acc.ChestBigTime}
	js, _ := json.Marshal(data)
	acc.DB.ShouldExec("UPDATE #DB#.users SET chests=? WHERE uid=?", string(js), acc.Uid)
}

func (acc *CAccount) LoadVessels() {
	var vessels string
	acc.DB.MustQueryRow("SELECT iconType,vessels FROM #DB#.users WHERE uid=?", acc.Uid).Scan(&acc.IconType, &vessels)
	json.Unmarshal([]byte(vessels), acc)
	var clrs map[string]int
	json.Unmarshal([]byte(vessels), &clrs)
	acc.ColorPrimary = clrs["clr_primary"]
	acc.ColorSecondary = clrs["clr_secondary"]
	acc.ColorGlow = clrs["clr_glow"]
}

func (acc *CAccount) PushVessels() {
	data := map[string]int{"clr_primary": acc.ColorPrimary, "clr_secondary": acc.ColorSecondary, "clr_glow": acc.ColorGlow,
		"cube": acc.Cube, "ship": acc.Ship, "ball": acc.Ball, "ufo": acc.Ufo, "wave": acc.Wave, "robot": acc.Robot,
		"spider": acc.Spider, "swing": acc.Swing, "jetpack": acc.Jetpack, "trace": acc.Trace, "death": acc.Death}
	js, _ := json.Marshal(data)
	acc.DB.ShouldExec("UPDATE #DB#.users SET vessels=?, iconType=? WHERE uid=?", string(js), acc.IconType, acc.Uid)
}

func (acc *CAccount) LoadStats() {
	acc.DB.MustQueryRow("SELECT stars,diamonds,coins,ucoins,demons,cpoints,orbs,moons,special,lvlsCompleted FROM #DB#.users WHERE uid=?", acc.Uid).Scan(
		&acc.Stars, &acc.Diamonds, &acc.Coins, &acc.UCoins, &acc.Demons, &acc.CPoints, &acc.Orbs, &acc.Moons, &acc.Special, &acc.LvlsCompleted)
}

func (acc *CAccount) PushStatsAndExtra() {
	sex, _ := json.Marshal(acc.ExtraData)
	acc.ExtraDataString = string(sex)
	acc.DB.ShouldExec("UPDATE #DB#.users SET stars=?,diamonds=?,coins=?,ucoins=?,demons=?,cpoints=?,orbs=?,moons=?,special=?,lvlsCompleted=?, extraData=? WHERE uid=?",
		acc.Stars, acc.Diamonds, acc.Coins, acc.UCoins, acc.Demons, acc.CPoints, acc.Orbs, acc.Moons, acc.Special, acc.LvlsCompleted, acc.ExtraDataString, acc.Uid)
}

func (acc *CAccount) PushExtra() {
	sex, _ := json.Marshal(acc.ExtraData)
	acc.ExtraDataString = string(sex)
	acc.DB.ShouldExec("UPDATE #DB#.users SET extraData=? WHERE uid=?", acc.ExtraDataString, acc.Uid)
}

func (acc *CAccount) LoadAuth(method int) {
	var req *sql.Row
	switch method {
	case CAUTH_UID:
		req = acc.DB.MustQueryRow("SELECT uid,uname,passhash,gjphash,email,role_id,isBanned FROM #DB#.users WHERE uid=?", acc.Uid)
	case CAUTH_UNAME:
		req = acc.DB.MustQueryRow("SELECT uid,uname,passhash,gjphash,email,role_id,isBanned FROM #DB#.users WHERE uname=?", acc.Uname)
	case CAUTH_EMAIL:
		req = acc.DB.MustQueryRow("SELECT uid,uname,passhash,gjphash,email,role_id,isBanned FROM #DB#.users WHERE email=?", acc.Email)
	default:
		return
	}
	req.Scan(&acc.Uid, &acc.Uname, &acc.Passhash, &acc.GjpHash, &acc.Email, &acc.Role_id, &acc.IsBanned)
}

func (acc *CAccount) LoadTechnical() {
	acc.Migrations()
	acc.DB.MustQueryRow("SELECT regDate,accessDate,lastIP,gameVer, extraData FROM #DB#.users WHERE uid=?", acc.Uid).Scan(
		&acc.RegDate, &acc.AccessDate, &acc.LastIP, &acc.GameVer, &acc.ExtraDataString)
	json.Unmarshal([]byte(acc.ExtraDataString), &acc.ExtraData)
}

func (acc *CAccount) LoadSocial() {
	acc.DB.MustQueryRow("SELECT blacklist,friends_cnt,friendship_ids FROM #DB#.users WHERE uid=?", acc.Uid).Scan(
		&acc.Blacklist, &acc.FriendsCount, &acc.FriendshipIds)
	acc.Blacklist = strings.ReplaceAll(acc.Blacklist, ",,", ",")
	acc.FriendshipIds = strings.ReplaceAll(acc.FriendshipIds, ",,", ",")
}

func (acc *CAccount) LoadAll() {
	acc.Migrations()
	var vessels, settings string
	acc.DB.MustQueryRow("SELECT uid,uname,passhash,gjphash,email,role_id,isBanned,stars,diamonds,coins,ucoins,"+
		"demons,cpoints,orbs,moons,special,lvlsCompleted,regDate,accessDate,lastIP,gameVer,blacklist,friends_cnt,friendship_ids,"+
		"iconType,vessels,settings, extraData FROM #DB#.users WHERE uid=?", acc.Uid).Scan(
		&acc.Uid, &acc.Uname, &acc.Passhash, &acc.GjpHash, &acc.Email, &acc.Role_id, &acc.IsBanned, &acc.Stars, &acc.Diamonds, &acc.Coins,
		&acc.UCoins, &acc.Demons, &acc.CPoints, &acc.Orbs, &acc.Moons, &acc.Special, &acc.LvlsCompleted, &acc.RegDate, &acc.AccessDate,
		&acc.LastIP, &acc.GameVer, &acc.Blacklist, &acc.FriendsCount, &acc.FriendshipIds, &acc.IconType, &vessels, &settings, &acc.ExtraDataString)
	json.Unmarshal([]byte(vessels), acc)
	var clrs map[string]int
	json.Unmarshal([]byte(vessels), &clrs)
	acc.ColorPrimary = clrs["clr_primary"]
	acc.ColorSecondary = clrs["clr_secondary"]
	acc.ColorGlow = clrs["clr_glow"]
	json.Unmarshal([]byte(settings), acc)
	acc.Blacklist = QuickComma(acc.Blacklist)
	acc.FriendshipIds = QuickComma(acc.FriendshipIds)
	json.Unmarshal([]byte(acc.ExtraDataString), &acc.ExtraData)
}

func (acc *CAccount) GetUIDByUname(uname string, autoSave bool) int {
	var uid int
	acc.DB.ShouldQueryRow("SELECT uid FROM #DB#.users WHERE uname=?", uname).Scan(&uid)
	if uid == 0 {
		return -1
	}
	if autoSave {
		acc.Uid = uid
	}
	return uid
}

func (acc *CAccount) GetUnameByUID(uid int) string {
	var uname string
	acc.DB.ShouldQueryRow("SELECT uname FROM #DB#.users WHERE uid=?", uid).Scan(&uname)
	if uname == "" {
		return "-1"
	}
	return uname
}

func (acc *CAccount) UpdateIP(ip string) {
	acc.LastIP = ip
	acc.DB.ShouldExec("UPDATE #DB#.users SET lastIP=?, accessDate=? WHERE uid=?", ip, time.Now().Format("2006-01-02 15:04:05"), acc.Uid)
}

func (acc *CAccount) UpdateGJP2(gjp2 string) {
	acc.GjpHash = gjp2
	acc.DB.ShouldExec("UPDATE #DB#.users SET gjphash=? WHERE uid=?", acc.GjpHash, acc.Uid)
}

func (acc *CAccount) CountIPs(ip string) int {
	var cnt int
	acc.DB.MustQueryRow("SELECT count(*) as cnt FROM #DB#.users WHERE lastIP=?", ip).Scan(&cnt)
	return cnt
}

func (acc *CAccount) UpdateBlacklist(action int, uid int) {
	acc.LoadSocial()
	blacklist := strings.Split(acc.Blacklist, ",")
	if action == CBLACKLIST_BLOCK && !slices.Contains(blacklist, strconv.Itoa(uid)) {
		blacklist = append(blacklist, strconv.Itoa(uid))
	}
	if action == CBLACKLIST_UNBLOCK && slices.Contains(blacklist, strconv.Itoa(uid)) {
		i := slices.Index(blacklist, strconv.Itoa(uid))
		blacklist = sliceRemove(blacklist, i)
	}
	acc.Blacklist = strings.Join(blacklist, ",")
	acc.DB.ShouldExec("UPDATE #DB#.users SET blacklist=? WHERE uid=?", acc.Blacklist, acc.Uid)
}

func (acc *CAccount) UpdateFriendships(action int, uid int) int {
	acc.LoadSocial()
	friendships := strings.Split(acc.FriendshipIds, ",")
	if action == CFRIENDSHIP_ADD && !slices.Contains(friendships, strconv.Itoa(uid)) {
		acc.FriendsCount++
		friendships = append(friendships, strconv.Itoa(uid))
	} else if action == CFRIENDSHIP_REMOVE && slices.Contains(friendships, strconv.Itoa(uid)) {
		acc.FriendsCount--
		i := slices.Index(friendships, strconv.Itoa(uid))
		friendships = sliceRemove(friendships, i)
	} else {
		return -1
	}
	acc.FriendshipIds = strings.Join(friendships, ",")
	acc.DB.ShouldExec("UPDATE #DB#.users SET friends_cnt=?, friendship_ids=? WHERE uid=?", acc.FriendsCount, acc.FriendshipIds, acc.Uid)
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
	case 7:
		return acc.Swing
	case 0:
	default:
		return acc.Cube
	}
	return acc.Cube
}

func (acc *CAccount) GetLeaderboardRank() int {
	var cnt int
	acc.DB.ShouldQueryRow("SELECT count(*) as cnt FROM #DB#.users WHERE stars>=? AND isBanned=0", acc.Stars).Scan(&cnt)
	return cnt
}

func (acc *CAccount) GetLeaderboard(atype int, grep []string, globalStars int, limit int) []int {
	var query string
	switch atype {
	case CLEADERBOARD_BY_STARS:
		query = "SELECT uid FROM #DB#.users WHERE stars>0 AND isBanned=0 ORDER BY stars DESC, uname ASC LIMIT " + strconv.Itoa(limit)
	case CLEADERBOARD_BY_CPOINTS:
		query = "SELECT uid FROM #DB#.users WHERE cpoints>0 AND isBanned=0 ORDER BY cpoints DESC, uname ASC LIMIT " + strconv.Itoa(limit)
	case CLEADERBOARD_GLOBAL:
		query = "SELECT X.uid as uid,X.stars FROM ((SELECT uid,stars,uname FROM #DB#.users WHERE stars>" + strconv.Itoa(globalStars) + " AND isBanned=0 ORDER BY stars ASC LIMIT 50)"
		query += " UNION (SELECT uid,stars,uname FROM #DB#.users WHERE stars<=" + strconv.Itoa(globalStars) + " AND stars>0 AND isBanned=0 ORDER BY stars DESC LIMIT 50)) as X ORDER BY X.stars DESC, X.uname ASC"
	case CLEADERBOARD_FRIENDS:
		friends := strings.Join(grep, ",")
		query = "SELECT uid FROM #DB#.users WHERE stars>0 AND isBanned=0 and uid IN (" + friends + ") ORDER BY stars DESC, uname ASC"
	default:
		query = "SELECT uid FROM #DB#.users WHERE 1=0" //IDK WHY I DID THIS
	}
	rows := acc.DB.ShouldQuery(query)
	defer rows.Close()

	var users []int
	for rows.Next() {
		var uid int
		if atype == CLEADERBOARD_GLOBAL {
			var stars int
			rows.Scan(&uid, &stars) //Workaround for Globals
		} else {
			rows.Scan(&uid)
		}
		users = append(users, uid)
	}
	return users
}

func (acc *CAccount) UpdateRole(role_id int) {
	acc.Role_id = role_id
	acc.DB.ShouldExec("UPDATE #DB#.users SET role_id=? WHERE uid=?", role_id, acc.Uid)
}

func (acc *CAccount) GetRoleObj(fetchPrivs bool) Role {
	role := Role{}
	if acc.Role_id == 0 {
		return role
	}
	if fetchPrivs {
		var privs string
		acc.DB.MustQueryRow("SELECT roleName,commentColor,modLevel,privs FROM #DB#.roles WHERE id=?", acc.Role_id).Scan(
			&role.RoleName, &role.CommentColor, &role.ModLevel, &privs)
		json.Unmarshal([]byte(privs), &role.Privs)
	} else {
		acc.DB.MustQueryRow("SELECT roleName,commentColor,modLevel FROM #DB#.roles WHERE id=?", acc.Role_id).Scan(
			&role.RoleName, &role.CommentColor, &role.ModLevel)
	}
	return role
}

func (acc *CAccount) UpdateAccessTime() {
	acc.DB.ShouldExec("UPDATE #DB#.users SET accessDate=?  WHERE uid=?", time.Now().Format("2006-01-02 15:04:05"), acc.Uid)
}

func (acc *CAccount) BanUser(action int) {
	var ban int
	switch action {
	case CBAN_BAN:
		ban = 2
	case CBAN_UNBAN:
		ban = 0
	default:
		ban = 1
	}
	acc.IsBanned = ban
	acc.DB.ShouldExec("UPDATE #DB#.users SET isBanned=? WHERE uid=?", ban, acc.Uid)
}

func (acc *CAccount) ChangePassword(passhash string) {
	acc.DB.ShouldExec("UPDATE #DB#.users SET passhash=? WHERE uid=?", passhash, acc.Uid)
	SendMessageDiscord(fmt.Sprintf("[%d] %s Switched algorithms. Size 36->%d", acc.Uid, acc.Uname, len(passhash)))
}

func (acc *CAccount) LogIn(uname string, pass string, ip string, uid int) int {
	if uid == 0 {
		uid = acc.GetUIDByUname(uname, false)
	}
	if uid > 0 {
		acc.Uid = uid
		acc.LoadAuth(CAUTH_UID)
		if acc.IsBanned > 0 {
			return -12
		}
		passx := SHA256(SHA512(pass) + "SaltyTruth:sob:")
		if len(acc.Passhash) == 36 {
			acc.ChangePassword(passx)
			passx = MD5(MD5(pass+"HalogenCore1704")+"ae07") + MD5(pass)[:4]
		}

		if acc.Passhash == passx {
			acc.UpdateIP(ip)
			return uid
		}
	}
	return -1
}

func (acc *CAccount) LogIn22(uname string, gjp string, ip string, uid int) int {
	if uid == 0 {
		uid = acc.GetUIDByUname(uname, false)
	}
	if uid > 0 {
		acc.Uid = uid
		acc.LoadAuth(CAUTH_UID)
		if acc.IsBanned > 0 {
			return -12
		}

		if acc.GjpHash == gjp {
			acc.UpdateIP(ip)
			return uid
		}
	}
	return -1
}

func (acc *CAccount) Register(uname string, pass string, email string, ip string, autoVerify bool) int {
	isBanned := "1"
	if autoVerify {
		isBanned = "0"
	}
	if len(uname) > 16 || !FilterEmail(email) {
		return -1
	}
	if acc.GetUIDByUname(uname, false) != -1 {
		return -2
	}
	var uid int
	acc.DB.ShouldQueryRow("SELECT uid FROM #DB#.users WHERE email=?", email).Scan(&uid)
	if uid != 0 {
		return -3
	}
	//passx := MD5(MD5(pass+"HalogenCore1704")+"ae07") + MD5(pass)[:4]
	passx := SHA256(SHA512(pass) + "SaltyTruth:sob:")

	rdate := time.Now().Format("2006-01-02 15:04:05")
	sreq := acc.DB.MustPrepareExec(
		"INSERT INTO #DB#.users (uname,passhash,gjphash,email,regDate,accessDate,isBanned) VALUES (?,?,?,?,?,?,?)",
		uname, passx, DoGjp2(pass), email, rdate, rdate, isBanned)
	vuid, _ := sreq.LastInsertId()
	acc.Uid = int(vuid)
	acc.UpdateIP(ip)
	return 1
}

func (acc *CAccount) VerifySession(uid int, ip string, gjp string, is22 bool) bool {
	var (
		aDate    string
		lastIP   string
		isBanned int
	)
	acc.DB.ShouldQueryRow("SELECT accessDate, lastIP, isBanned FROM #DB#.users WHERE uid=?", uid).Scan(&aDate, &lastIP, &isBanned)
	if aDate == "" || isBanned > 0 {
		acc.IsBanned = isBanned
		return false
	}
	ptime, _ := time.ParseInLocation("2006-01-02 15:04:05", aDate, loc)
	if ip == lastIP && (time.Now().Unix()-ptime.Unix()) < 3600 {
		acc.Uid = uid
		acc.LoadAuth(CAUTH_UID)
		return true
	}
	if is22 {
		acc.Uid = uid
		acc.LoadAuth(CAUTH_UID)
		if acc.IsBanned > 0 {
			return false
		}
		if acc.GjpHash == gjp {
			acc.UpdateIP(ip)
			return true
		}
	} else {
		gjp = strings.ReplaceAll(strings.ReplaceAll(gjp, "_", "/"), "-", "+")
		vgjp, _ := base64.StdEncoding.DecodeString(gjp)
		gjp = DoXOR(string(vgjp), "37526")
		if acc.LogIn("", gjp, ip, uid) > 0 {
			return true
		}
	}
	return false
}

func (acc *CAccount) PerformGJPAuth(Post url.Values, IPAddr string) bool {
	var uid int
	TryInt(&uid, Post.Get("accountID"))
	if GetGDVersion(Post) == 22 {
		gjp := ClearGDRequest(Post.Get("gjp2"))
		if !acc.VerifySession(uid, IPAddr, gjp, true) {
			return false
		}
	} else {
		gjp := ClearGDRequest(Post.Get("gjp"))
		if !acc.VerifySession(uid, IPAddr, gjp, false) {
			return false
		}
	}
	acc.Migrations()
	return true
}

func (acc *CAccount) Migrations() {
	acc.DB.ShouldExec("ALTER TABLE #DB#.users ADD COLUMN IF NOT EXISTS extraData LONGTEXT DEFAULT '{}'")
}

// Role
type Role struct {
	RoleName     string         `json:"role_name"`
	CommentColor string         `json:"comment_color"`
	ModLevel     int            `json:"mod_level"`
	Privs        map[string]int `json:"privileges,omitempty"`
}
