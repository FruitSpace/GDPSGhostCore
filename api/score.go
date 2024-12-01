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

	//Post:=ReadPost(req)
	db := &core.MySQLConn{}

	if logger.Should(db.ConnectBlob(config)) != nil {
		serverError(connector)
		return
	}
	acc := core.CAccount{DB: db}
	users := acc.GetLeaderboard(core.CLEADERBOARD_BY_CPOINTS, []string{}, 0, config.ServerConfig.TopSize)
	if len(users) == 0 {
		connector.Error("-2", "No users found")
	} else {
		xacc := core.CAccount{DB: db}
		connector.Score_GetLeaderboard(users, xacc)
	}
}

func GetLevelScores(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if core.CheckGDAuth(Post) && Post.Get("levelID") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid credentials")
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
			connector.Error("-2", "No scores found")
			return
		}
		connector.Score_GetScores(scores, "default")
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func GetLevelPlatScores(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if core.CheckGDAuth(Post) && Post.Get("levelID") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid credentials")
			return
		}

		cs := core.CScores{DB: db}

		var percent, time, points, attempts, coins, lvlId, xtype, mode int
		core.TryInt(&lvlId, Post.Get("levelID"))
		core.TryInt(&xtype, Post.Get("type"))
		core.TryInt(&mode, Post.Get("mode"))
		core.TryInt(&percent, Post.Get("percent"))
		core.TryInt(&points, Post.Get("points"))
		core.TryInt(&time, Post.Get("time"))
		core.TryInt(&attempts, Post.Get("s1"))
		core.TryInt(&coins, Post.Get("s9"))

		// COINS = POINTS
		// ATTEMPTS = TIME
		//1: Username
		//2: playerID
		//3: время прохождения в миллисекундах или поинты
		//6: ранг
		//9: иконка
		//10: цвет1
		//11: цвет2
		//14: тип иконки
		//15: special
		//16: accountID
		//42: как давно
		percent = int(math.Abs(float64(percent))) % 101
		attempts = int(math.Abs(float64(attempts)))
		if percent > 0 && attempts > 0 {
			// Upload score
			cs.Attempts = time
			cs.Coins = points
			cs.Uid = xacc.Uid
			cs.LvlId = lvlId
			cs.Percent = percent
			if cs.ScoreExistsByUid(xacc.Uid, lvlId) {
				cs.UpdateLevelScore()
			} else {
				cs.UploadLevelScore()
			}
		}
		//Retrieve all scores
		scores := cs.GetScoresForPlatformerLevelId(lvlId, xtype%4+500, mode == 1, xacc)
		if len(scores) == 0 {
			connector.Error("-2", "No scores found")
			return
		}
		modes := "attempts"
		if mode == 1 {
			modes = "coins"
		}
		connector.Score_GetScores(scores, modes)
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func GetScores(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	xType := Post.Get("type")
	if xType == "" {
		xType = "top"
	}
	db := &core.MySQLConn{}

	if logger.Should(db.ConnectBlob(config)) != nil {
		serverError(connector)
		return
	}
	acc := core.CAccount{DB: db}
	var users []int
	switch xType {
	case "relative":
		if !acc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid credentials")
			return
		}
		acc.LoadStats()
		users = acc.GetLeaderboard(core.CLEADERBOARD_GLOBAL, []string{}, acc.Stars, config.ServerConfig.TopSize)
	case "friends":
		if !acc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid credentials")
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
		connector.Error("-2", "No users found")
	} else {
		connector.Score_GetLeaderboard(users, acc)
	}
}

func UpdateUserScore(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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

	if conf.MaintenanceMode {
		config.SecurityConfig.DisableProtection = false
	}

	if core.CheckIPBan(IPAddr, config) {
		connector.Error("-1", "Banned")
		return
	}

	Post := ReadPost(req)
	if core.CheckGDAuth(Post) {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Success("Invalid Credentials, but as per Geometry Dash API we should return 1 no matter what")
			return
		}
		xacc.LoadStats()
		core.TryInt(&xacc.ColorPrimary, Post.Get("color1"))
		core.TryInt(&xacc.ColorSecondary, Post.Get("color2"))
		core.TryInt(&xacc.ColorGlow, Post.Get("color3"))
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
		core.TryInt(&xacc.Swing, Post.Get("accSwing"))
		core.TryInt(&xacc.Jetpack, Post.Get("accJetpack"))
		core.TryInt(&xacc.Trace, Post.Get("accGlow"))
		core.TryInt(&xacc.Death, Post.Get("accExplosion"))

		protect := core.CProtect{DB: db, Savepath: conf.SavePath + "/" + vars["gdps"], DisableProtection: config.SecurityConfig.DisableProtection}
		protect.LoadModel(conf, config)
		if !protect.DetectStats(xacc.Uid, xacc.Stars, xacc.Diamonds, xacc.Demons, xacc.Coins, xacc.UCoins) {
			connector.Error("-1", "Invalid stats breh") // Le trolling
			return
		}

		// 2.2 demon stats
		{
			core.TryInt(&xacc.ExtraData.DemonStats.Weeklies, Post.Get("dinfow"))
			core.TryInt(&xacc.ExtraData.DemonStats.Gauntlets, Post.Get("dinfog"))
			cf := &core.CLevelFilter{DB: db}
			data := cf.CountDemonTypes(core.Decompose(
				core.CleanDoubles(core.ClearGDRequest(Post.Get("dinfo")), ","),
				","))

			xacc.ExtraData.DemonStats.Standard.Easy = data.Standard.Easy
			xacc.ExtraData.DemonStats.Standard.Medium = data.Standard.Medium
			xacc.ExtraData.DemonStats.Standard.Hard = data.Standard.Hard
			xacc.ExtraData.DemonStats.Standard.Insane = data.Standard.Insane
			xacc.ExtraData.DemonStats.Standard.Extreme = data.Standard.Extreme

			xacc.ExtraData.DemonStats.Platformer.Easy = data.Platformer.Easy
			xacc.ExtraData.DemonStats.Platformer.Medium = data.Platformer.Medium
			xacc.ExtraData.DemonStats.Platformer.Hard = data.Platformer.Hard
			xacc.ExtraData.DemonStats.Platformer.Insane = data.Platformer.Insane
			xacc.ExtraData.DemonStats.Platformer.Extreme = data.Platformer.Extreme
		}

		// 2.2 standard stats
		{
			core.TryInt(&xacc.ExtraData.StandardStats.Daily, Post.Get("sinfod"))
			core.TryInt(&xacc.ExtraData.StandardStats.Gauntlet, Post.Get("sinfog"))
			sinfo := core.Decompose(core.CleanDoubles(core.ClearGDRequest(Post.Get("sinfo")), ","), ",")
			if len(sinfo) == 12 {
				xacc.ExtraData.StandardStats.Auto = sinfo[0]
				xacc.ExtraData.StandardStats.Easy = sinfo[1]
				xacc.ExtraData.StandardStats.Normal = sinfo[2]
				xacc.ExtraData.StandardStats.Hard = sinfo[3]
				xacc.ExtraData.StandardStats.Harder = sinfo[4]
				xacc.ExtraData.StandardStats.Insane = sinfo[5]
				xacc.ExtraData.PlatformerStats.Auto = sinfo[6]
				xacc.ExtraData.PlatformerStats.Easy = sinfo[7]
				xacc.ExtraData.PlatformerStats.Normal = sinfo[8]
				xacc.ExtraData.PlatformerStats.Hard = sinfo[9]
				xacc.ExtraData.PlatformerStats.Harder = sinfo[10]
				xacc.ExtraData.PlatformerStats.Insane = sinfo[11]
				d := 0
				for _, v := range sinfo {
					d += v
				}
				if xacc.Demons > d {
					// Still have no idea why
					xacc.ExtraData.StandardStats.Hard += min(xacc.Demons-d, 5)
				}
			}
		}

		xacc.PushVessels()
		xacc.PushStatsAndExtra()
		connector.NumberedSuccess(xacc.Uid)
	} else {
		connector.Success("Bad Request, but as per Geometry Dash API we should return 1 no matter what")
	}
}
