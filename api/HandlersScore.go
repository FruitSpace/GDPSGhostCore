package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetCreators(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	//Post:=ReadPost(req)
	db := core.MySQLConn{}
	if logger.Should(db.ConnectBlob(config)) != nil {return}
	acc:=core.CAccount{DB: db}
	users:=acc.GetLeaderboard(core.CLEADERBOARD_BY_CPOINTS,[]string{},0)
	if len(users)==0 {
		io.WriteString(resp,"-2")
	}else{
		var lk int
		out:=""
		for _,user:=range users {
			xacc:=core.CAccount{DB: db,Uid: user}
			lk++
			out+=connectors.GetAccLeaderboardItem(xacc,lk)
		}
		io.WriteString(resp,out[:len(out)-1])
	}
}

func GetLevelScores(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func GetScores(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	xType:=Post.Get("type")
	if xType=="" {xType="top"}
	db := core.MySQLConn{}
	if logger.Should(db.ConnectBlob(config)) != nil {return}
	acc:=core.CAccount{DB: db}
	var users []int
	switch xType {
	case "relative":
	case "friends":
	case "creators":
		users=acc.GetLeaderboard(core.CLEADERBOARD_BY_CPOINTS,[]string{},0)
		break
	default:
		users=acc.GetLeaderboard(core.CLEADERBOARD_BY_STARS,[]string{},0)
	}
	if len(users)==0 {
		io.WriteString(resp,"-2")
	}else{
		var lk int
		out:=""
		for _,user:=range users {
			xacc:=core.CAccount{DB: db,Uid: user}
			lk++
			out+=connectors.GetAccLeaderboardItem(xacc,lk)
		}
		io.WriteString(resp,out[:len(out)-1])
	}
}

func UpdateUserScore(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}