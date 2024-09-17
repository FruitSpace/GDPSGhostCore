package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	gorilla "github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func LevelListDelete(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if core.CheckGDAuth(Post) && Post.Get("listID") != "" {
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
		var list_id int
		core.TryInt(&list_id, Post.Get("listID"))
		cl := core.CLevelList{DB: db, ID: list_id}
		if !cl.IsOwnedBy(xacc.Uid) {
			connector.Error("-1", "List not found or not owned by user")
			return
		}
		cl.DeleteList()
		connector.Success("List deleted")
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func LevelListUpload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if core.CheckGDAuth(Post) && Post.Get("listLevels") != "" {
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
		cl := core.CLevelList{DB: db}

		cl.UID = xacc.Uid
		core.TryInt(&cl.ID, Post.Get("listID"))
		cl.Name = core.ClearGDRequest(Post.Get("listName"))
		cl.Description = core.ClearGDRequest(Post.Get("listDesc"))
		cl.Levels = core.ClearGDRequest(Post.Get("listLevels"))
		core.TryInt(&cl.Difficulty, Post.Get("difficulty"))
		core.TryInt(&cl.Version, Post.Get("listVersion"))
		//core.TryInt(&cl., Post.Get("original")) //???
		core.TryInt(&cl.Unlisted, Post.Get("unlisted"))

		if cl.Name == "" {
			cl.Name = "Unnamed"
		}
		if cl.Version == 0 {
			cl.Version = 1
		}

		cl.Preload() // Because we have no tables

		if cl.IsOwnedBy(xacc.Uid) {
			res := cl.UpdateList()
			if res > 0 {
				connector.Level_UploadList(res)
			} else {
				connector.Error("-1", "Failed to update list")
			}
		} else {
			if !cl.CheckParams() {
				connector.Error("-1", "Invalid parameters")
				return
			}
			res := cl.UploadList()
			if res > 0 {
				connector.Level_UploadList(res)
			} else {
				connector.Error("-1", "Failed to upload list")
			}
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func LevelListSearch(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	log.Printf("%s: %s\n", vars["gdps"], Post)

	var mode, page int
	core.TryInt(&mode, Post.Get("type"))
	core.TryInt(&page, Post.Get("page"))

	Params := make(map[string]string)
	if sterm := Post.Get("str"); sterm != "" {
		Params["sterm"] = core.ClearGDRequest(Post.Get("str"))
	}

	//Difficulty selector
	if diff := Post.Get("diff"); diff != "" {
		preg, err := regexp.Compile("[^0-9,-]")
		if logger.Should(err) != nil {
			connector.Error("-1", "Invalid difficulty filter")
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
				diffl = append(diffl, sdiff)
			}
			Params["diff"] = strings.Join(diffl, ",")
		}
	}

	//Other params

	var star int

	core.TryInt(&star, Post.Get("star"))
	if star != 0 {
		Params["star"] = "1"
	}

	db := &core.MySQLConn{}

	if logger.Should(db.ConnectBlob(config)) != nil {
		serverError(connector)
		return
	}
	filter := core.CLevelListFilter{DB: db}
	var lists []int

	switch Post.Get("type") {
	case "1":
		lists = filter.SearchLists(page, Params, core.CLEVELLISTFILTER_MOSTDOWNLOADED)
	case "3":
		lists = filter.SearchLists(page, Params, core.CLEVELLISTFILTER_TRENDING)
	case "4":
		lists = filter.SearchLists(page, Params, core.CLEVELLISTFILTER_LATEST)
	case "5":
		lists = filter.SearchUserLists(page, Params, false) //User lists (uid in sterm)
	case "7":
		lists = filter.SearchLists(page, Params, core.CLEVELLISTFILTER_MAGIC) // Robtop lobotomy
	case "11":
		lists = filter.SearchLists(page, Params, core.CLEVELLISTFILTER_AWARDED) //Awarded tab
	case "12":
		//Follow levels
		preg, err := regexp.Compile("[^0-9,-]")
		if logger.Should(err) != nil {
			return
		}
		Params["followList"] = preg.ReplaceAllString(core.ClearGDRequest(Post.Get("followed")), "")
		if Params["followList"] == "" {
			break
		}
		lists = filter.SearchUserLists(page, Params, true)
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
		lists = filter.SearchUserLists(page, Params, false)
	case "27":
		lists = filter.SearchLists(page, Params, core.CLEVELLISTFILTER_SENT)

	default:
		lists = filter.SearchLists(page, Params, core.CLEVELLISTFILTER_MOSTLIKED)
	}

	//Output, begins!
	if len(lists) == 0 {
		connector.Error("-2", "No results found")
		return
	}

	listCore := core.CLevelList{DB: db}
	llistsX := listCore.LoadBulkSearch(lists)
	connector.Level_SearchList(lists, llistsX, filter.Count, page)
}
