package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
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
		core.TryInt(&acc.Uid,Post.Get("accountID"))
		if core.GetGDVersion(Post)==22{
			gjp:=core.ClearGDRequest(Post.Get("gjp2"))
			if !acc.VerifySession(acc.Uid,IPAddr,gjp,true) {
				users=[]int{}
				break
			}
		}else{
			gjp:=core.ClearGDRequest(Post.Get("gjp"))
			if !acc.VerifySession(acc.Uid,IPAddr,gjp,false) {
				users=[]int{}
				break
			}
		}
		acc.LoadStats()
		users=acc.GetLeaderboard(core.CLEADERBOARD_GLOBAL,[]string{},acc.Stars)
		break
	case "friends":
		core.TryInt(&acc.Uid,Post.Get("accountID"))
		if core.GetGDVersion(Post)==22{
			gjp:=core.ClearGDRequest(Post.Get("gjp2"))
			if !acc.VerifySession(acc.Uid,IPAddr,gjp,true) {
				io.WriteString(resp,"-2")
				return
			}
		}else{
			gjp:=core.ClearGDRequest(Post.Get("gjp"))
			if !acc.VerifySession(acc.Uid,IPAddr,gjp,false) {
				io.WriteString(resp,"-2")
				return
			}
		}
		acc.LoadSocial()
		if acc.FriendsCount==0 {
			users=[]int{}
			break
		}
		cf:=core.CFriendship{DB: db}
		frs:=strings.Split(acc.FriendshipIds,",")
		var friends []string
		for _,fr:=range frs {
			id,err:=strconv.Atoi(fr)
			if err!=nil {continue}
			uid1,uid2:=cf.GetFriendByFID(id)
			if uid1==0 {continue}
			xuid:=uid1
			if acc.Uid==uid1 {xuid=uid2}
			friends=append(friends,strconv.Itoa(xuid))
		}
		friends=append(friends, strconv.Itoa(acc.Uid))
		users=acc.GetLeaderboard(core.CLEADERBOARD_FRIENDS,friends,0)
		break
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