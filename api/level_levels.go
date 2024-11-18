package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	"encoding/base64"
	gorilla "github.com/gorilla/mux"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func LevelDelete(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
		if !cl.IsOwnedBy(xacc.Uid) {
			connector.Error("-1", "Level not found or not owned by you")
			return
		}
		cl.DeleteLevel() //!Fetch before that shit
		cl.RecalculateCPoints(xacc.Uid)
		core.OnLevel(db, conf, config)
		core.RegisterAction(core.ACTION_LEVEL_DELETE, xacc.Uid, lvl_id, map[string]string{"uname": xacc.Uname, "type": "Delete:Owner"}, db)
		connector.Success("Level deleted")
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func LevelDownload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if Post.Get("levelID") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		var lvl_id, quest_id int
		core.TryInt(&lvl_id, Post.Get("levelID"))
		if lvl_id < 0 {
			cq := core.CQuests{DB: db}
			if !cq.Exists(lvl_id) {
				connector.Error("-2", "Level not found")
				return
			}
			switch lvl_id {
			case -1:
				quest_id, lvl_id = cq.GetDaily()
			case -2:
				quest_id, lvl_id = cq.GetWeekly()
			case -3:
				quest_id, lvl_id = cq.GetEvent()
			default:
				connector.Error("-2", "Level not found")
				return
			}
		}
		cl := core.CLevel{DB: db, Id: lvl_id}
		if !cl.Exists(cl.Id) {
			connector.Error("-2", "Level not found")
			return
		}
		cl.LoadAll()
		cl.OnDownloadLevel()
		passwd := "0"
		phash := cl.Password
		if phash != "0" {
			passwd = base64.StdEncoding.EncodeToString([]byte(core.DoXOR(cl.Password, "26364")))
		}
		if core.CheckGDAuth(Post) {
			var uid int
			core.TryInt(&uid, Post.Get("accountID"))
			xacc := core.CAccount{DB: db}
			if xacc.PerformGJPAuth(Post, IPAddr) {
				role := xacc.GetRoleObj(true)
				if len(role.Privs) > 0 && role.Privs["cLvlAccess"] > 0 {
					passwd = base64.StdEncoding.EncodeToString([]byte(core.DoXOR("1", "26364")))
					phash = "1"
				}
			}
		}

		if cl.SuggestDifficultyCnt > 0 && cl.StarsGot == 0 {
			diffCount := int(math.Round(cl.SuggestDifficulty))
			diffName := "Unspecified"
			//! Change that to array and get %11 index
			switch diffCount {
			case 1:
				diffName = "Auto"
			case 2:
				diffName = "Easy"
			case 3:
				diffName = "Normal"
			case 4:
				fallthrough
			case 5:
				diffName = "Hard"
			case 6:
				fallthrough
			case 7:
				diffName = "Harder"
			case 8:
				fallthrough
			case 9:
				diffName = "Insane"
			case 10:
				diffName = "Demon"
			}
			t, _ := base64.StdEncoding.DecodeString(cl.Description)
			cl.Description = base64.StdEncoding.EncodeToString([]byte(string(t) + " [Suggest: " + diffName + " (" + strconv.Itoa(diffCount) + ")]"))
		}
		connector.Level_GetLevelFull(cl, passwd, phash, quest_id)
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func LevelGetDaily(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	db := &core.MySQLConn{}

	if logger.Should(db.ConnectBlob(config)) != nil {
		serverError(connector)
		return
	}
	var xtype int
	if w := Post.Get("weekly"); w == "1" {
		xtype = 1
	}
	if t := Post.Get("type"); t != "" {
		xtype = 0
		if t == "1" {
			xtype = 1
		}
		if t == "2" {
			xtype = -1
		}
	}
	cq := core.CQuests{DB: db}
	id, left := cq.GetSpecialLevel(xtype)
	connector.Level_GetSpecials(id, left)
}

func LevelGetLevels(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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

	metrics := core.NewGoMetrics()
	defer func() {
		log.Printf("### %s perfromance debug V2\n%s", vars["gdps"], metrics.DumpText())
	}()
	metrics.NewStep("Parsing")

	Post := ReadPost(req)
	//FIXME: Cache doesn't work well with connectors
	// Check cache
	//t := "gd"
	//if req.URL.Query().Has("json") {
	//	t = "json"
	//}
	//cacheKey := fmt.Sprintf("%s/getLevels/%s/%s", vars["gdps"], t, Post.Encode())
	//if res, errx := cached(cacheKey); errx == nil {
	//	io.WriteString(resp, res)
	//	return
	//}

	var mode, page int
	core.TryInt(&mode, Post.Get("type"))
	core.TryInt(&page, Post.Get("page"))

	s := strconv.Itoa
	Params := make(map[string]string)
	Params["versionGame"] = s(core.GetGDVersion(Post))
	if sterm := Post.Get("str"); sterm != "" {
		Params["sterm"] = core.ClearGDRequest(Post.Get("str"))
	}

	//Difficulty selector
	if diff := Post.Get("diff"); diff != "" {
		preg, err := regexp.Compile("[^0-9,-]")
		if logger.Should(err) != nil {
			serverError(connector)
			return
		}
		diff = core.CleanDoubles(preg.ReplaceAllString(diff, ""), ",")
		if diff != "-" && diff != "," {
			// The real diff filter begins
			difflist := strings.Split(diff, ",")
			var diffl []string
			for _, sdiff := range difflist {
				if sdiff == "" || sdiff == "-" {
					continue
				}
				switch sdiff {
				case "-1":
					diffl = append(diffl, "0") //N/A
				case "-2":
					//! Change switch to array with index %6
					switch Post.Get("demonFilter") {
					case "1":
						Params["demonDiff"] = "3"
					case "2":
						Params["demonDiff"] = "4"
					case "3":
						Params["demonDiff"] = "0"
					case "4":
						Params["demonDiff"] = "5"
					case "5":
						Params["demonDiff"] = "6"
					default:
						Params["demonDiff"] = "0"
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
					diffl = append(diffl, sdiff+"0")
					break
				default:
					diffl = append(diffl, "-1") //AUTO
				}
			}
			Params["diff"] = strings.Join(diffl, ",")
		}
	}

	//Other params
	if plen := Post.Get("len"); plen != "" {
		preg, err := regexp.Compile("[^0-9,-]")
		if logger.Should(err) != nil {
			serverError(connector)
			return
		}
		plen = core.CleanDoubles(preg.ReplaceAllString(plen, ""), ",")
		if plen != "-" && plen != "," {
			Params["length"] = plen
		}
	}
	var uncompleted, onlyCompleted, featured, original, twoPlayer, coins, epic, star, noStar, song, Gauntlet int
	core.TryInt(&uncompleted, Post.Get("uncompleted"))
	core.TryInt(&onlyCompleted, Post.Get("onlyCompleted"))
	if uncompleted != 0 {
		Params["completed"] = "0"
	}
	if onlyCompleted != 0 {
		Params["completed"] = "1"
	}
	if completed := Post.Get("completedLevels"); completed != "" {
		preg, err := regexp.Compile("[^0-9,-]")
		if logger.Should(err) != nil {
			serverError(connector)
			return
		}
		completed = core.CleanDoubles(preg.ReplaceAllString(completed, ""), ",")
		Params["completedLevels"] = completed
	} else {
		delete(Params, "completed")
	}

	core.TryInt(&featured, Post.Get("featured"))
	if featured != 0 {
		Params["isFeatured"] = "1"
	}
	core.TryInt(&epic, Post.Get("epic"))
	if epic != 0 {
		Params["isEpic"] = "1"
	}
	core.TryInt(&original, Post.Get("original"))
	if original != 0 {
		Params["isOrig"] = "1"
	}
	core.TryInt(&twoPlayer, Post.Get("twoPlayer"))
	if twoPlayer != 0 {
		Params["is2p"] = "1"
	}
	core.TryInt(&coins, Post.Get("coins"))
	if coins != 0 {
		Params["coins"] = "1"
	}
	core.TryInt(&star, Post.Get("star"))
	if star != 0 {
		Params["star"] = "1"
	}
	core.TryInt(&noStar, Post.Get("noStar"))
	if noStar != 0 {
		Params["star"] = "0"
	}
	core.TryInt(&song, Post.Get("song"))
	if song != 0 {
		if !Post.Has("songCustom") {
			song *= -1
		}
		Params["songid"] = strconv.Itoa(song)
	}

	db := &core.MySQLConn{}

	if logger.Should(db.ConnectBlob(config)) != nil {
		serverError(connector)
		return
	}
	filter := core.CLevelFilter{DB: db}
	var levels []int

	core.TryInt(&Gauntlet, Post.Get("gauntlet"))

	metrics.NewStep("Searching")

	if Gauntlet != 0 {
		//get GAU levels
		levels = filter.GetGauntletLevels(Gauntlet)
	} else {
		switch Post.Get("type") {
		case "1":
			levels = filter.SearchLevels(page, Params, core.CLEVELFILTER_MOSTDOWNLOADED)
		case "3":
			levels = filter.SearchLevels(page, Params, core.CLEVELFILTER_TRENDING)
		case "4":
			levels = filter.SearchLevels(page, Params, core.CLEVELFILTER_LATEST)
		case "5":
			levels = filter.SearchUserLevels(page, Params, false) //User levels (uid in sterm)
		case "6":
			fallthrough
		case "17":
			Params["isFeatured"] = "1"
			levels = filter.SearchLevels(page, Params, core.CLEVELFILTER_LATEST) //Search featured
		case "7":
			levels = filter.SearchLevels(page, Params, core.CLEVELFILTER_MAGIC) //Magic (New+Old) | Old = >=10k obj & long
		case "10":
			fallthrough
		case "19":
			levels = filter.SearchListLevels(page, Params) //List levels (id1,id2,... in sterm)
		case "11":
			Params["star"] = "1"
			levels = filter.SearchLevels(page, Params, core.CLEVELFILTER_LATEST) //Awarded tab
		case "12":
			//Follow levels
			preg, err := regexp.Compile("[^0-9,-]")
			if logger.Should(err) != nil {
				serverError(connector)
				return
			}
			Params["followList"] = preg.ReplaceAllString(core.ClearGDRequest(Post.Get("followed")), "")
			if Params["followList"] == "" {
				break
			}
			levels = filter.SearchUserLevels(page, Params, true)
		case "13":
			//Friend levels
			xacc := core.CAccount{DB: db}
			if !(core.CheckGDAuth(Post) && xacc.PerformGJPAuth(Post, IPAddr)) {
				break
			}
			xacc.LoadSocial()
			if xacc.FriendsCount == 0 {
				break
			}
			fr := core.CFriendship{DB: db}
			friendships := core.Decompose(core.CleanDoubles(xacc.FriendshipIds, ","), ",")
			friends := []int{xacc.Uid}
			for _, frid := range friendships {
				id1, id2 := fr.GetFriendByFID(frid)
				fid := id1
				if id1 == xacc.Uid {
					fid = id2
				}
				friends = append(friends, fid)
			}
			Params["followList"] = strings.Join(core.ArrTranslate(friends), ",")
			levels = filter.SearchUserLevels(page, Params, false)
		case "16":
			levels = filter.SearchLevels(page, Params, core.CLEVELFILTER_HALL)
		case "21":
			levels = filter.SearchLevels(page, Params, core.CLEVELFILTER_SAFE_DAILY)
		case "22":
			levels = filter.SearchLevels(page, Params, core.CLEVELFILTER_SAFE_WEEKLY)
		case "23":
			levels = filter.SearchLevels(page, Params, core.CLEVELFILTER_SAFE_EVENT)
		case "25":
			var lid int // Fuck robtop, that's for getting levels from lists
			core.TryInt(&lid, Post.Get("str"))
			clist := core.CLevelList{DB: db, ID: lid}
			clist.OnDownloadList()
			clist.Load(lid)
			levels = core.Decompose(core.QuickComma(clist.Levels), ",")
		case "27":
			levels = filter.SearchLevels(page, Params, core.CLEVELFILTER_SENT) //SENT
		default:
			levels = filter.SearchLevels(page, Params, core.CLEVELFILTER_MOSTLIKED)
		}
	}

	//Output, begins!
	if len(levels) == 0 {
		metrics.Done()
		connector.Error("-2", "No levels found")

		// FIXME
		//io.WriteString(resp, withCache(cacheKey, "-2"))
		return
	}

	metrics.NewStep("Levels Fetch")
	lvlCore := core.CLevel{DB: db}
	lvlsX := lvlCore.LoadBulkSearch(levels)
	mus := &core.CMusic{DB: db, ConfBlob: config, Config: conf}

	connector.Level_SearchLevels(levels, lvlsX, mus, filter.Count, page, core.GetGDVersion(Post), Gauntlet)

	metrics.Done()

}

func LevelReport(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if Post.Get("levelID") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		var lvl_id int
		core.TryInt(&lvl_id, Post.Get("levelID"))
		cl := core.CLevel{DB: db, Id: lvl_id}
		if cl.Exists(lvl_id) {
			cl.ReportLevel()
		}
		connector.Success("Level was reported")
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func LevelUpdateDescription(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
		if !cl.IsOwnedBy(xacc.Uid) {
			connector.Error("-1", "Level doesn't exist or you don't own it")
			return
		}
		desc := core.ClearGDRequest(Post.Get("levelDesc"))
		if cl.UpdateDescription(desc) {
			connector.Success("Description updated")
		} else {
			connector.Error("-1", "Failed to update description")
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func LevelUpload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if core.CheckGDAuth(Post) && Post.Get("levelString") != "" {
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
		cl := core.CLevel{DB: db}

		var pwd, is2p, isUnlisted, isFUnlisted, isLDM int
		cl.Uid = xacc.Uid
		cl.VersionGame = core.GetGDVersion(Post)
		cl.StringLevel = core.ClearGDRequest(Post.Get("levelString"))
		cl.Name = core.ClearGDRequest(Post.Get("levelName"))
		if cl.Name == "" {
			cl.Name = "Unnamed"
		}
		cl.Description = core.ClearGDRequest(Post.Get("levelDesc"))
		core.TryInt(&cl.Version, Post.Get("levelVersion"))
		if cl.Version == 0 {
			cl.Version = 1
		}
		core.TryInt(&cl.Length, Post.Get("levelLength"))
		core.TryInt(&cl.TrackId, Post.Get("audioTrack"))
		core.TryInt(&pwd, Post.Get("password"))
		cl.Password = strconv.Itoa(pwd)
		core.TryInt(&cl.OrigId, Post.Get("original"))
		core.TryInt(&cl.SongId, Post.Get("songID"))
		core.TryInt(&cl.Objects, Post.Get("objects"))
		core.TryInt(&cl.Ucoins, Post.Get("coins"))
		core.TryInt(&cl.StarsRequested, Post.Get("requestedStars"))
		if cl.StarsRequested == 0 {
			cl.StarsRequested = 1
		}
		core.TryInt(&is2p, Post.Get("original"))
		cl.Is2p = is2p != 0
		core.TryInt(&isUnlisted, Post.Get("unlisted"))
		if unl := Post.Get("unlisted1"); unl != "" {
			core.TryInt(&isUnlisted, unl)
		}
		core.TryInt(&isFUnlisted, Post.Get("unlisted2"))
		cl.IsUnlisted = isUnlisted%2 + isFUnlisted%2
		core.TryInt(&isLDM, Post.Get("ldm"))
		cl.IsLDM = isLDM != 0
		cl.StringExtra = core.ClearGDRequest(Post.Get("extraString"))
		if cl.StringExtra == "" {
			cl.StringExtra = "29_29_29_40_29_29_29_29_29_29_29_29_29_29_29_29"
		}
		cl.StringLevelInfo = core.ClearGDRequest(Post.Get("levelInfo"))
		core.TryInt(&cl.VersionBinary, Post.Get("binaryVersion"))
		core.TryInt(&cl.Id, Post.Get("levelID"))
		cl.StringSettings = core.ClearGDRequest(Post.Get("songIDs")) + ";" + core.ClearGDRequest(Post.Get("sfxIDs"))

		cl.UnlockLevelObject = config.SecurityConfig.NoLevelLimits

		if cl.IsOwnedBy(xacc.Uid) {
			res := cl.UpdateLevel()
			if res > 0 {
				connector.NumberedSuccess(res)
				if config.ServerConfig.EnableModules["discord"] {
					desc, _ := base64.StdEncoding.DecodeString(cl.Description)
					data := map[string]string{
						"id":      strconv.Itoa(res),
						"name":    cl.Name + " (v" + strconv.Itoa(cl.Version) + ")",
						"builder": xacc.Uname,
						"desc":    string(desc),
					}
					core.SendAPIWebhook(vars["gdps"], "newlevel", data)
				}
				core.RegisterAction(core.ACTION_LEVEL_UPDATE, xacc.Uid, res, map[string]string{
					"name": cl.Name, "version": strconv.Itoa(cl.Version),
					"objects": strconv.Itoa(cl.Objects), "starsReq": strconv.Itoa(cl.StarsRequested),
				}, db)
				//!Here be plug
			} else {
				connector.Error("-1", "Level Update Failed")
			}
		} else {
			if !cl.CheckParams() {
				connector.Error("-1", "Invalid level data")
				return
			}
			if !core.OnLevel(db, conf, config) {
				connector.Error("-1", "Level Update Failed")
				return
			}
			protect := core.CProtect{DB: db, Savepath: conf.SavePath + "/" + vars["gdps"], DisableProtection: config.SecurityConfig.DisableProtection}
			protect.LoadModel(conf, config)
			res := -1
			if protect.DetectLevelModel(xacc.Uid) {
				res = cl.UploadLevel()
			}
			if res > 0 {
				connector.NumberedSuccess(res)
				if config.ServerConfig.EnableModules["discord"] {
					desc, _ := base64.StdEncoding.DecodeString(cl.Description)
					data := map[string]string{
						"id":      strconv.Itoa(res),
						"name":    cl.Name,
						"builder": xacc.Uname,
						"desc":    string(desc),
					}
					core.SendAPIWebhook(vars["gdps"], "newlevel", data)
				}
				core.RegisterAction(core.ACTION_LEVEL_UPLOAD, xacc.Uid, res, map[string]string{
					"name": cl.Name, "version": strconv.Itoa(cl.Version),
					"objects": strconv.Itoa(cl.Objects), "starsReq": strconv.Itoa(cl.StarsRequested),
				}, db)
				//!Here be plug
			} else {
				connector.Error("-1", "Level Upload Failed")
			}
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}
