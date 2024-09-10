package connectors

import (
	"HalogenGhostCore/core"
	"time"
)

var loc, _ = time.LoadLocation("Europe/Moscow")

type Connector interface {
	Output() string

	Error(code string, reason string)
	Success(message string)
	Account_Sync(savedata string)
	Account_Login(uid int)
	Comment_AccountGet(comments []core.CComment, count int, page int)
	Comment_LevelGet(comments []core.CComment, count int, page int)
	Comment_HistoryGet(comments []core.CComment, acc core.CAccount, role core.Role, count int, page int)
	Communication_FriendGetRequests(reqs []map[string]string, count int, page int)
	Communication_MessageGet(message core.CMessage, uid int)
}

func NewConnector(isJson bool) Connector {
	if isJson {
		return &JSONConnector{output: make(map[string]interface{})}
	} else {
		return &GDConnector{}
	}
}
