package core

import (
	_ "embed"
	"net/http"
	"strconv"
)

//func PushMusicNotify(db MySQLConn, conf *GlobalConfig, blob ConfigBlob, songID int) {
//	plug:= modules.PluginCore{}
//	mus:= CMusic{DB: db, Logger: db.logger, Config: conf, ConfBlob: blob}
//
//}

//Count stuff

func CountUsers(db *MySQLConn) int {
	var cnt int
	db.ShouldQueryRow("SELECT COUNT(*) as cnt FROM #DB#.users").Scan(&cnt)
	return cnt
}

func CountLevels(db *MySQLConn) int {
	var cnt int
	db.ShouldQueryRow("SELECT COUNT(*) as cnt FROM #DB#.levels").Scan(&cnt)
	return cnt
}

func CountPosts(db *MySQLConn) int {
	var cnt int
	db.ShouldQueryRow("SELECT COUNT(*) as cnt FROM #DB#.acccomments").Scan(&cnt)
	return cnt
}

func CountComments(db *MySQLConn) int {
	var cnt int
	db.ShouldQueryRow("SELECT COUNT(*) as cnt FROM #DB#.comments").Scan(&cnt)
	return cnt
}

//Trigger stuff

func OnRegister(db *MySQLConn, config *GlobalConfig, blob ConfigBlob) bool {
	cnt := CountUsers(db)
	if blob.ServerConfig.MaxUsers == -1 {
		http.Get(config.ApiEndpoint + "?srvid=" + blob.ServerConfig.SrvID + "&key=" + blob.ServerConfig.SrvKey + "&action=stats.users&value=" + strconv.Itoa(cnt+1))
		return true
	}
	if cnt > blob.ServerConfig.MaxUsers {
		return false
	}
	http.Get(config.ApiEndpoint + "?srvid=" + blob.ServerConfig.SrvID + "&key=" + blob.ServerConfig.SrvKey + "&action=stats.users&value=" + strconv.Itoa(cnt+1))
	return true
}

func OnLevel(db *MySQLConn, config *GlobalConfig, blob ConfigBlob) bool {
	cnt := CountLevels(db)
	if blob.ServerConfig.MaxLevels == -1 {
		http.Get(config.ApiEndpoint + "?srvid=" + blob.ServerConfig.SrvID + "&key=" + blob.ServerConfig.SrvKey + "&action=stats.levels&value=" + strconv.Itoa(cnt+1))
		return true
	}
	if cnt > blob.ServerConfig.MaxLevels {
		return false
	}
	http.Get(config.ApiEndpoint + "?srvid=" + blob.ServerConfig.SrvID + "&key=" + blob.ServerConfig.SrvKey + "&action=stats.levels&value=" + strconv.Itoa(cnt+1))
	return true
}

func OnPost(db *MySQLConn, config *GlobalConfig, blob ConfigBlob) bool {
	cnt := CountPosts(db)
	if blob.ServerConfig.MaxPosts == -1 {
		http.Get(config.ApiEndpoint + "?srvid=" + blob.ServerConfig.SrvID + "&key=" + blob.ServerConfig.SrvKey + "&action=stats.posts&value=" + strconv.Itoa(cnt+1))
		return true
	}
	if cnt > blob.ServerConfig.MaxPosts {
		return false
	}
	http.Get(config.ApiEndpoint + "?srvid=" + blob.ServerConfig.SrvID + "&key=" + blob.ServerConfig.SrvKey + "&action=stats.posts&value=" + strconv.Itoa(cnt+1))
	return true
}

func OnComment(db *MySQLConn, config *GlobalConfig, blob ConfigBlob) bool {
	cnt := CountComments(db)
	if blob.ServerConfig.MaxComments == -1 {
		http.Get(config.ApiEndpoint + "?srvid=" + blob.ServerConfig.SrvID + "&key=" + blob.ServerConfig.SrvKey + "&action=stats.comments&value=" + strconv.Itoa(cnt+1))
		return true
	}
	if cnt > blob.ServerConfig.MaxComments {
		return false
	}
	http.Get(config.ApiEndpoint + "?srvid=" + blob.ServerConfig.SrvID + "&key=" + blob.ServerConfig.SrvKey + "&action=stats.comments&value=" + strconv.Itoa(cnt+1))
	return true
}
