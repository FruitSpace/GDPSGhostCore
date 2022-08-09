package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkLevelGetLevelsC(b *testing.B) {
	conf:=core.GlobalConfig{
		"https://halhost.cc/app/api/gdps_api.php",
		"stdout",
		"null",
		false,
		"localhost",
		"6379",
		"3XTR4OrD1nArY_K3Y_1907",
		7,
		"./",

		map[string]string{"rabbitmq_host":"auto"},
	}
	for i:=0; i< b.N; i++ {
		LevelGetLevelsС(&conf)
	}
}


func ReadPostS(body string) url.Values {
	if len(body)==0 || strings.Count(string(body),"=")==0 { return url.Values{}}
	vals:=make(url.Values)
	pairs:=strings.Split(string(body),"&")
	for _,val:= range pairs {
		if !strings.Contains(val,"=") {continue}
		m:=strings.SplitN(val,"=",2)
		//fmt.Println(m)
		rval,_:=url.QueryUnescape(m[1])
		rkey,_:=url.QueryUnescape(m[0])
		vals[rkey]=append(vals[rkey],rval)
	}
	return vals

}

func LevelGetLevelsС(conf *core.GlobalConfig){
	//resp:=os.Stderr
	IPAddr:="127.0.0.1"
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById("000S")
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPostS("accountID=3&targetAccountID=3&levelID=9&type=2&gameVersion=21&secret=1&str=sonic")

	var mode, page int
	core.TryInt(&mode,Post.Get("type"))
	core.TryInt(&page,Post.Get("page"))

	s:=strconv.Itoa
	Params:= make(map[string]string)
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
				case "-2":
					//! Change switch to array with index %6
					switch Post.Get("demonFilter") {
					case "1":
						Params["demonDiff"]="3"
					case "2":
						Params["demonDiff"]="4"
					case "3":
						Params["demonDiff"]="0"
					case "4":
						Params["demonDiff"]="5"
					case "5":
						Params["demonDiff"]="6"
					default:
						Params["demonDiff"]="0"
					}
					break
				case "1": //EASY
					fallthrough
				case "2": //NORMAL
					fallthrough
				case "3": //HARD
					fallthrough
				case "4": //HARDER
					fallthrough
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
		case "3":
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_TRENDING)
		case "4":
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_LATEST)
		case "5":
			levels=filter.SearchUserLevels(page,Params,false) //User levels (uid in sterm)
		case "6":
			fallthrough
		case "17":
			Params["isFeatured"]="1"
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_LATEST) //Search featured
		case "7":
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_MAGIC) //Magic (New+Old) | Old = >=10k obj & long
		case "10":
			levels=filter.SearchListLevels(page,Params) //List levels (id1,id2,... in sterm)
		case "11":
			Params["star"]="1"
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_LATEST) //Awarded tab
		case "12":
			//Follow levels
			preg, err := regexp.Compile("[^0-9,-]")
			if logger.Should(err) != nil {return}
			Params["followList"]=preg.ReplaceAllString(core.ClearGDRequest(Post.Get("followed")), "")
			if Params["followList"]=="" {break}
			levels=filter.SearchUserLevels(page,Params,true)
		case "13":
			//Friend levels
			xacc:=core.CAccount{DB: db}
			if ! (core.CheckGDAuth(Post) && xacc.PerformGJPAuth(Post, IPAddr)) {break}
			xacc.LoadSocial()
			if xacc.FriendsCount==0 {break}
			fr:=core.CFriendship{DB: db}
			friendships:=core.Decompose(core.CleanDoubles(xacc.FriendshipIds,","),",")
			friends:=[]int{xacc.Uid}
			for _, frid := range friendships {
				id1,id2:=fr.GetFriendByFID(frid)
				fid:=id1
				if id1==xacc.Uid {fid=id2}
				friends=append(friends,fid)
			}
			Params["followList"]=strings.Join(core.ArrTranslate(friends),",")
			levels=filter.SearchUserLevels(page,Params,true)
		case "16":
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_HALL)
		case "21":
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_SAFE_DAILY)
		case "22":
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_SAFE_WEEKLY)
		case "23":
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_SAFE_EVENT)
		default:
			levels=filter.SearchLevels(page,Params,core.CLEVELFILTER_MOSTLIKED)
		}
	}

	//Output, begins!
	if len(levels)==0 {
		//io.WriteString(resp,"-2")
		return
	}

	//fmt.Println(levels)
	//fmt.Println(Params)
	out:=""
	lvlHash:=""
	usrstring:=""
	musStr:=""
	for _,lvl:= range levels {
		cl:=core.CLevel{DB: db, Id: lvl}
		cl.LoadAll()
		lvlS,lvlH,usrH:=connectors.GetLevelSearch(cl, Gauntlet!=0)
		out+=lvlS
		lvlHash+=lvlH
		usrstring+=usrH
		mus:=core.CMusic{DB: db, ConfBlob: config, Config: conf}
		if cl.SongId!=0 &&mus.GetSong(cl.SongId){
			musStr+=connectors.GetMusic(mus)+"~:~"
		}

	}
	if len(musStr)==0 {musStr="lll"}
	//io.WriteString(resp,out[:len(out)-1]+"#"+
	//	usrstring[:len(usrstring)-1]+"#"+
	//	musStr[:len(musStr)-3]+"#"+
	//	s(filter.Count)+":"+s(page*10)+":10#"+
	//	core.HashSolo2(lvlHash))

}