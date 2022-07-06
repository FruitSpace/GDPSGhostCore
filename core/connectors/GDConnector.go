package connectors

import (
	"HalogenGhostCore/core"
	"strconv"
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