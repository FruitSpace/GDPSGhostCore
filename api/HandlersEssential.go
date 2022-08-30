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

func GetAccountUrl(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,"http://s.halhost.cc/"+vars["gdps"]+"/database")
}

func GetSongInfo(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
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
	if Post.Get("songID")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		mus:=core.CMusic{DB: db, ConfBlob: config, Config: conf}
		var id int
		core.TryInt(&id,Post.Get("songID"))
		if mus.GetSong(id) {
			io.WriteString(resp,connectors.GetMusic(mus))
		}else{
			io.WriteString(resp,"-1")
		}
	}else{
		io.WriteString(resp,"-1")
	}
}

func GetTopArtists(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
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
	db:=core.MySQLConn{}
	if logger.Should(db.ConnectBlob(config))!=nil {return}
	page:=0
	core.TryInt(&page,Post.Get("page"))
	if page<0 {page=0}
	if logger.Should(db.ConnectBlob(config))!=nil {return}
	mus:=core.CMusic{DB: db, ConfBlob: config, Config: conf}
	artists:=mus.GetTopArtists()
	io.WriteString(resp,connectors.GetTopArtists(artists)+"#"+strconv.Itoa(len(artists))+"0:"+strconv.Itoa(len(artists)))
}

func LikeItem(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
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
	if core.CheckGDAuth(Post) && Post.Get("itemID")!="" && Post.Get("type")!=""{
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		var itemId, cType int
		like:=Post.Get("like")=="1"
		core.TryInt(&itemId,Post.Get("itemID"))
		core.TryInt(&cType,Post.Get("type"))
		switch cType {
		case 1:
			cl:=core.CLevel{DB: db}
			if cl.Exists(itemId) {
				likeAction:=core.CLEVEL_ACTION_DISLIKE
				if like {likeAction=core.CLEVEL_ACTION_LIKE}
				cl.LikeLevel(itemId, xacc.Uid, likeAction)
				io.WriteString(resp, "1")
			}else{
				io.WriteString(resp, "-1")
			}
		case 2:
			comm:=core.CComment{DB: db}
			if comm.ExistsLevelComment(itemId) {
				comm.LikeLevelComment(itemId,xacc.Uid,like)
				io.WriteString(resp, "1")
			}else{
				io.WriteString(resp, "-1")
			}
		case 3:
			comm:=core.CComment{DB: db}
			if comm.ExistsAccComment(itemId) {
				comm.LikeAccComment(itemId,xacc.Uid,like)
				io.WriteString(resp, "1")
			}else{
				io.WriteString(resp, "-1")
			}
		default:
			io.WriteString(resp,"-1")
		}
	}else{
		io.WriteString(resp, "-1")
	}
}

func RequestMod(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
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
		role:=xacc.GetRoleObj(true)
		if len(role.Privs)>0 && role.Privs["aReqMod"]>0 {
			io.WriteString(resp,"1")
		}else{
			io.WriteString(resp,"-1")
		}
	}else{
		io.WriteString(resp, "-1")
	}
}