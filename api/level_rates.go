package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
)

func RateDemon(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
			connector.Error("-1", "Invalid Credentials")
			return
		}
		var lvl_id int
		core.TryInt(&lvl_id, Post.Get("levelID"))
		cl := core.CLevel{DB: db, Id: lvl_id}
		if !cl.Exists(cl.Id) {
			connector.Error("-1", "Level not found")
			return
		}
		role := xacc.GetRoleObj(true)
		var diff, mode int
		core.TryInt(&mode, Post.Get("mode"))
		core.TryInt(&diff, Post.Get("rating"))
		if len(role.Privs) > 0 && role.Privs["aRateDemon"] > 0 && mode != 0 {
			cl.RateDemon(diff % 6)
			if config.ServerConfig.EnableModules["discord"] {
				cl.LoadMain()
				cl.LoadParams()
				cl.LoadStats()
				builder := "[deleted]"
				acc := core.CAccount{DB: db, Uid: cl.Uid}
				if acc.Exists(acc.Uid) {
					acc.LoadAuth(core.CAUTH_UID)
				}
				if len(acc.Uname) > 0 {
					builder = acc.Uname
				}
				data := map[string]string{
					"id":        strconv.Itoa(cl.Id),
					"name":      cl.Name,
					"builder":   builder,
					"diff":      core.DiffToText(cl.StarsGot, cl.DemonDifficulty, cl.IsFeatured, cl.IsEpic),
					"stars":     strconv.Itoa(cl.StarsGot),
					"likes":     strconv.Itoa(cl.Likes),
					"downloads": strconv.Itoa(cl.Downloads),
					"len":       strconv.Itoa(cl.Length),
					"rateuser":  xacc.Uname,
				}
				core.SendAPIWebhook(vars["gdps"], "rate", data)
			}
			connector.NumberedSuccess(cl.Id)
		} else {
			connector.Error("-1", "Insufficient privileges")
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func RateStar(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
			connector.Error("-1", "Invalid Credentials")
			return
		}
		var lvl_id int
		core.TryInt(&lvl_id, Post.Get("levelID"))
		cl := core.CLevel{DB: db, Id: lvl_id}
		if !cl.Exists(cl.Id) {
			connector.Error("-1", "Level not found")
			return
		}
		var diff int
		core.TryInt(&diff, Post.Get("stars"))
		cl.LoadMain()
		cl.DoSuggestDifficulty(diff % 11)
		connector.Success("Difficulty suggested")
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func SuggestStars(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
			connector.Error("-1", "Invalid Credentials")
			return
		}
		var lvl_id int
		core.TryInt(&lvl_id, Post.Get("levelID"))
		cl := core.CLevel{DB: db, Id: lvl_id}
		if !cl.Exists(cl.Id) {
			connector.Error("-1", "Level not found")
			return
		}
		cl.LoadMain()
		role := xacc.GetRoleObj(true)
		var diff, isFeature int
		core.TryInt(&isFeature, Post.Get("feature"))
		core.TryInt(&diff, Post.Get("stars"))
		if len(role.Privs) > 0 {
			if role.Privs["aRateStars"] > 0 {
				diff = diff % 11
				if v, _ := role.Privs["aRateNoDemon"]; v > 0 && diff == 10 {
					connector.Error("-1", "Insufficient privileges")
					return
				}
				cl.RateLevel(diff)
				cl.FeatureLevel(isFeature % 5)
				switch isFeature {
				case 2:
					cl.EpicLevel(true)
				case 3:
					cl.LegendaryLevel(true)
				case 4:
					cl.MythicLevel(true)
				default:
					cl.EpicLevel(false)
				}
				if config.ServerConfig.EnableModules["discord"] {
					cl.LoadMain()
					cl.LoadParams()
					cl.LoadStats()
					builder := "[deleted]"
					acc := core.CAccount{DB: db, Uid: cl.Uid}
					if acc.Exists(acc.Uid) {
						acc.LoadAuth(core.CAUTH_UID)
					}
					if len(acc.Uname) > 0 {
						builder = acc.Uname
					}
					data := map[string]string{
						"id":        strconv.Itoa(cl.Id),
						"name":      cl.Name,
						"builder":   builder,
						"diff":      core.DiffToText(cl.StarsGot, 3, isFeature, cl.IsEpic),
						"stars":     strconv.Itoa(diff),
						"likes":     strconv.Itoa(cl.Likes),
						"downloads": strconv.Itoa(cl.Downloads),
						"len":       strconv.Itoa(cl.Length),
						"rateuser":  xacc.Uname,
					}
					core.SendAPIWebhook(vars["gdps"], "rate", data)
				}
				core.RegisterAction(core.ACTION_LEVEL_RATE, xacc.Uid, cl.Id, map[string]string{
					"uname": xacc.Uname, "type": "StarRate:" + strconv.Itoa(diff%11),
				}, db)
				if isFeature != 0 {
					core.RegisterAction(core.ACTION_LEVEL_RATE, xacc.Uid, cl.Id, map[string]string{
						"uname": xacc.Uname, "type": "Feature",
					}, db)
				}
			} else if role.Privs["aRateReq"] > 0 {
				cl.SendReq(xacc.Uid, diff%11, isFeature)
			} else {
				connector.Error("-1", "Insufficient privileges")
				return
			}
			connector.Success("Rated successfully")
		} else {
			connector.Error("-1", "Insufficient privileges")
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}
