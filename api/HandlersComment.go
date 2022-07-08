package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	"encoding/base64"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func AccountCommentDelete(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("commentID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		var uid int
		core.TryInt(&uid,Post.Get("accountID"))
		xacc:=core.CAccount{DB: db}
		if core.GetGDVersion(Post)==22{
			gjp:=core.ClearGDRequest(Post.Get("gjp2"))
			if !xacc.VerifySession(uid,IPAddr,gjp,true) {
				io.WriteString(resp,"-1")
				return
			}
		}else{
			gjp:=core.ClearGDRequest(Post.Get("gjp"))
			if !xacc.VerifySession(uid,IPAddr,gjp,false) {
				io.WriteString(resp,"-1")
				return
			}
		}
		cc:=core.CComment{DB: db}
		var id int
		core.TryInt(&id,Post.Get("commentID"))
		cc.DeleteAccComment(id, uid)
		io.WriteString(resp,"1")
	}else{
		io.WriteString(resp,"-1")
	}
}

func AccountCommentGet(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Get("accountID")!="" {
		page:=0
		if Post.Has("page") {
			if c, err:= strconv.Atoi(Post.Get("page")); err==nil {page=c}
		}
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		var uid int
		core.TryInt(&uid,Post.Get("accountID"))
		cc:=core.CComment{DB: db}
		comments:=cc.GetAllAccComments(uid,page)
		if len(comments)==0 {
			io.WriteString(resp,"#0:0:0")
		}else{
			output:=""
			for _,comm:= range comments {
				output+=connectors.GetAccountComment(comm)
			}
			io.WriteString(resp,output[:len(output)-1]+"#"+strconv.Itoa(cc.CountAccComments(uid))+":"+strconv.Itoa(page*10)+":10")
		}

	}else{
		io.WriteString(resp,"-1")
	}
}

func AccountCommentUpload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("comment")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		var uid int
		core.TryInt(&uid,Post.Get("accountID"))
		xacc:=core.CAccount{DB: db}
		if core.GetGDVersion(Post)==22{
			gjp:=core.ClearGDRequest(Post.Get("gjp2"))
			if !xacc.VerifySession(uid,IPAddr,gjp,true) {
				io.WriteString(resp,"-1")
				return
			}
		}else{
			gjp:=core.ClearGDRequest(Post.Get("gjp"))
			if !xacc.VerifySession(uid,IPAddr,gjp,false) {
				io.WriteString(resp,"-1")
				return
			}
		}
		comment:=core.ClearGDRequest(Post.Get("comment"))
		cc:=core.CComment{DB: db, Uid: uid, Comment: comment}
		c:="-1"
		if cc.PostAccComment() {c="1"}
		io.WriteString(resp,c)
	}else{
		io.WriteString(resp,"-1")
	}
}


func CommentDelete(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("commentID")!="" && Post.Get("levelID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		var uid int
		core.TryInt(&uid,Post.Get("accountID"))
		xacc:=core.CAccount{DB: db}
		if core.GetGDVersion(Post)==22{
			gjp:=core.ClearGDRequest(Post.Get("gjp2"))
			if !xacc.VerifySession(uid,IPAddr,gjp,true) {
				io.WriteString(resp,"-1")
				return
			}
		}else{
			gjp:=core.ClearGDRequest(Post.Get("gjp"))
			if !xacc.VerifySession(uid,IPAddr,gjp,false) {
				io.WriteString(resp,"-1")
				return
			}
		}
		cc:=core.CComment{DB: db}
		cl:=core.CLevel{DB: db}
		var id, lvl_id int
		core.TryInt(&id,Post.Get("commentID"))
		core.TryInt(&lvl_id,Post.Get("levelID"))
		if cl.IsOwnedBy(uid) {
			cc.DeleteOwnerLevelComment(id,lvl_id)
		}else{
			cc.DeleteLevelComment(id,uid)
		}
		io.WriteString(resp,"1")
	}else{
		io.WriteString(resp,"-1")
	}
}

func CommentGet(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Get("levelID")!="" {
		page:=0
		if Post.Has("page") {
			if c, err:= strconv.Atoi(Post.Get("page")); err==nil {page=c}
		}
		mode:=false
		if Post.Has("mode") && Post.Get("mode")!="0" {mode=true}
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		var lvlId int
		core.TryInt(&lvlId,Post.Get("levelID"))
		cc:=core.CComment{DB: db}
		comments:=cc.GetAllLevelComments(lvlId,page,mode)
		if len(comments)==0 {
			io.WriteString(resp,"#0:0:0")
		}else{
			output:=""
			for _,comm:= range comments {
				output+=connectors.GetLevelComment(comm)
			}
			io.WriteString(resp,output[:len(output)-1]+"#"+strconv.Itoa(cc.CountLevelComments(lvlId))+":"+strconv.Itoa(page*10)+":10")
		}

	}else{
		io.WriteString(resp,"-1")
	}
}

func CommentGetHistory(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Get("userID")!="" {
		page:=0
		if Post.Has("page") {
			if c, err:= strconv.Atoi(Post.Get("page")); err==nil {page=c}
		}
		mode:=false
		if Post.Has("mode") && Post.Get("mode")!="0" {mode=true}
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		acc:=core.CAccount{DB: db}
		core.TryInt(&acc.Uid,Post.Get("userID"))
		if !acc.Exists(acc.Uid) {
			io.WriteString(resp,"-1")
			return
		}
		acc.LoadAuth(core.CAUTH_UID)
		acc.LoadStats()
		acc.LoadVessels()
		role:=acc.GetRoleObj(false)
		cc:=core.CComment{DB: db}
		comments:=cc.GetAllCommentsHistory(acc.Uid,page,mode)
		if len(comments)==0 {
			io.WriteString(resp,"#0:0:0")
		}else{
			output:=""
			for _,comm:= range comments {
				output+=connectors.GetCommentHistory(comm,acc,role)
			}
			io.WriteString(resp,output[:len(output)-1]+"#"+strconv.Itoa(cc.CountCommentHistory(acc.Uid))+":"+strconv.Itoa(page*10)+":10")
		}

	}else{
		io.WriteString(resp,"-1")
	}
}

func CommentUpload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("comment")!="" && Post.Get("levelID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		var uid int
		core.TryInt(&uid,Post.Get("accountID"))
		xacc:=core.CAccount{DB: db}
		if core.GetGDVersion(Post)==22{
			gjp:=core.ClearGDRequest(Post.Get("gjp2"))
			if !xacc.VerifySession(uid,IPAddr,gjp,true) {
				io.WriteString(resp,"-1")
				return
			}
		}else{
			gjp:=core.ClearGDRequest(Post.Get("gjp"))
			if !xacc.VerifySession(uid,IPAddr,gjp,false) {
				io.WriteString(resp,"-1")
				return
			}
		}
		comment:=core.ClearGDRequest(Post.Get("comment"))
		percent:=0
		if Post.Has("percent") {
			core.TryInt(&percent,Post.Get("percent"))
		}
		percent%=101
		cl:=core.CLevel{DB: db}
		core.TryInt(&cl.Id,Post.Get("levelID"))
		if !cl.Exists(cl.Id) {
			io.WriteString(resp,"-1")
			return
		}
		acc:=core.CAccount{DB: db, Uid: uid}
		acc.LoadAuth(core.CAUTH_UID)
		role:=acc.GetRoleObj(true)
		isOwned:=cl.IsOwnedBy(uid)
		if len(role.Privs)>0 || isOwned {
			modCommentByte,err:=base64.StdEncoding.DecodeString(comment)
			modComment:=string(modCommentByte)
			if err==nil && modComment[0]=='!' {
				cl.LoadMain()
				if core.InvokeCommands(db,cl,acc,modComment,isOwned, role) {
					io.WriteString(resp,"1")
				}else{
					io.WriteString(resp,"-1")
				}
			}else{
				cc:=core.CComment{DB: db, Uid: uid, LvlId: cl.Id, Comment: comment, Percent: percent}
				if cc.PostLevelComment() {
					io.WriteString(resp,"1")
				}else{
					io.WriteString(resp,"-1")
				}
			}

		}else{
			cc:=core.CComment{DB: db, Uid: uid, LvlId: cl.Id, Comment: comment, Percent: percent}
			if cc.PostLevelComment() {
				io.WriteString(resp,"1")
			}else{
				io.WriteString(resp,"-1")
			}
		}
	}else{
		io.WriteString(resp,"-1")
	}
}