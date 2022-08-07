package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	"encoding/base64"
	gorilla "github.com/gorilla/mux"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
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
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		var lvl_id int
		core.TryInt(&lvl_id,Post.Get("levelID"))
		cl:=core.CLevel{DB: db, Id: lvl_id}
		if !cl.IsOwnedBy(xacc.Uid) {
			io.WriteString(resp,"-1")
			return
		}
		cl.DeleteLevel() //!Fetch before that shit
		cl.RecalculateCPoints(xacc.Uid)
		core.RegisterAction(core.ACTION_LEVEL_DELETE, xacc.Uid, lvl_id, map[string]string{"uname":xacc.Uname,"type":"Delete:Owner"},db)
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
		passwd:="0"
		phash:=cl.Password
		if cl.Password!="0" {passwd=base64.StdEncoding.EncodeToString([]byte(core.DoXOR(cl.Password,"26364")))}
		if core.CheckGDAuth(Post){
			var uid int
			core.TryInt(&uid,Post.Get("accountID"))
			xacc:=core.CAccount{DB: db}
			if xacc.PerformGJPAuth(Post, IPAddr){
				role:=xacc.GetRoleObj(true)
				if len(role.Privs)>0 && role.Privs["aReqMod"]>0 {
					passwd=base64.StdEncoding.EncodeToString([]byte(core.DoXOR("1","26364")))
					phash="1"
				}
			}
		}

		if cl.SuggestDifficultyCnt>0 && cl.StarsGot==0 {
			diffCount:=int(math.Round(cl.SuggestDifficulty))
			diffName:="Unspecified"
			switch diffCount {
			case 1:
				diffName="Auto"
				break
			case 2:
				diffName="Easy"
				break
			case 3:
				diffName="Normal"
				break
			case 4:
			case 5:
				diffName="Hard"
				break
			case 6:
			case 7:
				diffName="Harder"
				break
			case 8:
			case 9:
				diffName="Insane"
				break
			case 10:
				diffName="Demon"
				break
			}
			t,_:=base64.StdEncoding.DecodeString(cl.Description)
			cl.Description=base64.StdEncoding.EncodeToString([]byte(string(t)+" [Suggest: "+diffName+" ("+strconv.Itoa(diffCount)+")]"))
		}
		io.WriteString(resp,connectors.GetLevelFull(cl,passwd,phash,quest_id))
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