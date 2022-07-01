package api

import "encoding/base64"

func BannedTemplateComment(text string) string {
	text=base64.StdEncoding.EncodeToString([]byte(text))
	return "2~"+text+"~3~1~4~0~7~0~10~0~9~1 second~6~1:1~M41dss~9~98~10~35~11~3~14~0~15~0~16~0"
}
