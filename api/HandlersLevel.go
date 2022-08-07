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
	"regexp"
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
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)

	var mode, page int
	core.TryInt(&mode,Post.Get("type"))
	core.TryInt(&page,Post.Get("page"))

	s:=strconv.Itoa
	var Params map[string]string
	Params["versionGame"]=s(core.GetGDVersion(Post))
	if sterm:=Post.Get("str"); sterm!="" {Params["sterm"]=core.ClearGDRequest(Post.Get("str"))}

	//Difficulty selector
	if diff:=Post.Get("diff"); diff!=""{
		preg, err:= regexp.Compile("[^0-9,-]")
		if logger.Should(err)!=nil {return}
		diff=core.CleanDoubles(preg.ReplaceAllString(diff,""),",")
		if diff!="-" && diff!="," {
			// The real diff filter begins
			difflist:=strings.Split(diff,",")
			var diffl []string
			for _,sdiff := range difflist {
				if sdiff=="" || sdiff=="-" {continue}
				switch sdiff {
				case "-1":
					diffl=append(diffl,"0") //N/A
					break
				case "-2":
					switch Post.Get("demonFilter") {
					case "1":
						Params["demonDiff"]="3"
						break
					case "2":
						Params["demonDiff"]="4"
						break
					case "3":
						Params["demonDiff"]="0"
						break
					case "4":
						Params["demonDiff"]="5"
						break
					case "5":
						Params["demonDiff"]="6"
						break
					}
					break
				case "1": //EASY
				case "2": //NORMAL
				case "3": //HARD
				case "4": //HARDER
				case "5": //INSANE
					diffl=append(diffl,sdiff+"0")
					break
				default:
					diffl=append(diffl,"-1") //AUTO
				}
			}
			Params["diff"]=strings.Join(diffl,",")
		}
	}

	//Other params
	if plen:=Post.Get("len"); plen!="" {
		preg, err := regexp.Compile("[^0-9,-]")
		if logger.Should(err) != nil {return}
		plen = core.CleanDoubles(preg.ReplaceAllString(plen, ""), ",")
		if plen != "-" && plen != "," {
			Params["length"]=plen
		}
	}
	var uncompleted, onlyCompleted, featured, original, twoPlayer, coins, epic, star, noStar, song, Gauntlet int
	core.TryInt(&uncompleted,Post.Get("uncompleted"))
	core.TryInt(&onlyCompleted,Post.Get("onlyCompleted"))
	if uncompleted!=0 {Params["completed"]="0"}
	if onlyCompleted!=0 {Params["completed"]="1"}
	if completed:=Post.Get("completedLevels"); completed!="" {
		preg, err := regexp.Compile("[^0-9,-]")
		if logger.Should(err) != nil {return}
		completed = core.CleanDoubles(preg.ReplaceAllString(completed, ""), ",")
		Params["completedLevels"]=completed
	}else{
		delete(Params,"completed")
	}

	core.TryInt(&featured,Post.Get("featured"))
	if featured!=0 {Params["isFeatured"]="1"}
	core.TryInt(&epic,Post.Get("epic"))
	if epic!=0 {Params["isEpic"]="1"}
	core.TryInt(&original,Post.Get("original"))
	if original!=0 {Params["isOrig"]="1"}
	core.TryInt(&twoPlayer,Post.Get("twoPlayer"))
	if twoPlayer!=0 {Params["is2p"]="1"}
	core.TryInt(&coins,Post.Get("coins"))
	if coins!=0 {Params["coins"]="1"}
	core.TryInt(&star,Post.Get("star"))
	if star!=0 {Params["star"]="1"}
	core.TryInt(&noStar,Post.Get("noStar"))
	if noStar!=0 {Params["star"]="0"}
	core.TryInt(&song,Post.Get("song"))
	if song!=0 {
		if !Post.Has("songCustom"){song*=-1}
		Params["songid"]=strconv.Itoa(song)
	}


	db:=core.MySQLConn{}
	if logger.Should(db.ConnectBlob(config))!=nil {return}
	filter:=core.CLevelFilter{DB: db}
	var levels []int

	core.TryInt(&Gauntlet,Post.Get("gauntlet"))

	if Gauntlet!=0 {
		//get GAU levels
		levels=filter.GetGauntletLevels(Gauntlet)
	}else{
		switch Post.Get("type") {
		case "1":
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_MOSTDOWNLOADED)
			break
		case "3":
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_TRENDING)
			break
		case "4":
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_LATEST)
			break
		case "5":
			levels=filter.SearchUserLevels(page,Params,false) //User levels (uid in sterm)
			break
		case "6":
		case "17":
			Params["isFeatured"]="1"
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_LATEST) //Search featured
			break
		case "7":
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_MAGIC) //Magic (New+Old)
			break //Old = >=10k obj & long
		case "10":
			levels=filter.SearchListLevels(page,Params) //List levels (id1,id2,... in sterm)
			break
		case "11":
			Params["star"]="1"
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_LATEST) //Awarded tab
			break
		case "12":
			//Follow levels
			preg, err := regexp.Compile("[^0-9,-]")
			if logger.Should(err) != nil {return}
			Params["followList"]=preg.ReplaceAllString(core.ClearGDRequest(Post.Get("followed")), "")
			if Params["followList"]=="" {break}
			levels=filter.SearchUserLevels(page,Params,true)
			break
		case "13":
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_MOSTDOWNLOADED)
			break
		default:
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_MOSTLIKED)
		}
	}




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
		desc:=core.ClearGDRequest(Post.Get("levelDesc"))
		if cl.UpdateDescription(desc) {
			io.WriteString(resp, "1")
		}else{
			io.WriteString(resp,"-1")
		}
	}else{
		io.WriteString(resp,"-1")
	}
}

func LevelUpload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("levelString")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		cl:=core.CLevel{DB: db}

		var pwd, is2p, isUnlisted, isFUnlisted, isLDM int
		cl.Uid=xacc.Uid
		cl.VersionGame=core.GetGDVersion(Post)
		cl.StringLevel=core.ClearGDRequest(Post.Get("levelString"))
		cl.Name=core.ClearGDRequest(Post.Get("levelName"))
		if cl.Name=="" {cl.Name="Unnamed"}
		cl.Description=core.ClearGDRequest(Post.Get("levelDesc"))
		core.TryInt(&cl.Version,Post.Get("levelVersion"))
		if cl.Version==0 {cl.Version=1}
		core.TryInt(&cl.Length,Post.Get("levelLength"))
		core.TryInt(&cl.TrackId,Post.Get("audioTrack"))
		core.TryInt(&pwd,Post.Get("audioTrack"))
		cl.Password=strconv.Itoa(pwd)
		core.TryInt(&cl.OrigId,Post.Get("original"))
		core.TryInt(&cl.SongId,Post.Get("songID"))
		core.TryInt(&cl.Objects,Post.Get("objects"))
		core.TryInt(&cl.Ucoins,Post.Get("coins"))
		core.TryInt(&cl.StarsRequested,Post.Get("requestedStars"))
		if cl.StarsRequested==0 {cl.StarsRequested=1}
		core.TryInt(&is2p,Post.Get("original"))
		cl.Is2p=is2p!=0
		core.TryInt(&isUnlisted,Post.Get("levelVersion"))
		core.TryInt(&isFUnlisted,Post.Get("levelVersion"))
		cl.IsUnlisted=isUnlisted%2+isFUnlisted%2
		core.TryInt(&isLDM,Post.Get("ldm"))
		cl.IsLDM=isLDM!=0
		cl.StringExtra=core.ClearGDRequest(Post.Get("extraString"))
		if cl.StringExtra=="" {cl.StringExtra="29_29_29_40_29_29_29_29_29_29_29_29_29_29_29_29"}
		cl.StringLevelInfo=core.ClearGDRequest(Post.Get("levelInfo"))
		core.TryInt(&cl.VersionBinary,Post.Get("binaryVersion"))
		core.TryInt(&cl.Id,Post.Get("levelID"))

		if cl.IsOwnedBy(xacc.Uid){
			res:=cl.UpdateLevel()
			io.WriteString(resp,strconv.Itoa(res))
			if res>0 {
				core.RegisterAction(core.ACTION_LEVEL_UPDATE, xacc.Uid, res, map[string]string{
					"name": cl.Name, "version": strconv.Itoa(cl.Version),
					"objects": strconv.Itoa(cl.Objects), "starsReq": strconv.Itoa(cl.StarsRequested),
				}, db)
				//!Here be plug
			}else{
				io.WriteString(resp,"-1")
			}
		}else{
			res:=cl.UploadLevel()
			io.WriteString(resp,strconv.Itoa(res))
			if res>0 {
				core.RegisterAction(core.ACTION_LEVEL_UPLOAD, xacc.Uid, res, map[string]string{
					"name": cl.Name, "version": strconv.Itoa(cl.Version),
					"objects": strconv.Itoa(cl.Objects), "starsReq": strconv.Itoa(cl.StarsRequested),
				}, db)
				//!Here be plug
			}else{
				io.WriteString(resp,"-1")
			}
		}
		io.WriteString(resp,"1")
	}else{
		io.WriteString(resp,"-1")
	}
}

func RateDemon(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
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
		if !cl.Exists(cl.Id) {
			io.WriteString(resp,"-1")
			return
		}
		role:=xacc.GetRoleObj(true)
		var diff, mode int
		core.TryInt(&mode, Post.Get("mode"))
		core.TryInt(&diff, Post.Get("rating"))
		if len(role.Privs)>0 && role.Privs["aRateDemon"]>0 && mode!=0 {
			cl.RateDemon(diff%6)
			io.WriteString(resp,strconv.Itoa(cl.Id))
		}else{
			io.WriteString(resp,"-1")
		}
	}else{
		io.WriteString(resp,"-1")
	}
}

func RateStar(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
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
		if !cl.Exists(cl.Id) {
			io.WriteString(resp,"-1")
			return
		}
		var diff int
		core.TryInt(&diff, Post.Get("stars"))
		cl.LoadMain()
		cl.DoSuggestDifficulty(diff%11)
		io.WriteString(resp,"1")
	}else{
		io.WriteString(resp,"-1")
	}
}

func SuggestStars(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
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
		if !cl.Exists(cl.Id) {
			io.WriteString(resp,"-1")
			return
		}
		cl.LoadMain()
		role:=xacc.GetRoleObj(true)
		var diff, isFeature int
		core.TryInt(&isFeature, Post.Get("feature"))
		core.TryInt(&diff, Post.Get("stars"))
		if len(role.Privs)>0 {
			if role.Privs["aRateStars"]>0 {
				cl.RateLevel(diff%11)
				cl.FeatureLevel(isFeature!=0)
				core.RegisterAction(core.ACTION_LEVEL_RATE,xacc.Uid,cl.Id, map[string]string{
					"uname": xacc.Uname, "type": "StarRate:"+strconv.Itoa(diff%11),
				},db)
				if isFeature!=0 {
					core.RegisterAction(core.ACTION_LEVEL_RATE,xacc.Uid,cl.Id, map[string]string{
						"uname": xacc.Uname, "type": "Feature",
					},db)
				}
			}else if role.Privs["aRateReq"]>0 {
				cl.SendReq(xacc.Uid,diff%11,isFeature)
			}else{
				io.WriteString(resp,"-1")
				return
			}
			io.WriteString(resp,"1")
		}else{
			io.WriteString(resp,"-1")
		}
	}else{
		io.WriteString(resp,"-1")
	}
}