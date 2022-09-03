package core

import (
	_ "embed"
	"net/http"
	"os"
	"strconv"
)

//func PushMusicNotify(db MySQLConn, conf *GlobalConfig, blob ConfigBlob, songID int) {
//	plug:= modules.PluginCore{}
//	mus:= CMusic{DB: db, Logger: db.logger, Config: conf, ConfBlob: blob}
//
//}

//go:embed resources/database.sql
var gdpsDatabase string

func HalInitialize(configBlob ConfigBlob, glob *GlobalConfig) {
	db:=MySQLConn{}
	logger:=Logger{}
	//Create DB
	if logger.Should(db.ConnectMultiBlob(configBlob))!=nil {return}
	db.ShouldQuery(gdpsDatabase)
	//Create paths
	os.MkdirAll(glob.SavePath+"/"+configBlob.ServerConfig.SrvID+"/savedata",0777)
	os.Create(glob.SavePath+"/"+configBlob.ServerConfig.SrvID+"levelModel.json")
}

//Count stuff

func CountUsers(db MySQLConn) int {
	var cnt int
	db.ShouldQueryRow("SELECT COUNT(*) as cnt FROM users").Scan(&cnt)
	return cnt
}

func CountLevels(db MySQLConn) int {
	var cnt int
	db.ShouldQueryRow("SELECT COUNT(*) as cnt FROM levels").Scan(&cnt)
	return cnt
}

func CountPosts(db MySQLConn) int {
	var cnt int
	db.ShouldQueryRow("SELECT COUNT(*) as cnt FROM acccomments").Scan(&cnt)
	return cnt
}

func CountComments(db MySQLConn) int {
	var cnt int
	db.ShouldQueryRow("SELECT COUNT(*) as cnt FROM comments").Scan(&cnt)
	return cnt
}

//Trigger stuff

func OnRegister(db MySQLConn, config *GlobalConfig, blob ConfigBlob) bool {
	cnt:=CountUsers(db)
	if cnt>blob.ServerConfig.MaxUsers {return false}
	http.Get(config.ApiEndpoint+"?id="+blob.ServerConfig.SrvID+"&key="+blob.ServerConfig.SrvKey+"&action=stats.users&value="+strconv.Itoa(cnt))
	return true
}

func OnLevel(db MySQLConn, config *GlobalConfig, blob ConfigBlob) bool {
	cnt:=CountLevels(db)
	if cnt>blob.ServerConfig.MaxLevels {return false}
	http.Get(config.ApiEndpoint+"?id="+blob.ServerConfig.SrvID+"&key="+blob.ServerConfig.SrvKey+"&action=stats.levels&value="+strconv.Itoa(cnt))
	return true
}

func OnPost(db MySQLConn, config *GlobalConfig, blob ConfigBlob) bool {
	cnt:=CountPosts(db)
	if cnt>blob.ServerConfig.MaxPosts {return false}
	http.Get(config.ApiEndpoint+"?id="+blob.ServerConfig.SrvID+"&key="+blob.ServerConfig.SrvKey+"&action=stats.posts&value="+strconv.Itoa(cnt))
	return true
}

func OnComment(db MySQLConn, config *GlobalConfig, blob ConfigBlob) bool {
	cnt:=CountComments(db)
	if cnt>blob.ServerConfig.MaxComments {return false}
	http.Get(config.ApiEndpoint+"?id="+blob.ServerConfig.SrvID+"&key="+blob.ServerConfig.SrvKey+"&action=stats.comments&value="+strconv.Itoa(cnt))
	return true
}