package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	"encoding/base64"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"time"
)

func GetChallenges(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	config, err := conf.LoadById(vars["gdps"])
	if logger.Should(err) != nil {
		connector.Error("-1", "Not Found")
		return
	}
	if core.CheckIPBan(IPAddr, config) {
		connector.Error("-1", "Banned")
		return
	}

	Post := ReadPost(req)
	if Post.Get("chk") != "" && Post.Get("udid") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		cq := core.CQuests{DB: db}
		if cq.Exists(core.QUEST_TYPE_CHALLENGE) {
			chalk, _ := base64.StdEncoding.DecodeString(Post.Get("chk")[5:])
			chk := core.DoXOR(string(chalk), "19847")
			var uid int
			core.TryInt(&uid, Post.Get("accountID"))
			connector.Rewards_ChallengesOutput(cq, uid, chk, Post.Get("udid"))
		} else {
			connector.Error("-2", "Challenge not found")
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func GetRewards(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	config, err := conf.LoadById(vars["gdps"])
	if logger.Should(err) != nil {
		connector.Error("-1", "Not Found")
		return
	}
	if core.CheckIPBan(IPAddr, config) {
		connector.Error("-1", "Banned")
		return
	}

	Post := ReadPost(req)
	if core.CheckGDAuth(Post) && len(Post.Get("chk")) > 5 && Post.Get("udid") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		var chestType int
		core.TryInt(&chestType, Post.Get("rewardType"))
		chestType %= 3 //Strip to 2 options
		xacc.LoadChests()
		chalk, _ := base64.StdEncoding.DecodeString(Post.Get("chk")[5:])
		chk := core.DoXOR(string(chalk), "59182")

		chestSmallLeft := core.MaxInt(0, config.ChestConfig.ChestSmallWait-100+xacc.ChestSmallTime-int(time.Now().Unix())) //!+10800
		chestBigLeft := core.MaxInt(0, config.ChestConfig.ChestBigWait-100+xacc.ChestBigTime-int(time.Now().Unix()))       //!+10800
		switch chestType {
		case 1:
			if chestSmallLeft == 0 {
				xacc.ChestSmallCount++
				xacc.ChestSmallTime = int(time.Now().Unix())
				xacc.PushChests()
				chestSmallLeft = config.ChestConfig.ChestSmallWait
			} else {
				connector.Error("-1", "Small chest is not ready yet")
				return
			}
		case 2:
			if chestBigLeft == 0 {
				xacc.ChestBigCount++
				xacc.ChestBigTime = int(time.Now().Unix())
				xacc.PushChests()
				chestBigLeft = config.ChestConfig.ChestBigWait
			} else {
				connector.Error("-1", "Big chest is not ready yet")
				return
			}
		}
		connector.Rewards_ChestOutput(xacc, config, Post.Get("udid"), chk, chestSmallLeft, chestBigLeft, chestType)
	} else {
		connector.Error("-1", "Bad Request")
	}
}
