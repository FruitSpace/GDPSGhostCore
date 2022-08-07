package api

import (
	"HalogenGhostCore/core"
	"encoding/base64"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetGauntlets(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	//Post:=ReadPost(req)
	db:=core.MySQLConn{}
	if logger.Should(db.ConnectBlob(config))!=nil {return}
	filter:=core.CLevelFilter{DB: db}
	io.WriteString(resp,filter.GetGauntlets())
}

func GetMapPacks(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	db:=core.MySQLConn{}
	if logger.Should(db.ConnectBlob(config))!=nil {return}
	filter:=core.CLevelFilter{DB: db}
	var page int
	core.TryInt(&page,Post.Get("page"))
	io.WriteString(resp,filter.GetMapPacks(page))
}



func LevelDelete(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("levelID")!="" {
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
		var lvl_id int
		core.TryInt(&lvl_id,Post.Get("levelID"))
		cl:=core.CLevel{DB: db, Id: lvl_id}
		if !cl.IsOwnedBy(uid) {
			io.WriteString(resp,"-1")
			return
		}
		cl.DeleteLevel() //!Fetch before that shit
		cl.RecalculateCPoints(uid)
		core.RegisterAction(core.ACTION_LEVEL_DELETE, uid, lvl_id, map[string]string{"uname":xacc.Uname,"type":"Delete:Owner"},db)
		io.WriteString(resp,"1")
	}else{
		io.WriteString(resp,"-1")
	}
}

func LevelDownload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
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
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		var lvl_id, quest_id int
		core.TryInt(&lvl_id,Post.Get("levelID"))
		cl:=core.CLevel{DB: db}
		if lvl_id<0 {
			cq:=core.CQuests{DB: db}
			if !cq.Exists(lvl_id){
				io.WriteString(resp,"-2")
				return
			}
			switch lvl_id {
			case -1:
				lvl_id,quest_id = cq.GetDaily()
				break
			case -2:
				lvl_id,quest_id = cq.GetWeekly()
				break
			case -3:
				lvl_id,quest_id = cq.GetEvent()
				break
			default:
				io.WriteString(resp,"-2")
				return
			}
		}
		cl:=core.CLevel{DB: db, Id: lvl_id}
		if !cl.Exists(cl.Id) {
			io.WriteString(resp,"-1")
			return
		}
		cl.LoadAll()
		cl.OnDownloadLevel()
		var auto int
		if cl.Difficulty<0 {
			auto=1
			cl.Difficulty=0
		}
		passwd:="0"
		if cl.Password!="0" {passwd:=base64.StdEncoding.EncodeToString([]byte(core.DoXOR(cl.Password,"26364")))}
		if core.CheckGDAuth(Post){
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
		}

	}else{
		io.WriteString(resp,"-1")
	}
}


func LevelGetDaily(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	db:=core.MySQLConn{}
	if logger.Should(db.ConnectBlob(config))!=nil {return}
	var xtype int
	if w:=Post.Get("weekly"); w=="1" {xtype=1}
	if t:=Post.Get("type"); t!="" {
		xtype=0
		if t=="1" {xtype=1}
		if t=="2" {xtype=-1}
	}
	cq:=core.CQuests{DB: db}
	io.WriteString(resp,cq.GetSpecialLevel(xtype))
}

func LevelGetLevels(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func LevelReport(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
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
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		var lvl_id int
		core.TryInt(&lvl_id,Post.Get("levelID"))
		cl:=core.CLevel{DB: db, Id: lvl_id}
		if cl.Exists(lvl_id) {cl.ReportLevel()}
		io.WriteString(resp,"1")
	}else{
		io.WriteString(resp,"-1")
	}
}

func LevelUpdateDescription(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func LevelUpload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func RateDemon(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func RateStar(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func SuggestStars(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}