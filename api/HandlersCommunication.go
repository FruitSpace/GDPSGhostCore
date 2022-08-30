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

func BlockUser(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	if core.CheckIPBan(IPAddr,config) {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post)  && Post.Get("targetAccountID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		var uidTarget int
		core.TryInt(&uidTarget,Post.Get("targetAccountID"))
		if uidTarget>0 {
			xacc.UpdateBlacklist(core.CBLACKLIST_BLOCK,uidTarget)
		}
		io.WriteString(resp,"1")
	}else{
		io.WriteString(resp, "-1")
	}
}

func UnblockUser(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	if core.CheckIPBan(IPAddr,config) {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("targetAccountID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		var uidTarget int
		core.TryInt(&uidTarget,Post.Get("targetAccountID"))
		if uidTarget>0 {
			xacc.UpdateBlacklist(core.CBLACKLIST_UNBLOCK,uidTarget)
		}
		io.WriteString(resp,"1")
	}else{
		io.WriteString(resp, "-1")
	}
}


func FriendAcceptRequest(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	if core.CheckIPBan(IPAddr,config) {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("requestID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		var requestId int
		core.TryInt(&requestId,Post.Get("requestID"))
		if requestId>0 {
			cf:=core.CFriendship{DB: db}
			io.WriteString(resp,strconv.Itoa(cf.AcceptFriendRequest(requestId,xacc.Uid)))
		}else{
			io.WriteString(resp,"-1")
		}
	}else{
		io.WriteString(resp, "-1")
	}
}

func FriendRejectRequest(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	if core.CheckIPBan(IPAddr,config) {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("targetAccountID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		var targetId int
		core.TryInt(&targetId,Post.Get("targetAccountID"))
		if targetId>0 {
			cf:=core.CFriendship{DB: db}
			issender:= Post.Get("isSender")=="1"
			cf.RejectFriendRequestByUid(xacc.Uid,targetId,issender)
		}
		io.WriteString(resp,"1")
	}else{
		io.WriteString(resp, "-1")
	}
}

func FriendGetRequests(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	if core.CheckIPBan(IPAddr,config) {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		page:=0
		core.TryInt(&page,Post.Get("page"))
		getSent:= Post.Get("getSent")=="1"
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		cf:=core.CFriendship{DB: db}
		count, frqs:=cf.GetFriendRequests(xacc.Uid,page,getSent)
		if len(frqs)==0 {
			io.WriteString(resp,"-2")
		}else{
			output:=""
			for _,frq := range frqs {
				output += connectors.GetFriendRequest(frq)
			}
			io.WriteString(resp,output[:len(output)-1]+"#"+strconv.Itoa(count)+":"+strconv.Itoa(page*10)+":10")
		}
	}else{
		io.WriteString(resp, "-1")
	}
}

func FriendReadRequest(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	if core.CheckIPBan(IPAddr,config) {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("requestID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		var requestId int
		core.TryInt(&requestId,Post.Get("requestID"))
		if requestId>0 {
			cf:=core.CFriendship{DB: db}
			cf.ReadFriendRequest(requestId)
		}
		io.WriteString(resp,"1")
	}else{
		io.WriteString(resp, "-1")
	}
}

func FriendRemove(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	if core.CheckIPBan(IPAddr,config) {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("targetAccountID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		var targetId int
		core.TryInt(&targetId,Post.Get("targetAccountID"))
		if targetId>0 {
			cf:=core.CFriendship{DB: db}
			cf.DeleteFriendship(xacc.Uid,targetId)
		}
		io.WriteString(resp,"1")
	}else{
		io.WriteString(resp, "-1")
	}
}

func FriendRequest(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	if core.CheckIPBan(IPAddr,config) {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("toAccountID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		var targetId int
		core.TryInt(&targetId,Post.Get("toAccountID"))
		if targetId>0 {
			cf:=core.CFriendship{DB: db}
			comment:=Post.Get("comment")
			comment=core.ClearGDRequest(comment)
			io.WriteString(resp,strconv.Itoa(cf.RequestFriend(xacc.Uid,targetId,comment)))
		}else{
			io.WriteString(resp,"-1")
		}
	}else{
		io.WriteString(resp, "-1")
	}
}


func MessageDelete(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	if core.CheckIPBan(IPAddr,config) {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("messageID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		var msgId int
		core.TryInt(&msgId,Post.Get("messageID"))
		if msgId>0 {
			cm:=core.CMessage{DB: db, Id: msgId}
			cm.DeleteMessage(xacc.Uid)
		}
		io.WriteString(resp,"1")
	}else{
		io.WriteString(resp, "-1")
	}
}

func MessageGet(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	if core.CheckIPBan(IPAddr,config) {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("messageID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		var msgId int
		core.TryInt(&msgId,Post.Get("messageID"))
		cm:=core.CMessage{DB: db}
		if cm.Exists(msgId) {
			cm.LoadMessageById(msgId)
			if xacc.Uid==cm.UidSrc || xacc.Uid==cm.UidDest {
				io.WriteString(resp,connectors.GetMessage(cm,xacc.Uid))
			}else{
				io.WriteString(resp,"1")
			}
		}else{
			io.WriteString(resp,"-1")
		}

	}else{
		io.WriteString(resp, "-1")
	}
}

func MessageGetAll(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	if core.CheckIPBan(IPAddr,config) {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		page:=0
		core.TryInt(&page,Post.Get("page"))
		getSent:= Post.Get("getSent")=="1"
		cm:=core.CMessage{DB: db}
		count,msgs:=cm.GetMessageForUid(xacc.Uid,page,getSent)
		if len(msgs)==0 {
			io.WriteString(resp, "-2")
		}else{
			output:=""
			for _,msg := range msgs {
				output+=connectors.GetMessageStr(msg,getSent)
			}
			io.WriteString(resp,output[:len(output)-1]+"#"+strconv.Itoa(count)+":"+strconv.Itoa(page*10)+":10")
		}

	}else{
		io.WriteString(resp, "-1")
	}
}

func MessageUpload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	if core.CheckIPBan(IPAddr,config) {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("toAccountID")!="" && Post.Get("body")!=""{
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		var uidDest int
		core.TryInt(&uidDest,Post.Get("toAccountID"))
		body:=core.ClearGDRequest(Post.Get("body"))
		subject:=core.ClearGDRequest(Post.Get("subject"))
		cm:=core.CMessage{
			DB: db,
			UidSrc: xacc.Uid,
			UidDest: uidDest,
			Subject: subject,
			Message: body,
		}
		protect:=core.CProtect{DB: db, DisableProtection: config.SecurityConfig.DisableProtection}
		if protect.DetectMessages(xacc.Uid) && cm.SendMessageObj(){
			io.WriteString(resp, "1")
		}else{
			io.WriteString(resp, "-1")
		}
	}else{
		io.WriteString(resp, "-1")
	}
}