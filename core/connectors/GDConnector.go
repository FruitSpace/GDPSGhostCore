package connectors

import (
	"HalogenGhostCore/core"
	"encoding/base64"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

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

func UserProfilePersonal(frReq int,msgNewCnt int) string {
	return ":38:"+strconv.Itoa(msgNewCnt)+":39:"+strconv.Itoa(frReq)+":40:0"
}

func UserListItem(acc core.CAccount) string {
	s:=strconv.Itoa
	return "1:"+acc.Uname+":2:"+s(acc.Uid)+":9:"+s(acc.GetShownIcon())+":10:"+s(acc.ColorPrimary)+":11:"+s(acc.ColorSecondary)+
		":14:"+s(acc.IconType)+":15:"+s(acc.Special)+":16:"+s(acc.Uid)+":18:0:41:1|"
}

func UserSearchItem(acc core.CAccount) string {
	s:=strconv.Itoa
	return "1:"+acc.Uname+":2:"+s(acc.Uid)+":3:"+s(acc.Stars)+":4:"+s(acc.Demons)+":8:"+s(acc.CPoints)+":9:"+s(acc.GetShownIcon())+
		":10:"+s(acc.ColorPrimary)+":11:"+s(acc.ColorSecondary)+":13:"+s(acc.Coins)+":14:"+s(acc.IconType)+":15:"+s(acc.Special)+
		":16:"+s(acc.Uid)+":17:"+s(acc.UCoins)+"#1:0:10"
}

func GetAccountComment(comment core.CComment) string {
	s:=strconv.Itoa
	t,err:=time.Parse("2006-01-02 15:04:05",comment.PostedTime)
	if err!=nil {t=time.Now()}
	age:=core.GetDateAgo(t.Unix())
	return "2~"+comment.Comment+"~3~"+s(comment.Uid)+"~4~"+s(comment.Likes)+"~5~0~6~"+s(comment.Id)+"~7~"+s(core.ToInt(comment.IsSpam))+"~9~"+age+"|"
}

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

func GetFriendRequest(frq map[string]string) string {
	t,err:=time.Parse("2006-01-02 15:04:05",frq["date"])
	if err!=nil {t=time.Now()}
	age:=core.GetDateAgo(t.Unix())
	return "1:"+frq["uname"]+":2:"+frq["uid"]+":9:"+frq["iconId"]+":10:"+frq["clr_primary"]+":11:"+frq["clr_secondary"]+
		":14:"+frq["iconType"]+":15:"+frq["special"]+":16:"+frq["uid"]+":32:"+frq["id"]+":35:"+frq["comment"]+":37:"+age+":41:"+frq["isNew"]+"|"
}

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

func GetMessageStr(msg map[string]string, getSent bool) string {
	t,err:=time.Parse("2006-01-02 15:04:05",msg["date"])
	if err!=nil {t=time.Now()}
	age:=core.GetDateAgo(t.Unix())
	return "1:"+msg["id"]+":2:"+msg["uid"]+":3:"+msg["uid"]+":4:"+msg["subject"]+":5:"+msg["message"]+":6:"+msg["uname"]+":7:"+age+
		":8:"+msg["isOld"]+":9:"+strconv.Itoa(core.ToInt(getSent))+"|"
}

func GetMusic(mus core.CMusic) string {
	return "1~|~"+strconv.Itoa(mus.Id)+"~|~2~|~"+mus.Name+"~|~3~|~1~|~4~|~"+mus.Artist+"~|~5~|~"+mus.Size+"~|~6~|~~|~10~|~"+mus.Url
}

func GetTopArtists(artists map[string]string) string {
	out:=""
	for artist, youtube := range artists {
		out+="4:"+artist+":7:"+youtube+"|"
	}
	return out[:len(out)-1]
}


func GenerateChestSmall(config core.ConfigBlob) string {
	s:=strconv.Itoa
	rand.Seed(time.Now().UnixNano())
	intR:= func(min, max int) int {return rand.Intn(max-min+1)+min}
	return s(intR(config.ChestConfig.ChestSmallOrbsMin,config.ChestConfig.ChestSmallOrbsMax))+","+
		s(intR(config.ChestConfig.ChestSmallDiamondsMin,config.ChestConfig.ChestSmallDiamondsMax))+ ","+
		s(config.ChestConfig.ChestSmallShards[rand.Intn(len(config.ChestConfig.ChestSmallShards))])+","+
		s(intR(config.ChestConfig.ChestSmallKeysMin,config.ChestConfig.ChestSmallKeysMax))
}

func GenerateChestBig(config core.ConfigBlob) string {
	s:=strconv.Itoa
	rand.Seed(time.Now().UnixNano())
	intR:= func(min, max int) int {return rand.Intn(max-min+1)+min}
	return s(intR(config.ChestConfig.ChestBigOrbsMin,config.ChestConfig.ChestBigOrbsMax))+","+
		s(intR(config.ChestConfig.ChestBigDiamondsMin,config.ChestConfig.ChestBigDiamondsMax))+ ","+
		s(config.ChestConfig.ChestBigShards[rand.Intn(len(config.ChestConfig.ChestBigShards))])+","+
		s(intR(config.ChestConfig.ChestBigKeysMin,config.ChestConfig.ChestBigKeysMax))
}

func ChestOutput(acc core.CAccount, config core.ConfigBlob, udid string, chk string, smallLeft int, bigLeft int, chestType int) string {
	s:=strconv.Itoa
	out:=core.RandStringBytes(5)+":"+s(acc.Uid)+":"+chk+":"+udid+":"+s(acc.Uid)+":"+s(smallLeft)+GenerateChestSmall(config)+":"+s(acc.ChestSmallCount)+":"+
		s(bigLeft)+GenerateChestBig(config)+":"+s(acc.ChestBigCount)+":"+s(chestType)
	out=strings.ReplaceAll(strings.ReplaceAll(base64.StdEncoding.EncodeToString([]byte(core.DoXOR(out,"59182"))),"/","_"),"+","-")
	return core.RandStringBytes(5)+out+"|"+core.HashSolo4(out)
}