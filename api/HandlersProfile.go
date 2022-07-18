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
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		var cType int
		core.TryInt(&cType,Post.Get("type"))
		xacc:=core.CAccount{DB: db}
		core.TryInt(&xacc.Uid,Post.Get("accountID"))
		if core.GetGDVersion(Post)==22{
			gjp:=core.ClearGDRequest(Post.Get("gjp2"))
			if !xacc.VerifySession(xacc.Uid,IPAddr,gjp,true) {
				io.WriteString(resp,"-1")
				return
			}
		}else{
			gjp:=core.ClearGDRequest(Post.Get("gjp"))
			if !xacc.VerifySession(xacc.Uid,IPAddr,gjp,false) {
				io.WriteString(resp,"-1")
					return
			}
		}
		xacc.LoadSocial()
		switch cType {
		case 0:
			if xacc.FriendsCount==0 {
				io.WriteString(resp,"-2")
			}else{
				flist:=strings.Split(xacc.FriendshipIds,",")
				out:=""
				for _,fid:= range flist {
					var Xfid int
					core.TryInt(&Xfid,fid)
					cf:=core.CFriendship{DB: db}
					uid1,uid2:=cf.GetFriendByFID(Xfid)
					acc:=core.CAccount{DB: db,Uid: uid1}
					if uid1==xacc.Uid {acc.Uid=uid2}
					acc.LoadAuth(core.CAUTH_UID)
					acc.LoadVessels()
					acc.LoadStats()
					out+=connectors.UserListItem(acc)
				}
				io.WriteString(resp,out[:len(out)-1])
			}
			break
		case 1:
			blacklist:=strings.Split(xacc.Blacklist,",")
			if len(blacklist)==0 {
				io.WriteString(resp,"-2")
			}else{
				out:=""
				for _,buid:= range blacklist {
					acc:=core.CAccount{DB: db}
					core.TryInt(&acc.Uid,buid)
					if acc.Uid==0 {continue}
					acc.LoadAuth(core.CAUTH_UID)
					acc.LoadVessels()
					acc.LoadStats()
					out+=connectors.UserListItem(acc)
				}
				io.WriteString(resp,out[:len(out)-1])
			}
		}
	}else{
		io.WriteString(resp,"-1")
	}
}

func GetUsers(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Get("str")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		acc:=core.CAccount{DB: db}
		acc.Uid=acc.SearchUsers(core.ClearGDRequest(Post.Get("str")))
		if acc.Uid==0 {
			io.WriteString(resp,"-1")
		}else{
			acc.LoadAuth(core.CAUTH_UID)
			acc.LoadVessels()
			acc.LoadStats()
			io.WriteString(resp,connectors.UserSearchItem(acc))
		}
	}else{
		io.WriteString(resp,"-1")
	}
}

func UpdateAccountSettings(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post){
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		core.TryInt(&xacc.Uid,Post.Get("accountID"))
		if core.GetGDVersion(Post)==22{
			gjp:=core.ClearGDRequest(Post.Get("gjp2"))
			if !xacc.VerifySession(xacc.Uid,IPAddr,gjp,true) {
				io.WriteString(resp,"-1")
				return
			}
		}else{
			gjp:=core.ClearGDRequest(Post.Get("gjp"))
			if !xacc.VerifySession(xacc.Uid,IPAddr,gjp,false) {
				io.WriteString(resp,"-1")
				return
			}
		}

		core.TryInt(&xacc.MS,Post.Get("mS"))
		core.TryInt(&xacc.FrS,Post.Get("frS"))
		core.TryInt(&xacc.CS,Post.Get("cS"))
		xacc.Youtube=core.ClearGDRequest(Post.Get("yt"))
		xacc.Twitter=core.ClearGDRequest(Post.Get("twitter"))
		xacc.Twitch=core.ClearGDRequest(Post.Get("twitch"))
		xacc.PushSettings()
		io.WriteString(resp,"1")
	}else{
		io.WriteString(resp,"-1")
	}
}