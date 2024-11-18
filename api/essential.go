package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
)

func GetAccountUrl(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	vars := gorilla.Vars(req)
	_, _ = io.WriteString(resp, "https://rugd.gofruit.space/"+vars["gdps"]+"/db")
}

func GetSongInfo(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	songid := Post.Get("songID")
	linkmode := false
	if songid == "" {
		songid = req.URL.Query().Get("id")
		linkmode = true
	}
	if songid != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		mus := core.CMusic{DB: db, ConfBlob: config, Config: conf}
		var id int
		core.TryInt(&id, songid)
		if mus.GetSong(id) {
			if linkmode {
				resp.Header().Set("Location", mus.Url)
				resp.WriteHeader(301)
			} else {
				connector.Essential_GetMusic(mus)
			}
		} else {
			connector.Error("-1", "Music Not Found")
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func GetTopArtists(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	page := 0
	core.TryInt(&page, Post.Get("page"))
	if page < 0 {
		page = 0
	}
	if logger.Should(db.ConnectBlob(config)) != nil {
		return
	}
	mus := core.CMusic{DB: db, ConfBlob: config, Config: conf}
	artists := mus.GetTopArtists()
	connector.Essential_GetTopArtists(artists)
}

func LikeItem(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if core.CheckGDAuth(Post) && Post.Get("itemID") != "" && Post.Get("type") != "" {
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
		var itemId, cType int
		like := Post.Get("like") == "1"
		core.TryInt(&itemId, Post.Get("itemID"))
		core.TryInt(&cType, Post.Get("type"))
		switch cType {
		case 1:
			cl := core.CLevel{DB: db}
			if cl.Exists(itemId) {
				likeAction := core.CLEVEL_ACTION_DISLIKE
				if like {
					likeAction = core.CLEVEL_ACTION_LIKE
				}
				cl.LikeLevel(itemId, xacc.Uid, likeAction)
				connector.Success("Liked level")
			} else {
				connector.Error("-1", "Level Not Found")
			}
		case 2:
			comm := core.CComment{DB: db}
			if comm.ExistsLevelComment(itemId) {
				comm.LikeLevelComment(itemId, xacc.Uid, like)
				connector.Success("Liked comment")
			} else {
				connector.Error("-1", "Comment Not Found")
			}
		case 3:
			comm := core.CComment{DB: db}
			if comm.ExistsAccComment(itemId) {
				comm.LikeAccComment(itemId, xacc.Uid, like)
				connector.Success("Liked account comment")
			} else {
				connector.Error("-1", "Account Comment Not Found")
			}
		case 4:
			clist := core.CLevelList{DB: db, ID: itemId}
			if clist.Exists(itemId) {
				likeAction := core.CLEVEL_ACTION_DISLIKE
				if like {
					likeAction = core.CLEVEL_ACTION_LIKE
				}
				clist.LikeList(itemId, xacc.Uid, likeAction)
				connector.Success("Liked list")
			} else {
				connector.Error("-1", "List Not Found")
			}
		default:
			connector.Error("-1", "Invalid Type")
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func RequestMod(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if core.CheckGDAuth(Post) {
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
		role := xacc.GetRoleObj(true)
		if len(role.Privs) > 0 && role.Privs["aReqMod"] > 0 {
			connector.Success("Request approved")
		} else {
			connector.Error("-1", "Insufficient Privileges")
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}
