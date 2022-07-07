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
	if core.CheckGDAuth(Post) && Post.Has("commentID") && Post.Get("commentID")!="" {
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
	if Post.Has("accountID") && Post.Get("accountID")!="" {
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
	if core.CheckGDAuth(Post) && Post.Has("comment") && Post.Get("comment")!="" {
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
	if core.CheckGDAuth(Post) && Post.Has("commentID") && Post.Get("commentID")!="" &&
		Post.Has("levelID") && Post.Get("levelID")!="" {
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
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func CommentGetHistory(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func CommentUpload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}