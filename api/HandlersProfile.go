package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	gorilla "github.com/gorilla/mux"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetUserInfo(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Get("targetAccountID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		acc:=core.CAccount{DB: db}
		var uidSelf int
		core.TryInt(&acc.Uid,Post.Get("targetAccountID"))
		if core.CheckGDAuth(Post) {
			xacc:=core.CAccount{DB: db}
			core.TryInt(&uidSelf,Post.Get("accountID"))
			if core.GetGDVersion(Post)==22{
				gjp:=core.ClearGDRequest(Post.Get("gjp2"))
				if !xacc.VerifySession(uidSelf,IPAddr,gjp,true) {uidSelf=0}
			}else{
				gjp:=core.ClearGDRequest(Post.Get("gjp"))
				if !xacc.VerifySession(uidSelf,IPAddr,gjp,false) {uidSelf=0}
			}
		}
		if !acc.Exists(acc.Uid) {
			io.WriteString(resp,"-1")
			return
		}
		acc.LoadAll()
		blacklist:=strings.Split(acc.Blacklist,",")
		if uidSelf>0 && slices.Contains(blacklist,strconv.Itoa(uidSelf)) {
			io.WriteString(resp,"-1")
			return
		}
		cf:=core.CFriendship{DB: db}
		data:=connectors.GetUserProfile(acc,cf.IsAlreadyFriend(acc.Uid,uidSelf))
		if acc.Uid==uidSelf {
			cm:=core.CMessage{DB: db}
			data+=connectors.UserProfilePersonal(cf.CountFriendRequests(acc.Uid,true),cm.CountMessages(acc.Uid,true))
		}
		io.WriteString(resp,data)
	}else{
		io.WriteString(resp,"-1")
	}
}

func GetUserList(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func GetUsers(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func UpdateAccountSettings(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}