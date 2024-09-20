package api

import (
	"encoding/base64"
)

func BannedTemplateComment(text string) string {
	text=base64.StdEncoding.EncodeToString([]byte(text))
	return "2~"+text+"~3~1~4~0~7~0~10~0~9~1 second~6~1:1~M41dss~9~98~10~35~11~3~14~0~15~2~16~0"
}

func BannedTemplateUserProfile() string {
	uname:="M41dss"
	return "1:"+uname+":2:0:3:41:4:0:6:1:7:0:8:41:9:98:10:35:11::13:41:14:0:15:2:16:0:17:41:18:0:19:0:20::21:98:22:SHIP_ID"+
		":23:BALL_ID:24:UFO_ID:25:WAVE_ID:26:ROBOT_ID:28:1:29:1:30:1:31:0:43:SPIDER_ID:44::45::46:41:48:1:49:2:50:0:38:0:39:0:40:0"
}

func BannedTemplateUserListItem() string {
	uname:="M41dss"
	return "1:"+uname+":2:0:9:98:10:35:11:3:14:0:15:2:16:0:18:0:41:1|"
}