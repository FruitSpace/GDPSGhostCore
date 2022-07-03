package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	"github.com/go-redis/redis/v8"
	gorilla "github.com/gorilla/mux"
	"golang.org/x/exp/slices"
	"io"
	"log"
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
	if err!=nil{
		if err==redis.Nil {return}
		io.WriteString(resp,"There was an error")
		log.Panicln(err.Error())
	}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Has("targetAccountID") && Post.Get("targetAccountID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		if err:=db.ConnectBlob(config); err!=nil {log.Fatalln(err.Error())}
		acc:=core.CAccount{DB: db}
		var uidSelf int
		core.TryInt(&acc.Uid,Post.Get("targetAccountID"))
		if Post.Has("accountID") && Post.Has("gjp") &&
			Post.Get("accountID")!="" && Post.Get("gjp")!="" {
			xacc:=core.CAccount{DB: db}
			core.TryInt(&uidSelf,Post.Get("accountID"))
			gjp:=core.ClearGDRequest(Post.Get("gjp"))
			if !xacc.VerifySession(uidSelf,IPAddr,gjp) {uidSelf=0}
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
		//! ADD CMessages Messages Check
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