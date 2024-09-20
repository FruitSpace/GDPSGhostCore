// Package connectors allow translating beautiful typed data to a hell of a mess RobTop format
// and also to communicate with outside world
package connectors

import (
	"HalogenGhostCore/core"
	"encoding/base64"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type GDConnector struct {
	output string
}

func (c *GDConnector) Output() string {
	log.Println("Output: " + c.output)
	return c.output
}

func (c *GDConnector) Error(code string, reason string) {
	log.Println("Error: " + code + " " + reason)
	c.output = code
}

func (c *GDConnector) Success(message string) {
	c.output = "1"
}

func (c *GDConnector) NumberedSuccess(id int) {
	c.output = strconv.Itoa(id)
}

func (c *GDConnector) Account_Sync(savedata string) {
	c.output = savedata + ";21;30;a;a"
}

func (c *GDConnector) Account_Login(uid int) {
	c.output = fmt.Sprintf("%d,%d", uid, uid)
}

func (c *GDConnector) Comment_AccountGet(comments []core.CComment, count int, page int) {
	if len(comments) == 0 {
		c.output = "#0:0:0"
	} else {
		for _, comm := range comments {
			c.output += c.getAccountComment(comm)
		}
		c.output = fmt.Sprintf("%s#%d:%d:10", c.output[:len(c.output)-1], count, page*10)
	}
}

func (c *GDConnector) Comment_LevelGet(comments []core.CComment, count int, page int) {
	if len(comments) == 0 {
		c.output = "#0:0:0"
	} else {
		for _, comm := range comments {
			c.output += c.getLevelComment(comm)
		}
		c.output = fmt.Sprintf("%s#%d:%d:10", c.output[:len(c.output)-1], count, page*10)
	}
}

func (c *GDConnector) Comment_HistoryGet(comments []core.CComment, acc core.CAccount, role core.Role, count int, page int) {
	if len(comments) == 0 {
		c.output = "#0:0:0"
	} else {
		for _, comm := range comments {
			c.output += c.getCommentHistory(comm, acc, role)
		}
		c.output = fmt.Sprintf("%s#%d:%d:10", c.output[:len(c.output)-1], count, page*10)
	}
}

func (c *GDConnector) Communication_FriendGetRequests(reqs []map[string]string, count int, page int) {
	for _, frq := range reqs {
		c.output += c.getFriendRequest(frq)
	}
	c.output = fmt.Sprintf("%s#%d:%d:10", c.output[:len(c.output)-1], count, page*10)
}

func (c *GDConnector) Communication_MessageGet(message core.CMessage, uid int) {
	c.output = c.getMessage(message, uid)
}

func (c *GDConnector) Communication_MessageGetAll(messages []map[string]string, getSent bool, count int, page int) {
	for _, msg := range messages {
		c.output += c.getMessageStr(msg, getSent)
	}
	c.output = fmt.Sprintf("%s#%d:%d:10", c.output[:len(c.output)-1], count, page*10)
}

// GetMusic used to get simple music string (w/o traling hash)
func (c *GDConnector) Essential_GetMusic(mus core.CMusic) {
	//convert size to string
	size := mus.Size
	size = math.Round(size*100) / 100

	//convert size to string
	sizeStr := strconv.FormatFloat(size, 'f', 2, 64)
	mstr := "1~|~" + strconv.Itoa(mus.Id) + "~|~2~|~" + mus.Name + "~|~3~|~1~|~4~|~" + mus.Artist + "~|~5~|~" + sizeStr + "~|~6~|~~|~10~|~" +
		url.QueryEscape(mus.Url)
	c.output = strings.ReplaceAll(mstr, "#", "")
}

// used to get simple top artists string (w/o trailing hash)
func (c *GDConnector) Essential_GetTopArtists(artists map[string]string) {
	for artist, youtube := range artists {
		c.output += "4:" + artist + ":7:" + youtube + "|"
	}
	c.output = c.output[:len(c.output)-1] + "#" + strconv.Itoa(len(artists)) + ":0:10"
}

func (c *GDConnector) Level_GetGauntlets(gaus []map[string]string, hash string) {
	for _, gau := range gaus {
		c.output += "1:" + gau["pack_name"] + ":3:" + gau["levels"] + "|"
	}
	c.output = c.output[:len(c.output)-1] + "#" + hash
}

func (c *GDConnector) Level_SearchList(intlists []int, lists []core.CLevelList, count int, page int) {
	lvlHash := ""
	usrstring := ""
	var llists []core.CLevelList
	for _, lid := range intlists {
		for i, list := range lists {
			if list.ID == lid {
				llists = append(llists, list)
				lists = append(lists[:i], lists[i+1:]...)
				break
			}
		}
	}
	for _, list := range llists {
		lvlS, usrH, lvlH := c.getListSearch(list)
		c.output += lvlS
		lvlHash += lvlH
		usrstring += usrH
	}

	if len(c.output) == 0 {
		c.output = "x"
		usrstring = "x"
	}

	c.output = c.output[:len(c.output)-1] + "#" +
		usrstring[:len(usrstring)-1] + "#" +
		strconv.Itoa(count) + ":" + strconv.Itoa(page*10) + ":10#" +
		core.HashSolo2(lvlHash)
}

func (c *GDConnector) Level_GetMapPacks(packs []core.LevelPack, count int, page int) {
	hashstr := ""
	for _, pack := range packs {
		c.output += fmt.Sprintf("1:%d:2:%s:3:%s:4:%d:5:%d:6:%d:7:%s:8:%s|",
			pack.Id, pack.PackName, pack.Levels, pack.PackStars, pack.PackCoins, pack.PackDifficulty,
			pack.PackColor, pack.PackColor)
		id_s := strconv.Itoa(pack.Id)
		hashstr += fmt.Sprintf("%s%s%d%d", string(id_s[0]), string(id_s[len(id_s)-1]), pack.PackStars, pack.PackCoins)

	}
	c.output = c.output[:len(c.output)-1] + "#" + strconv.Itoa(count) + ":" + strconv.Itoa(page*10) + ":10#" + core.HashSolo2(hashstr)
}

// Levels_GetLevelFull used to retrieve full Level data (w/ trailing hash)
func (c *GDConnector) Level_GetLevelFull(cl core.CLevel, password string, phash string, quest_id int) {
	s := strconv.Itoa
	t, err := time.ParseInLocation("2006-01-02 15:04:05", cl.UploadDate, loc)
	if err != nil {
		t = time.Now()
	}
	uplAge := core.GetDateAgo(t.Unix())
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", cl.UpdateDate, loc)
	if err != nil {
		t2 = time.Now()
	}
	updAge := core.GetDateAgo(t2.Unix())
	diffNom := 0
	if cl.Difficulty > 0 {
		diffNom = 10
	}
	var auto int
	if cl.Difficulty < 0 {
		auto = 1
		cl.Difficulty = 0
	}
	coinsVer := 0
	if cl.Coins > 0 {
		coinsVer = 1
	}
	demonDiff := 3
	isDemon := 0
	if cl.DemonDifficulty >= 0 {
		isDemon = 1
		demonDiff = cl.DemonDifficulty
	}
	quest := ""
	questHash := ""
	if quest_id > 0 {
		quest = ":41:" + s(quest_id)
		acc := core.CAccount{DB: cl.DB, Uid: cl.Uid}
		acc.LoadAuth(core.CAUTH_UID)
		questHash = "#" + s(acc.Uid) + ":" + acc.Uname + ":" + s(acc.Uid)
	}
	sfxSongs := strings.Split(cl.StringSettings, ";")
	if len(sfxSongs) == 1 {
		sfxSongs = append(sfxSongs, "")
	}
	hash := s(cl.Uid) + "," + s(cl.StarsGot) + "," + s(isDemon) + "," + s(cl.Id) + "," + s(coinsVer) + "," + s(cl.IsFeatured) + "," + phash +
		"," + s(quest_id)
	c.output = "1:" + s(cl.Id) + ":2:" + cl.Name + ":3:" + cl.Description + ":4:" + cl.StringLevel + ":5:" + s(cl.Version) + ":6:" + s(cl.Uid) + ":8:" + s(diffNom) +
		":9:" + s(cl.Difficulty) + ":10:" + s(cl.Downloads) + ":12:" + s(cl.TrackId) + ":13:" + s(cl.VersionGame) + ":14:" + s(cl.Likes) +
		":15:" + s(cl.Length) + ":17:" + s(isDemon) + ":18:" + s(cl.StarsGot) + ":19:" + s(cl.IsFeatured) + ":25:" + s(auto) + ":26:" + cl.StringLevelInfo +
		":27:" + password + ":28:" + uplAge + ":29:" + updAge + ":30:" + s(cl.OrigId) + ":31:" + s(core.ToInt(cl.Is2p)) + ":35:" + s(cl.SongId) +
		":36:" + cl.StringExtra + ":37:" + s(cl.Ucoins) + ":38:" + s(coinsVer) + ":39:" + s(cl.StarsRequested) + ":40:" + s(core.ToInt(cl.IsLDM)) +
		":42:" + s(cl.IsEpic) + ":43:" + s(demonDiff) + ":45:" + s(cl.Objects) + ":46:1:47:2:48::52:" + sfxSongs[0] + ":53:" + sfxSongs[1] + quest +
		"#" + core.HashSolo(cl.StringLevel) + "#" + core.HashSolo2(hash) + questHash

	//44 isGauntlet
}

func (c *GDConnector) Level_GetSpecials(id int, left int) {
	c.output = strconv.Itoa(id) + "|" + strconv.Itoa(left)
}

func (c *GDConnector) Level_SearchLevels(
	intlevels []int, levels []core.CLevel, mus *core.CMusic,
	count int, page int, gdVersion int, gauntlet int,
) {
	lvlHash := ""
	usrstring := ""
	musStr := ""
	var musQueue []int

	var lvls []core.CLevel
	for _, lvlid := range intlevels {
		for i, lvl := range levels {
			if lvl.Id == lvlid {
				lvls = append(lvls, lvl)
				levels = append(levels[:i], levels[i+1:]...)
				break
			}
		}
	}
	for _, lvl := range lvls {
		if gdVersion == 22 {
			lvl.VersionGame = 21
		}
		lvlS, lvlH, usrH := c.getLevelSearch(lvl, gauntlet != 0)
		c.output += lvlS
		lvlHash += lvlH
		usrstring += usrH

		if lvl.SongId != 0 {
			musQueue = append(musQueue, lvl.SongId)
		}
	}
	if len(musQueue) > 0 {
		songs := mus.GetBulkSongs(musQueue)
		for _, sng := range songs {
			mc := GDConnector{}
			mc.Essential_GetMusic(sng)
			musStr += mc.Output() + "~:~"
		}
	}

	if len(musStr) == 0 {
		musStr = "lll"
	}

	if len(c.output) == 0 {
		c.output = "x"
		usrstring = "x"
	}

	c.output = fmt.Sprintf(
		"%s#%s#%s#%d:%d:10#%s",
		c.output[:len(c.output)-1],
		usrstring[:len(usrstring)-1],
		musStr[:len(musStr)-3],
		count, page*10, core.HashSolo2(lvlHash),
	)
}

// Rewards_ChallengesOutput used to retrieve all quests/challenges data (w/ trailing hash)
func (c *GDConnector) Rewards_ChallengesOutput(cq core.CQuests, uid int, chk string, udid string) {
	s := strconv.Itoa
	virt := core.RandStringBytes(5)
	tme, _ := time.ParseInLocation("2006-01-02 15:04:05", strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")[0]+" 00:00:00", loc)
	//!Additional 10800 Review is needed
	timeLeft := int(tme.AddDate(0, 0, 1).Unix() - (time.Now().Unix()))
	out := virt + ":" + s(uid) + ":" + chk + ":" + udid + ":" + s(uid) + ":" + s(timeLeft) + ":" + cq.GetQuests(uid)
	out = strings.ReplaceAll(strings.ReplaceAll(base64.StdEncoding.EncodeToString([]byte(core.DoXOR(out, "19847"))), "/", "_"), "+", "-")
	c.output = virt + out + "|" + core.HashSolo3(out)
}

// Rewards_ChestOutput used to retrieve all chest data (w/ trailing hash)
func (c *GDConnector) Rewards_ChestOutput(acc core.CAccount, config core.ConfigBlob, udid string, chk string, smallLeft int, bigLeft int, chestType int) {
	s := strconv.Itoa
	config.ChestConfig.ChestSmallOrbsMax = core.MaxInt(config.ChestConfig.ChestSmallOrbsMax, config.ChestConfig.ChestSmallOrbsMin)
	config.ChestConfig.ChestSmallDiamondsMax = core.MaxInt(config.ChestConfig.ChestSmallDiamondsMax, config.ChestConfig.ChestSmallDiamondsMin)
	config.ChestConfig.ChestSmallKeysMax = core.MaxInt(config.ChestConfig.ChestSmallKeysMax, config.ChestConfig.ChestSmallKeysMin)

	config.ChestConfig.ChestBigOrbsMax = core.MaxInt(config.ChestConfig.ChestBigOrbsMax, config.ChestConfig.ChestBigOrbsMin)
	config.ChestConfig.ChestBigDiamondsMax = core.MaxInt(config.ChestConfig.ChestBigDiamondsMax, config.ChestConfig.ChestBigDiamondsMin)
	config.ChestConfig.ChestBigKeysMax = core.MaxInt(config.ChestConfig.ChestBigKeysMax, config.ChestConfig.ChestBigKeysMin)

	out := core.RandStringBytes(5) + ":" + s(acc.Uid) + ":" + chk + ":" + udid + ":" + s(acc.Uid) + ":" + s(smallLeft) + ":" + c.generateChestSmall(config) + ":" + s(acc.ChestSmallCount) + ":" +
		s(bigLeft) + ":" + c.generateChestBig(config) + ":" + s(acc.ChestBigCount) + ":" + s(chestType)
	out = strings.ReplaceAll(strings.ReplaceAll(base64.StdEncoding.EncodeToString([]byte(core.DoXOR(out, "59182"))), "/", "_"), "+", "-")
	c.output = core.RandStringBytes(5) + out + "|" + core.HashSolo4(out)
}

func (c *GDConnector) Profile_GetUserProfile(acc core.CAccount, selfUid int) {
	cf := core.CFriendship{DB: acc.DB}
	cm := core.CMessage{DB: acc.DB}
	c.output = c.getUserProfile(acc, cf.IsAlreadyFriend(acc.Uid, selfUid))
	if acc.Uid == selfUid {
		c.output += c.userProfilePersonal(cf.CountFriendRequests(acc.Uid, true), cm.CountMessages(acc.Uid, true))
	}
}

func (c *GDConnector) Profile_ListUserProfiles(accs []core.CAccount) {
	for _, acc := range accs {
		c.output += c.userListItem(acc)
	}
	c.output = c.output[:len(c.output)-1]
}

func (c *GDConnector) Profile_GetSearchableUsers(accs []core.CAccount, count int, page int) {
	for _, acc := range accs {
		c.output += c.getSearchableUser(acc)
	}
	c.output = fmt.Sprintf("%s#%d:%d:10", c.output[:len(c.output)-1], count, page*10)
}

// -- Internals --

func (c *GDConnector) getSearchableUser(acc core.CAccount) string {
	s := strconv.Itoa
	acc.LoadAuth(core.CAUTH_UID)
	acc.LoadVessels()
	acc.LoadStats()
	return "1:" + acc.Uname + ":2:" + s(acc.Uid) + ":3:" + s(acc.Stars) + ":4:" + s(acc.Demons) + ":8:" + s(acc.CPoints) + ":9:" + s(acc.GetShownIcon()) +
		":10:" + s(acc.ColorPrimary) + ":11:" + s(acc.ColorSecondary) + ":13:" + s(acc.Coins) + ":14:" + s(acc.IconType) + ":15:" + s(acc.Special) +
		":16:" + s(acc.Uid) + ":17:" + s(acc.UCoins) + ":52:" + s(acc.Moons) + "|"
}

// getAccountComment used to retrieve account comments (iterative, w/o hash)
func (c *GDConnector) getAccountComment(comment core.CComment) string {
	s := strconv.Itoa
	t, err := time.ParseInLocation("2006-01-02 15:04:05", comment.PostedTime, loc)
	if err != nil {
		t = time.Now()
	}
	age := core.GetDateAgo(t.Unix())
	return "2~" + comment.Comment + "~3~" + s(comment.Uid) + "~4~" + s(comment.Likes) + "~5~0~6~" + s(comment.Id) + "~7~" + s(core.ToInt(comment.IsSpam)) + "~9~" + age + "|"
}

// getLevelComment used to retrieve level comment (iterative, w/o hash)
func (c *GDConnector) getLevelComment(comment core.CComment) string {
	s := strconv.Itoa
	t, err := time.ParseInLocation("2006-01-02 15:04:05", comment.PostedTime, loc)
	if err != nil {
		t = time.Now()
	}
	age := core.GetDateAgo(t.Unix())
	acc := core.CAccount{DB: comment.DB, Uid: comment.Uid}
	if !acc.Exists(comment.Uid) {
		return ""
	}
	acc.LoadAuth(core.CAUTH_UID)
	acc.LoadStats()
	acc.LoadVessels()
	role := acc.GetRoleObj(false)
	if role.CommentColor != "" {
		role.CommentColor = "~12~" + role.CommentColor
	}
	return "2~" + comment.Comment + "~3~" + s(comment.Uid) + "~4~" + s(comment.Likes) + "~5~0~6~" + s(comment.Id) + "~7~" + s(core.ToInt(comment.IsSpam)) +
		"~8~" + s(comment.Uid) + "~9~" + age + "~10~" + s(comment.Percent) + "~11~" + s(role.ModLevel) + role.CommentColor + ":1~" + acc.Uname + "~9~" + s(acc.GetShownIcon()) +
		"~10~" + s(acc.ColorPrimary) + "~11~" + s(acc.ColorSecondary) + "~14~" + s(acc.IconType) + "~15~" + s(acc.Special) + s(acc.Uid) + "|"
}

// getCommentHistory used to retrieve level comment history of a user (iterative, w/o hash)
func (c *GDConnector) getCommentHistory(comment core.CComment, acc core.CAccount, role core.Role) string {
	s := strconv.Itoa
	t, err := time.ParseInLocation("2006-01-02 15:04:05", comment.PostedTime, loc)
	if err != nil {
		t = time.Now()
	}
	age := core.GetDateAgo(t.Unix())
	if role.CommentColor != "" {
		role.CommentColor = "~12~" + role.CommentColor
	}
	return "2~" + comment.Comment + "~3~" + s(comment.Uid) + "~4~" + s(comment.Likes) + "~5~0~6~" + s(comment.Id) + "~7~" + s(core.ToInt(comment.IsSpam)) +
		"~9~" + age + "~10~" + s(comment.Percent) + "~11~" + s(role.ModLevel) + "~12~" + role.CommentColor + ":1~" + acc.Uname + "~9~" + s(acc.GetShownIcon()) +
		"~10~" + s(acc.ColorPrimary) + "~11~" + s(acc.ColorSecondary) + "~14~" + s(acc.IconType) + "~15~" + s(acc.Special) + "~16~" + s(acc.Uid) + "|"
}

// getFriendRequest used to get friend request item (iterative, w/o hash)
func (c *GDConnector) getFriendRequest(frq map[string]string) string {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", frq["date"], loc)
	if err != nil {
		t = time.Now()
	}
	age := core.GetDateAgo(t.Unix())
	return "1:" + frq["uname"] + ":2:" + frq["uid"] + ":9:" + frq["iconId"] + ":10:" + frq["clr_primary"] + ":11:" + frq["clr_secondary"] +
		":14:" + frq["iconType"] + ":15:" + frq["special"] + ":16:" + frq["uid"] + ":32:" + frq["id"] + ":35:" + frq["comment"] + ":37:" + age + ":41:" + frq["isNew"] + "|"
}

// getMessage used to retrieve single message (w/o trailing hash)
func (c *GDConnector) getMessage(msg core.CMessage, uid int) string {
	s := strconv.Itoa
	t, err := time.ParseInLocation("2006-01-02 15:04:05", msg.PostedTime, loc)
	if err != nil {
		t = time.Now()
	}
	age := core.GetDateAgo(t.Unix())
	uidx := msg.UidDest
	if uid == msg.UidDest {
		uidx = msg.UidSrc
	}
	xacc := core.CAccount{DB: msg.DB, Uid: uidx}
	xacc.LoadAuth(core.CAUTH_UID)
	return "1:" + s(msg.Id) + ":2:" + s(uidx) + ":3:" + s(uidx) + ":4:" + msg.Subject + ":5:" + msg.Message + ":6:" + xacc.Uname + ":7:" + age +
		":8:" + s(core.ToInt(!msg.IsNew)) + ":9:" + s(core.ToInt(uid == msg.UidSrc))
}

// getMessageStr used to get message item (iterative, w/o hash)
func (c *GDConnector) getMessageStr(msg map[string]string, getSent bool) string {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", msg["date"], loc)
	if err != nil {
		t = time.Now()
	}
	age := core.GetDateAgo(t.Unix())
	return "1:" + msg["id"] + ":2:" + msg["uid"] + ":3:" + msg["uid"] + ":4:" + msg["subject"] + ":5:" + msg["message"] + ":6:" + msg["uname"] + ":7:" + age +
		":8:" + msg["isOld"] + ":9:" + strconv.Itoa(core.ToInt(getSent)) + "|"
}

func (c *GDConnector) getListSearch(cl core.CLevelList) (listStr string, user string, hash string) {
	s := strconv.Itoa

	t, err := time.ParseInLocation("2006-01-02 15:04:05", cl.UploadDate, loc)
	if err != nil {
		t = time.Now()
	}
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", cl.UpdateDate, loc)
	if err != nil {
		t2 = time.Now()
	}

	acc := core.CAccount{DB: cl.DB, Uid: cl.UID}
	if cl.SideloadUname == nil {
		if acc.Exists(acc.Uid) {
			acc.LoadAuth(core.CAUTH_UID)
		} else {
			acc.Uname = "[DELETED]"
		}
	} else {
		acc.Uname = *cl.SideloadUname
	}

	return "1:" + s(cl.ID) + ":2:" + cl.Name + ":3:" + cl.Description + ":5:" + s(cl.Version) + ":7:" + s(cl.Difficulty) +
			":10:" + s(cl.Downloads) + ":14:" + s(cl.Likes) + ":19:" + s(core.ToInt(cl.IsFeatured)) + ":28:" + s(int(t.Unix())) +
			":29:" + s(int(t2.Unix())) + ":49:" + s(cl.UID) + ":50:" + acc.Uname + ":51:" + cl.Levels + ":55:" + s(cl.Diamonds) +
			":56:" + s(cl.LevelDiamonds) + "|",
		s(acc.Uid) + ":" + acc.Uname + ":" + s(acc.Uid) + "|", ""

}

// getLevelSearch used to retrieve data about level in search (iterative, w/ half-hash), returns (lvlString, lvlHash, usrString)
func (c *GDConnector) getLevelSearch(cl core.CLevel, gau bool) (string, string, string) {
	s := strconv.Itoa
	diffNom := 0
	if cl.Difficulty > 0 {
		diffNom = 10
	}
	var auto int
	if cl.Difficulty < 0 {
		auto = 1
		cl.Difficulty = 0
	}
	coinsVer := 0
	if cl.Coins > 0 {
		coinsVer = 1
	}
	demonDiff := 3
	isDemon := 0
	if cl.DemonDifficulty >= 0 {
		isDemon = 1
		demonDiff = cl.DemonDifficulty
	}
	acc := core.CAccount{DB: cl.DB, Uid: cl.Uid}
	if cl.SideloadUname == nil {
		if acc.Exists(acc.Uid) {
			acc.LoadAuth(core.CAUTH_UID)
		} else {
			acc.Uname = "[DELETED]"
		}
	} else {
		acc.Uname = *cl.SideloadUname
	}

	gaustr := ""
	if gau {
		gaustr = ":44:1"
	}
	//lvlString
	strID := s(cl.Id)
	sliceL := len(strID) - 1
	//if sliceL==0 {sliceL=1}
	return "1:" + s(cl.Id) + ":2:" + cl.Name + ":3:" + cl.Description + ":5:" + s(cl.Version) + ":6:" + s(cl.Uid) + ":8:" + s(diffNom) +
			":9:" + s(cl.Difficulty) + ":10:" + s(cl.Downloads) + ":12:" + s(cl.TrackId) + ":13:" + s(cl.VersionGame) + ":14:" + s(cl.Likes) +
			":15:" + s(cl.Length) + ":17:" + s(isDemon) + ":18:" + s(cl.StarsGot) + ":19:" + s(cl.IsFeatured) + ":25:" + s(auto) +
			":30:" + s(cl.OrigId) + ":31:" + s(core.ToInt(cl.Is2p)) + ":35:" + s(cl.SongId) + ":37:" + s(cl.Ucoins) + ":38:" + s(coinsVer) +
			":39:" + s(cl.StarsRequested) + ":42:" + s(cl.IsEpic) + ":43:" + s(demonDiff) + gaustr + ":45:" + s(cl.Objects) + ":46:1:47:2|",
		//lvlHash
		string(strID[0]) + string(strID[sliceL]) + s(cl.StarsGot) + s(coinsVer),
		//usrString
		s(acc.Uid) + ":" + acc.Uname + ":" + s(acc.Uid) + "|"

	//44 isGauntlet
}

// generateChestSmall used to generate small chest loot
func (c *GDConnector) generateChestSmall(config core.ConfigBlob) string {
	s := strconv.Itoa
	rand.Seed(time.Now().UnixNano())
	intR := func(min, max int) int { return rand.Intn(core.MaxInt(max-min+1, 0)) + min }
	return s(intR(config.ChestConfig.ChestSmallOrbsMin, config.ChestConfig.ChestSmallOrbsMax)) + "," +
		s(intR(config.ChestConfig.ChestSmallDiamondsMin, config.ChestConfig.ChestSmallDiamondsMax)) + "," +
		s(config.ChestConfig.ChestSmallShards[rand.Intn(len(config.ChestConfig.ChestSmallShards))]) + "," +
		s(intR(config.ChestConfig.ChestSmallKeysMin, config.ChestConfig.ChestSmallKeysMax))
}

// generateChestBig used to generate big chest loot
func (c *GDConnector) generateChestBig(config core.ConfigBlob) string {
	s := strconv.Itoa
	rand.Seed(time.Now().UnixNano())
	intR := func(min, max int) int { return rand.Intn(max-min+1) + min }
	return s(intR(config.ChestConfig.ChestBigOrbsMin, config.ChestConfig.ChestBigOrbsMax)) + "," +
		s(intR(config.ChestConfig.ChestBigDiamondsMin, config.ChestConfig.ChestBigDiamondsMax)) + "," +
		s(config.ChestConfig.ChestBigShards[rand.Intn(len(config.ChestConfig.ChestBigShards))]) + "," +
		s(intR(config.ChestConfig.ChestBigKeysMin, config.ChestConfig.ChestBigKeysMax))
}

// getUserProfile used at getUserInfo (w/o trailing hash)
func (c *GDConnector) getUserProfile(acc core.CAccount, isFriend bool) string {
	s := strconv.Itoa
	// 51=13 - color3, 53=34 -> swingcopter, 54=5 -> jetpack?,
	role := acc.GetRoleObj(false)
	rank := acc.GetLeaderboardRank()
	return "1:" + acc.Uname + ":2:" + s(acc.Uid) + ":3:" + s(acc.Stars) + ":4:" + s(acc.Demons) + ":6:" + s(rank) + ":7:" + s(acc.Uid) +
		":8:" + s(acc.CPoints) + ":9:" + s(acc.GetShownIcon()) + ":10:" + s(acc.ColorPrimary) + ":11:" + s(acc.ColorSecondary) + ":13:" + s(acc.Coins) +
		":14:" + s(acc.IconType) + ":15:" + s(acc.Special) + ":16:" + s(acc.Uid) + ":17:" + s(acc.UCoins) + ":18:" + s(acc.MS) + ":19:" + s(acc.FrS) +
		":20:" + acc.Youtube + ":21:" + s(acc.Cube) + ":22:" + s(acc.Ship) + ":23:" + s(acc.Ball) + ":24:" + s(acc.Ufo) + ":25:" + s(acc.Wave) + ":26:" + s(acc.Robot) +
		":28:" + s(acc.Trace) + ":29:1:30:" + s(rank) + ":31:" + s(core.ToInt(isFriend)) + ":43:" + s(acc.Spider) + ":44:" + acc.Twitter +
		":45:" + acc.Twitch + ":46:" + s(acc.Diamonds) + ":48:" + s(acc.Death) + ":49:" + s(role.ModLevel) + ":50:" + s(acc.CS) + ":51:" + s(acc.ColorGlow) +
		":52:" + s(acc.Moons) + ":53:" + s(acc.Swing) + ":54:" + s(acc.Jetpack)
}

// userProfilePersonal used at getUserInfo to append some data if user is requesting themselves (w/o trailing hash)
func (c *GDConnector) userProfilePersonal(frReq int, msgNewCnt int) string {
	return ":38:" + strconv.Itoa(msgNewCnt) + ":39:" + strconv.Itoa(frReq) + ":40:0"
}

// userListItem used at getUserList to provide minimum data for user lists (iterative, w/o hash)
func (c *GDConnector) userListItem(acc core.CAccount) string {
	s := strconv.Itoa
	return "1:" + acc.Uname + ":2:" + s(acc.Uid) + ":9:" + s(acc.GetShownIcon()) + ":10:" + s(acc.ColorPrimary) + ":11:" + s(acc.ColorSecondary) +
		":14:" + s(acc.IconType) + ":15:" + s(acc.Special) + ":16:" + s(acc.Uid) + ":18:0:41:1|"
}

//
//
//
//
//

// GetAccLeaderboardItem used to retrieve user for leaderboards (iterative, w/o trailing hash)
func GetAccLeaderboardItem(acc core.CAccount, lk int) string {
	s := strconv.Itoa
	acc.LoadAll()
	return "1:" + acc.Uname + ":2:" + s(acc.Uid) + ":3:" + s(acc.Stars) + ":4:" + s(acc.Demons) + ":6:" + s(lk) + ":7:" + s(acc.Uid) +
		":8:" + s(acc.CPoints) + ":9:" + s(acc.GetShownIcon()) + ":10:" + s(acc.ColorPrimary) + ":11:" + s(acc.ColorSecondary) + ":13:" + s(acc.Coins) +
		":14:" + s(acc.IconType) + ":15:" + s(acc.Special) + ":16:" + s(acc.Uid) + ":17:" + s(acc.UCoins) + ":46:" + s(acc.Diamonds) + ":52:" + s(acc.Moons) + "|"
}

// GetLeaderboardScore used to retrieve leaderboard scores (iterative, w/o trailing hash)
func GetLeaderboardScore(score core.CScores) string {
	s := strconv.Itoa
	acc := core.CAccount{DB: score.DB, Uid: score.Uid}
	acc.LoadAuth(core.CAUTH_UID)
	acc.LoadVessels()
	acc.LoadStats()
	t, err := time.ParseInLocation("2006-01-02 15:04:05", score.PostedTime, loc)
	if err != nil {
		t = time.Now()
	}
	age := core.GetDateAgo(t.Unix())
	return "1:" + acc.Uname + ":2:" + s(acc.Uid) + ":3:" + s(score.Percent) + ":6:" + s(score.Ranking) + ":9:" + s(acc.GetShownIcon()) +
		":10:" + s(acc.ColorPrimary) + ":11:" + s(acc.ColorSecondary) + ":13:" + s(score.Coins) + ":14:" + s(acc.IconType) + ":15:" + s(acc.Special) +
		":16:" + s(acc.Uid) + ":42:" + age + "|"
}
