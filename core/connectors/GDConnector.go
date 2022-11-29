// Package connectors allow translating beautiful typed data to a hell of a mess RobTop format
// and also to communicate with outside world
package connectors

import (
	"HalogenGhostCore/core"
	"encoding/base64"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// GetUserProfile used at getUserInfo (w/o trailing hash)
func GetUserProfile(acc core.CAccount, isFriend bool) string {
	s:=strconv.Itoa
	role:=acc.GetRoleObj(false)
	return "1:"+acc.Uname+":2:"+s(acc.Uid)+":3:"+s(acc.Stars)+":4:"+s(acc.Demons)+":6:"+ s(acc.GetLeaderboardRank())+":7:"+s(acc.Uid)+
		":8:"+s(acc.CPoints)+":9:"+s(acc.GetShownIcon())+":10:"+s(acc.ColorPrimary)+":11:"+s(acc.ColorSecondary)+":13:"+s(acc.Coins)+
		":14:"+s(acc.IconType)+":15:"+s(acc.Special)+":16:"+s(acc.Uid)+":17:"+s(acc.UCoins)+":18:"+s(acc.MS)+":19:"+s(acc.FrS)+
		":20:"+acc.Youtube+":21:"+s(acc.Cube)+":22:"+s(acc.Ship)+":23:"+s(acc.Ball)+":24:"+s(acc.Ufo)+":25:"+s(acc.Wave)+":26:"+s(acc.Robot)+
		":28:"+s(acc.Trace)+":29:1:30:"+s(acc.GetLeaderboardRank())+":31:"+s(core.ToInt(isFriend))+":43:"+s(acc.Spider)+":44:"+acc.Twitter+
		":45:"+acc.Twitch+":46:"+s(acc.Diamonds)+":48:"+s(acc.Death)+":49:"+s(role.ModLevel)+":50:"+s(acc.CS)
}

// UserProfilePersonal used at getUserInfo to append some data if user is requesting themselves (w/o trailing hash)
func UserProfilePersonal(frReq int,msgNewCnt int) string {
	return ":38:"+strconv.Itoa(msgNewCnt)+":39:"+strconv.Itoa(frReq)+":40:0"
}

// UserListItem used at getUserList to provide minimum data for user lists (iterative, w/o hash)
func UserListItem(acc core.CAccount) string {
	s:=strconv.Itoa
	return "1:"+acc.Uname+":2:"+s(acc.Uid)+":9:"+s(acc.GetShownIcon())+":10:"+s(acc.ColorPrimary)+":11:"+s(acc.ColorSecondary)+
		":14:"+s(acc.IconType)+":15:"+s(acc.Special)+":16:"+s(acc.Uid)+":18:0:41:1|"
}

// UserSearchItem used at getUsers (w/ trailing hash)
func UserSearchItem(acc core.CAccount) string {
	s:=strconv.Itoa
	return "1:"+acc.Uname+":2:"+s(acc.Uid)+":3:"+s(acc.Stars)+":4:"+s(acc.Demons)+":8:"+s(acc.CPoints)+":9:"+s(acc.GetShownIcon())+
		":10:"+s(acc.ColorPrimary)+":11:"+s(acc.ColorSecondary)+":13:"+s(acc.Coins)+":14:"+s(acc.IconType)+":15:"+s(acc.Special)+
		":16:"+s(acc.Uid)+":17:"+s(acc.UCoins)+"#1:0:10"
}

// GetAccountComment used to retrieve account comments (iterative, w/o hash)
func GetAccountComment(comment core.CComment) string {
	s:=strconv.Itoa
	t,err:=time.Parse("2006-01-02 15:04:05",comment.PostedTime)
	if err!=nil {t=time.Now()}
	age:=core.GetDateAgo(t.Unix())
	return "2~"+comment.Comment+"~3~"+s(comment.Uid)+"~4~"+s(comment.Likes)+"~5~0~6~"+s(comment.Id)+"~7~"+s(core.ToInt(comment.IsSpam))+"~9~"+age+"|"
}

// GetLevelComment used to retrieve level comment (iterative, w/o hash)
func GetLevelComment(comment core.CComment) string {
	s:=strconv.Itoa
	t,err:=time.Parse("2006-01-02 15:04:05",comment.PostedTime)
	if err!=nil {t=time.Now()}
	age:=core.GetDateAgo(t.Unix())
	acc:=core.CAccount{DB: comment.DB, Uid: comment.Uid}
	if !acc.Exists(comment.Uid) {return ""}
	acc.LoadAuth(core.CAUTH_UID)
	acc.LoadStats()
	acc.LoadVessels()
	role:=acc.GetRoleObj(false)
	if role.CommentColor!="" {role.CommentColor="~12~"+role.CommentColor}
	return "2~"+comment.Comment+"~3~"+s(comment.Uid)+"~4~"+s(comment.Likes)+"~5~0~6~"+s(comment.Id)+"~7~"+s(core.ToInt(comment.IsSpam))+
		"~9~"+age+"~10~"+s(comment.Percent)+"~11~"+s(role.ModLevel)+role.CommentColor+":1~"+acc.Uname+"~9~"+s(acc.GetShownIcon())+
		"~10~"+s(acc.ColorPrimary)+"~11~"+s(acc.ColorSecondary)+"~14~"+s(acc.IconType)+"~15~"+s(acc.Special)+s(acc.Uid)+"|"
}

// GetCommentHistory used to retrieve level comment history of a user (iterative, w/o hash)
func GetCommentHistory(comment core.CComment, acc core.CAccount, role core.Role) string {
	s:=strconv.Itoa
	t,err:=time.Parse("2006-01-02 15:04:05",comment.PostedTime)
	if err!=nil {t=time.Now()}
	age:=core.GetDateAgo(t.Unix())
	if role.CommentColor!="" {role.CommentColor="~12~"+role.CommentColor}
	return "2~"+comment.Comment+"~3~"+s(comment.Uid)+"~4~"+s(comment.Likes)+"~5~0~6~"+s(comment.Id)+"~7~"+s(core.ToInt(comment.IsSpam))+
		"~9~"+age+"~10~"+s(comment.Percent)+"~11~"+s(role.ModLevel)+role.CommentColor+":1~"+acc.Uname+"~9~"+s(acc.GetShownIcon())+
		"~10~"+s(acc.ColorPrimary)+"~11~"+s(acc.ColorSecondary)+"~14~"+s(acc.IconType)+"~15~"+s(acc.Special)+s(acc.Uid)+"|"
}

// GetFriendRequest used to get friend request item (iterative, w/o hash)
func GetFriendRequest(frq map[string]string) string {
	t,err:=time.Parse("2006-01-02 15:04:05",frq["date"])
	if err!=nil {t=time.Now()}
	age:=core.GetDateAgo(t.Unix())
	return "1:"+frq["uname"]+":2:"+frq["uid"]+":9:"+frq["iconId"]+":10:"+frq["clr_primary"]+":11:"+frq["clr_secondary"]+
		":14:"+frq["iconType"]+":15:"+frq["special"]+":16:"+frq["uid"]+":32:"+frq["id"]+":35:"+frq["comment"]+":37:"+age+":41:"+frq["isNew"]+"|"
}

// GetMessage used to retrieve single message (w/o trailing hash)
func GetMessage(msg core.CMessage, uid int) string {
	s:=strconv.Itoa
	t,err:=time.Parse("2006-01-02 15:04:05",msg.PostedTime)
	if err!=nil {t=time.Now()}
	age:=core.GetDateAgo(t.Unix())
	uidx:=msg.UidDest
	if uid==msg.UidDest {uidx=msg.UidSrc}
	xacc:=core.CAccount{DB: msg.DB, Uid: uidx}
	xacc.LoadAuth(core.CAUTH_UID)
	return "1:"+s(msg.Id)+":2:"+s(uidx)+":3:"+s(uidx)+":4:"+msg.Subject+":5:"+msg.Message+":6:"+xacc.Uname+":7:"+age+
		":8:"+s(core.ToInt(!msg.IsNew))+":9:"+s(core.ToInt(uid==msg.UidSrc))
}

// GetMessageStr used to get message item (iterative, w/o hash)
func GetMessageStr(msg map[string]string, getSent bool) string {
	t,err:=time.Parse("2006-01-02 15:04:05",msg["date"])
	if err!=nil {t=time.Now()}
	age:=core.GetDateAgo(t.Unix())
	return "1:"+msg["id"]+":2:"+msg["uid"]+":3:"+msg["uid"]+":4:"+msg["subject"]+":5:"+msg["message"]+":6:"+msg["uname"]+":7:"+age+
		":8:"+msg["isOld"]+":9:"+strconv.Itoa(core.ToInt(getSent))+"|"
}

// GetMusic used to get simple music string (w/o traling hash)
func GetMusic(mus core.CMusic) string {
	size:=strconv.FormatFloat(mus.Size, 'f', 2, 64)
	return "1~|~"+strconv.Itoa(mus.Id)+"~|~2~|~"+mus.Name+"~|~3~|~1~|~4~|~"+mus.Artist+"~|~5~|~"+size+"~|~6~|~~|~10~|~"+
		url.QueryEscape(mus.Url)
}

//used to get simple top artists string (w/o trailing hash)
func GetTopArtists(artists map[string]string) string {
	out:=""
	for artist, youtube := range artists {
		out+="4:"+artist+":7:"+youtube+"|"
	}
	return out[:len(out)-1]
}

// GenerateChestSmall used to generate small chest loot
func GenerateChestSmall(config core.ConfigBlob) string {
	s:=strconv.Itoa
	rand.Seed(time.Now().UnixNano())
	intR:= func(min, max int) int {return rand.Intn(max-min+1)+min}
	return s(intR(config.ChestConfig.ChestSmallOrbsMin,config.ChestConfig.ChestSmallOrbsMax))+","+
		s(intR(config.ChestConfig.ChestSmallDiamondsMin,config.ChestConfig.ChestSmallDiamondsMax))+ ","+
		s(config.ChestConfig.ChestSmallShards[rand.Intn(len(config.ChestConfig.ChestSmallShards))])+","+
		s(intR(config.ChestConfig.ChestSmallKeysMin,config.ChestConfig.ChestSmallKeysMax))
}

// GenerateChestBig used to generate big chest loot
func GenerateChestBig(config core.ConfigBlob) string {
	s:=strconv.Itoa
	rand.Seed(time.Now().UnixNano())
	intR:= func(min, max int) int {return rand.Intn(max-min+1)+min}
	return s(intR(config.ChestConfig.ChestBigOrbsMin,config.ChestConfig.ChestBigOrbsMax))+","+
		s(intR(config.ChestConfig.ChestBigDiamondsMin,config.ChestConfig.ChestBigDiamondsMax))+ ","+
		s(config.ChestConfig.ChestBigShards[rand.Intn(len(config.ChestConfig.ChestBigShards))])+","+
		s(intR(config.ChestConfig.ChestBigKeysMin,config.ChestConfig.ChestBigKeysMax))
}

// ChestOutput used to retrieve all chest data (w/ trailing hash)
func ChestOutput(acc core.CAccount, config core.ConfigBlob, udid string, chk string, smallLeft int, bigLeft int, chestType int) string {
	s:=strconv.Itoa
	out:=core.RandStringBytes(5)+":"+s(acc.Uid)+":"+chk+":"+udid+":"+s(acc.Uid)+":"+s(smallLeft)+":"+GenerateChestSmall(config)+":"+s(acc.ChestSmallCount)+":"+
		s(bigLeft)+":"+GenerateChestBig(config)+":"+s(acc.ChestBigCount)+":"+s(chestType)
	out=strings.ReplaceAll(strings.ReplaceAll(base64.StdEncoding.EncodeToString([]byte(core.DoXOR(out,"59182"))),"/","_"),"+","-")
	return core.RandStringBytes(5)+out+"|"+core.HashSolo4(out)
}

// ChallengesOutput used to retrieve all quests/challenges data (w/ trailing hash)
func ChallengesOutput(cq core.CQuests, uid int, chk string, udid string) string{
	s:=strconv.Itoa
	virt:=core.RandStringBytes(5)
	tme,_:=time.Parse("2006-01-02 15:04:05",strings.Split(time.Now().Format("2006-01-02 15:04:05")," ")[0]+" 00:00:00")
	//!Additional 10800 Review is needed
	timeLeft:=int(tme.AddDate(0,0,1).Unix()-(time.Now().Unix()+10800))
	out:=virt+":"+s(uid)+":"+chk+":"+udid+":"+s(uid)+":"+s(timeLeft)+":"+cq.GetQuests(uid)
	out=strings.ReplaceAll(strings.ReplaceAll(base64.StdEncoding.EncodeToString([]byte(core.DoXOR(out,"19847"))),"/","_"),"+","-")
	return virt+out+"|"+core.HashSolo3(out)
}

// GetAccLeaderboardItem used to retrieve user for leaderboards (iterative, w/o trailing hash)
func GetAccLeaderboardItem(acc core.CAccount,lk int) string {
	s:=strconv.Itoa
	acc.LoadAll()
	return "1:"+acc.Uname+":2:"+s(acc.Uid)+":3:"+s(acc.Stars)+":4:"+s(acc.Demons)+":6:"+s(lk)+":7:"+s(acc.Uid)+
		":8:"+s(acc.CPoints)+":9:"+s(acc.GetShownIcon())+":10:"+s(acc.ColorPrimary)+":11:"+s(acc.ColorSecondary)+":13:"+s(acc.Coins)+
		":14:"+s(acc.IconType)+":15:"+s(acc.Special)+":16:"+s(acc.Uid)+":17:"+s(acc.UCoins)+":46:"+s(acc.Diamonds)+"|"
}

// GetLeaderboardScore used to retrieve leaderboard scores (iterative, w/o trailing hash)
func GetLeaderboardScore(score core.CScores) string {
	s:=strconv.Itoa
	acc:=core.CAccount{DB: score.DB, Uid: score.Uid}
	acc.LoadAuth(core.CAUTH_UID)
	acc.LoadVessels()
	acc.LoadStats()
	t,err:=time.Parse("2006-01-02 15:04:05",score.PostedTime)
	if err!=nil {t=time.Now()}
	age:=core.GetDateAgo(t.Unix())
	return "1:"+acc.Uname+":2:"+s(acc.Uid)+":3:"+s(score.Percent)+":6:"+s(score.Ranking)+":9:"+s(acc.GetShownIcon())+
		":10:"+s(acc.ColorPrimary)+":11:"+s(acc.ColorSecondary)+":13:"+s(score.Coins)+":14:"+s(acc.IconType)+":15:"+s(acc.Special)+
		":16:"+s(acc.Uid)+":42:"+age+"|"
}

// GetLevelFull used to retrieve full Level data (w/ trailing hash)
func GetLevelFull(cl core.CLevel, password string, phash string, quest_id int) string {
	s:=strconv.Itoa
	t,err:=time.Parse("2006-01-02 15:04:05",cl.UploadDate)
	if err!=nil {t=time.Now()}
	uplAge:=core.GetDateAgo(t.Unix())
	t2,err:=time.Parse("2006-01-02 15:04:05",cl.UpdateDate)
	if err!=nil {t2=time.Now()}
	updAge:=core.GetDateAgo(t2.Unix())
	diffNom:=0
	if cl.Difficulty>0 {diffNom=10}
	var auto int
	if cl.Difficulty<0 {
		auto=1
		cl.Difficulty=0
	}
	coinsVer:=0
	if cl.Coins>0 {coinsVer=1}
	demonDiff:=3
	isDemon:=0
	if cl.DemonDifficulty>=0 {
		isDemon=1
		demonDiff=cl.DemonDifficulty
	}
	quest:=""
	questHash:=""
	if quest_id>0{
		quest=":41:"+s(quest_id)
		acc:=core.CAccount{DB: cl.DB, Uid: cl.Uid}
		acc.LoadAuth(core.CAUTH_UID)
		questHash="#"+s(acc.Uid)+":"+acc.Uname+":"+s(acc.Uid)
	}
	hash:=s(cl.Uid)+","+s(cl.StarsGot)+","+s(isDemon)+","+s(cl.Id)+","+s(coinsVer)+","+s(core.ToInt(cl.IsFeatured))+","+phash+
		","+s(quest_id)
	return "1:"+s(cl.Id)+":2:"+cl.Name+":3:"+cl.Description+":4:"+cl.StringLevel+":5:"+s(cl.Version)+":6:"+s(cl.Uid)+":8:"+s(diffNom)+
		":9:"+s(cl.Difficulty)+":10:"+s(cl.Downloads)+":12:"+s(cl.TrackId)+":13:"+s(cl.VersionGame)+":14:"+s(cl.Likes)+
		":15:"+s(cl.Length)+":17:"+s(isDemon)+":18:"+s(cl.StarsGot)+":19:"+s(core.ToInt(cl.IsFeatured))+":25:"+s(auto)+
		":27:"+password+":28:"+uplAge+":29:"+updAge+":30:"+s(cl.OrigId)+":31:"+s(core.ToInt(cl.Is2p))+":35:"+s(cl.SongId)+
		":36:"+cl.StringExtra+":37:"+s(cl.Ucoins)+":38:"+s(coinsVer)+":39:"+s(cl.StarsRequested)+":40:"+s(core.ToInt(cl.IsLDM))+
		":42:"+s(cl.IsEpic)+":43:"+s(demonDiff)+":45:"+s(cl.Objects)+":46:1:47:2:48:"+cl.StringSettings+quest+
		"#"+core.HashSolo(cl.StringLevel)+"#"+core.HashSolo2(hash)+questHash

	//44 isGauntlet
}

// GetLevelSearch used to retrieve data about level in search (iterative, w/ half-hash), returns (lvlString, lvlHash, usrString)
func GetLevelSearch(cl core.CLevel, gau bool) (string, string, string) {
	s:=strconv.Itoa
	diffNom:=0
	if cl.Difficulty>0 {diffNom=10}
	var auto int
	if cl.Difficulty<0 {
		auto=1
		cl.Difficulty=0
	}
	coinsVer:=0
	if cl.Coins>0 {coinsVer=1}
	demonDiff:=3
	isDemon:=0
	if cl.DemonDifficulty>=0 {
		isDemon=1
		demonDiff=cl.DemonDifficulty
	}
	acc:=core.CAccount{DB: cl.DB, Uid: cl.Uid}
	if acc.Exists(acc.Uid) {
		acc.LoadAuth(core.CAUTH_UID)
	}else{
		acc.Uname="[DELETED]"
	}

	gaustr:=""
	if gau {gaustr=":44:1"}
	//lvlString
	strID:=s(cl.Id)
	sliceL:=len(strID)-1
	//if sliceL==0 {sliceL=1}
	return "1:"+s(cl.Id)+":2:"+cl.Name+":3:"+cl.Description+":5:"+s(cl.Version)+":6:"+s(cl.Uid)+":8:"+s(diffNom)+
		":9:"+s(cl.Difficulty)+":10:"+s(cl.Downloads)+":12:"+s(cl.TrackId)+":13:"+s(cl.VersionGame)+":14:"+s(cl.Likes)+
		":15:"+s(cl.Length)+":17:"+s(isDemon)+":18:"+s(cl.StarsGot)+":19:"+s(core.ToInt(cl.IsFeatured))+":25:"+s(auto)+
		":30:"+s(cl.OrigId)+":31:"+s(core.ToInt(cl.Is2p))+":35:"+s(cl.SongId)+":37:"+s(cl.Ucoins)+ ":38:"+s(coinsVer)+
		":39:"+s(cl.StarsRequested)+ ":42:"+s(cl.IsEpic)+":43:"+s(demonDiff)+gaustr+":45:"+s(cl.Objects)+":46:1:47:2|",
		//lvlHash
		string(strID[0])+string(strID[sliceL])+s(cl.StarsGot)+s(coinsVer),
		//usrString
		s(acc.Uid)+":"+acc.Uname+":"+s(acc.Uid)+"|"

	//44 isGauntlet
}