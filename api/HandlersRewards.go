package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	"encoding/base64"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetChallenges(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Get("chk")!="" && Post.Get("udid")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		cq:=core.CQuests{DB: db}
		if cq.Exists(core.QUEST_TYPE_CHALLENGE) {
			chalk, _ := base64.StdEncoding.DecodeString(Post.Get("chk")[5:])
			chk := core.DoXOR(string(chalk), "19847")
			var uid int
			core.TryInt(&uid,Post.Get("accountID"))
			io.WriteString(resp,connectors.ChallengesOutput(cq,uid,chk,Post.Get("udid")))
		}else{
			io.WriteString(resp,"-2")
		}
	}else{
		io.WriteString(resp,"-1")
	}
}

func GetRewards(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if core.CheckGDAuth(Post) && len(Post.Get("chk"))>5 && Post.Get("udid")!="" {
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		xacc:=core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr){
			io.WriteString(resp,"-1")
			return
		}
		var chestType int
		core.TryInt(&chestType,Post.Get("rewardType"))
		chestType%=3 //Strip to 2 options
		xacc.LoadChests()
		chalk,_:=base64.StdEncoding.DecodeString(Post.Get("chk")[5:])
		chk:=core.DoXOR(string(chalk),"59182")

		chestSmallLeft:=core.MaxInt(0, config.ChestConfig.ChestSmallWait-100+xacc.ChestSmallTime-int(time.Now().Unix())) //!+10800
		chestBigLeft:=core.MaxInt(0, config.ChestConfig.ChestBigWait-100+xacc.ChestBigTime-int(time.Now().Unix())) //!+10800
		switch chestType {
		case 1:
			if chestSmallLeft==0 {
				xacc.ChestSmallCount++
				xacc.ChestSmallTime=int(time.Now().Unix())
				xacc.PushChests()
				chestSmallLeft=config.ChestConfig.ChestSmallWait
			}else{
				io.WriteString(resp,"-1")
				return
			}
			break
		case 2:
			if chestBigLeft==0 {
				xacc.ChestBigCount++
				xacc.ChestBigTime=int(time.Now().Unix())
				xacc.PushChests()
				chestBigLeft=config.ChestConfig.ChestBigWait
			}else{
				io.WriteString(resp,"-1")
				return
			}
			break
		}
		io.WriteString(resp,connectors.ChestOutput(xacc,config,Post.Get("udid"),chk,chestSmallLeft,chestBigLeft,chestType))
	}else{
		io.WriteString(resp,"-1")
	}
}