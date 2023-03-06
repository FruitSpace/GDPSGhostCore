package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	gorilla "github.com/gorilla/mux"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetCreators(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := req.Header.Get("CF-Connecting-IP")
	if IPAddr == "" {
		IPAddr = req.Header.Get("X-Real-IP")
	}
	if IPAddr == "" {
		IPAddr = strings.Split(req.RemoteAddr, ":")[0]
	}
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	config, err := conf.LoadById(vars["gdps"])
	if logger.Should(err) != nil {
		return
	}
	if core.CheckIPBan(IPAddr, config) {
		return
	}
	//Get:=req.URL.Query()
	//Post:=ReadPost(req)
	db := &core.MySQLConn{}
	defer db.CloseDB()
	if logger.Should(db.ConnectBlob(config)) != nil {
		return
	}
	acc := core.CAccount{DB: db}
	users := acc.GetLeaderboard(core.CLEADERBOARD_BY_CPOINTS, []string{}, 0, config.ServerConfig.TopSize)
	if len(users) == 0 {
		io.WriteString(resp, "-2")
	} else {
		var lk int
		out := ""
		for _, user := range users {
			xacc := core.CAccount{DB: db, Uid: user}
			lk++
			out += connectors.GetAccLeaderboardItem(xacc, lk)
		}
		io.WriteString(resp, out[:len(out)-1])
	}
}

func GetLevelScores(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := req.Header.Get("CF-Connecting-IP")
	if IPAddr == "" {
		IPAddr = req.Header.Get("X-Real-IP")
	}
	if IPAddr == "" {
		IPAddr = strings.Split(req.RemoteAddr, ":")[0]
	}
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	config, err := conf.LoadById(vars["gdps"])
	if logger.Should(err) != nil {
		return
	}
	if core.CheckIPBan(IPAddr, config) {
		return
	}
	//Get:=req.URL.Query()
	Post := ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("levelID") != "" {
		db := &core.MySQLConn{}
		defer db.CloseDB()
		if logger.Should(db.ConnectBlob(config)) != nil {
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			io.WriteString(resp, "-1")
			return
		}

		cs := core.CScores{DB: db}

		var percent, attempts, coins, lvlId, mode int
		core.TryInt(&lvlId, Post.Get("levelID"))
		core.TryInt(&mode, Post.Get("mode"))
		core.TryInt(&percent, Post.Get("percent"))
		core.TryInt(&attempts, Post.Get("s1"))
		core.TryInt(&coins, Post.Get("s9"))
		percent = int(math.Abs(float64(percent))) % 101
		attempts = int(math.Abs(float64(attempts)))
		if percent > 0 && attempts > 0 {
			// Upload score
			if attempts < 8355 {
				attempts = 1
			} else {
				attempts -= 8354
			}
			if coins < 5820 {
				coins = 0
			} else {
				coins = (coins - 5819) % 4
			}
			cs.Uid = xacc.Uid
			cs.LvlId = lvlId
			cs.Percent = percent
			cs.Attempts = attempts
			cs.Coins = coins
			if cs.ScoreExistsByUid(xacc.Uid, lvlId) {
				cs.UpdateLevelScore()
			} else {
				cs.UploadLevelScore()
			}
		}
		//Retrieve all scores
		scores := cs.GetScoresForLevelId(lvlId, mode%4+400, xacc)
		if len(scores) == 0 {
			io.WriteString(resp, "-2")
			return
		}
		out := ""
		for _, score := range scores {
			out += connectors.GetLeaderboardScore(score)
		}
		io.WriteString(resp, out[:len(out)-1])
	} else {
		io.WriteString(resp, "-1")
	}
}

func GetScores(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := req.Header.Get("CF-Connecting-IP")
	if IPAddr == "" {
		IPAddr = req.Header.Get("X-Real-IP")
	}
	if IPAddr == "" {
		IPAddr = strings.Split(req.RemoteAddr, ":")[0]
	}
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	config, err := conf.LoadById(vars["gdps"])
	if logger.Should(err) != nil {
		return
	}
	if core.CheckIPBan(IPAddr, config) {
		return
	}
	//Get:=req.URL.Query()
	Post := ReadPost(req)
	xType := Post.Get("type")
	if xType == "" {
		xType = "top"
	}
	db := &core.MySQLConn{}
	defer db.CloseDB()
	if logger.Should(db.ConnectBlob(config)) != nil {
		return
	}
	acc := core.CAccount{DB: db}
	var users []int
	switch xType {
	case "relative":
		if !acc.PerformGJPAuth(Post, IPAddr) {
			io.WriteString(resp, "-1")
			return
		}
		acc.LoadStats()
		users = acc.GetLeaderboard(core.CLEADERBOARD_GLOBAL, []string{}, acc.Stars, config.ServerConfig.TopSize)
	case "friends":
		if !acc.PerformGJPAuth(Post, IPAddr) {
			io.WriteString(resp, "-1")
			return
		}
		acc.LoadSocial()
		if acc.FriendsCount == 0 {
			users = []int{}
			break
		}
		cf := core.CFriendship{DB: db}
		frs := strings.Split(acc.FriendshipIds, ",")
		var friends []string
		for _, fr := range frs {
			id, err := strconv.Atoi(fr)
			if err != nil {
				continue
			}
			uid1, uid2 := cf.GetFriendByFID(id)
			if uid1 == 0 {
				continue
			}
			xuid := uid1
			if acc.Uid == uid1 {
				xuid = uid2
			}
			friends = append(friends, strconv.Itoa(xuid))
		}
		friends = append(friends, strconv.Itoa(acc.Uid))
		users = acc.GetLeaderboard(core.CLEADERBOARD_FRIENDS, friends, 0, config.ServerConfig.TopSize)
	case "creators":
		users = acc.GetLeaderboard(core.CLEADERBOARD_BY_CPOINTS, []string{}, 0, config.ServerConfig.TopSize)
	default:
		users = acc.GetLeaderboard(core.CLEADERBOARD_BY_STARS, []string{}, 0, config.ServerConfig.TopSize)
	}
	if len(users) == 0 {
		io.WriteString(resp, "-1")
	} else {
		var lk int
		out := ""
		for _, user := range users {
			xacc := core.CAccount{DB: db, Uid: user}
			lk++
			out += connectors.GetAccLeaderboardItem(xacc, lk)
		}
		io.WriteString(resp, out[:len(out)-1])
	}
}

func UpdateUserScore(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := req.Header.Get("CF-Connecting-IP")
	if IPAddr == "" {
		IPAddr = req.Header.Get("X-Real-IP")
	}
	if IPAddr == "" {
		IPAddr = strings.Split(req.RemoteAddr, ":")[0]
	}
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	config, err := conf.LoadById(vars["gdps"])
	if logger.Should(err) != nil {
		return
	}
	if core.CheckIPBan(IPAddr, config) {
		return
	}
	//Get:=req.URL.Query()
	Post := ReadPost(req)
	if core.CheckGDAuth(Post) {
		db := &core.MySQLConn{}
		defer db.CloseDB()
		if logger.Should(db.ConnectBlob(config)) != nil {
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			io.WriteString(resp, "1") //! Weird thing
			return
		}
		xacc.LoadStats()
		core.TryInt(&xacc.ColorPrimary, Post.Get("color1"))
		core.TryInt(&xacc.ColorSecondary, Post.Get("color2"))
		core.TryInt(&xacc.Stars, Post.Get("stars"))
		core.TryInt(&xacc.Demons, Post.Get("demons"))
		core.TryInt(&xacc.Diamonds, Post.Get("diamonds"))
		core.TryInt(&xacc.IconType, Post.Get("iconType"))
		core.TryInt(&xacc.Coins, Post.Get("coins"))
		core.TryInt(&xacc.UCoins, Post.Get("userCoins"))
		core.TryInt(&xacc.Moons, Post.Get("moons"))
		core.TryInt(&xacc.Special, Post.Get("special"))
		core.TryInt(&xacc.Cube, Post.Get("accIcon"))
		core.TryInt(&xacc.Ship, Post.Get("accShip"))
		core.TryInt(&xacc.Wave, Post.Get("accDart"))
		core.TryInt(&xacc.Ball, Post.Get("accBall"))
		core.TryInt(&xacc.Ufo, Post.Get("accBird"))
		core.TryInt(&xacc.Robot, Post.Get("accRobot"))
		core.TryInt(&xacc.Spider, Post.Get("accSpider"))
		//!Unconfirmed jetpack and swing
		core.TryInt(&xacc.Swing, Post.Get("accSwing"))
		core.TryInt(&xacc.Jetpack, Post.Get("accJetpack"))

		core.TryInt(&xacc.Trace, Post.Get("accGlow"))
		core.TryInt(&xacc.Death, Post.Get("accExplosion"))
		protect := core.CProtect{DB: db, Savepath: conf.SavePath + "/" + vars["gdps"], DisableProtection: config.SecurityConfig.DisableProtection}
		protect.LoadModel(conf, config)
		if !protect.DetectStats(xacc.Uid, xacc.Stars, xacc.Diamonds, xacc.Demons, xacc.Coins, xacc.UCoins) {
			io.WriteString(resp, "-1")
			return
		}
		xacc.PushVessels()
		xacc.PushStats()
		io.WriteString(resp, strconv.Itoa(xacc.Uid))
	} else {
		io.WriteString(resp, "1") //! Temporary
	}
}
