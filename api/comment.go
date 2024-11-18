package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	"encoding/base64"
	"fmt"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
)

func AccountCommentDelete(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if core.CheckGDAuth(Post) && Post.Get("commentID") != "" {
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
		cc := core.CComment{DB: db}
		var id int
		core.TryInt(&id, Post.Get("commentID"))
		cc.DeleteAccComment(id, xacc.Uid)
		core.OnPost(db, conf, config)
		connector.Success("Comment deleted")
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func AccountCommentGet(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if Post.Get("accountID") != "" {
		page := 0
		core.TryInt(&page, Post.Get("page"))
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		var uid int
		xuid := Post.Get("accountID")
		fmt.Println(Post["accountID"])
		if len(Post["accountID"]) > 1 {
			xuid = Post["accountID"][1]
		}
		core.TryInt(&uid, xuid)
		cc := core.CComment{DB: db}
		comments := cc.GetAllAccComments(uid, page)
		connector.Comment_AccountGet(comments, cc.CountAccComments(uid), page)
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func AccountCommentUpload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if core.CheckGDAuth(Post) && Post.Get("comment") != "" {
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
		comment := core.ClearGDRequest(Post.Get("comment"))
		cc := core.CComment{DB: db, Uid: xacc.Uid, Comment: comment}
		protect := core.CProtect{DB: db, DisableProtection: config.SecurityConfig.DisableProtection}
		if !core.OnPost(db, conf, config) {
			connector.Error("-1", "Post limits exceeded")
			return
		}
		if protect.DetectPosts(xacc.Uid) && cc.PostAccComment() {
			connector.Success("Comment posted")
			return
		}
		serverError(connector)
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func CommentDelete(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if core.CheckGDAuth(Post) && Post.Get("commentID") != "" && Post.Get("levelID") != "" {
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
		cc := core.CComment{DB: db}
		var id, lvl_id int
		core.TryInt(&id, Post.Get("commentID"))
		core.TryInt(&lvl_id, Post.Get("levelID"))
		if lvl_id > 0 {
			// levels
			cl := core.CLevel{DB: db, Id: lvl_id}
			if cl.IsOwnedBy(xacc.Uid) {
				cc.DeleteOwnerLevelComment(id, lvl_id)
			} else {
				cc.DeleteLevelComment(id, xacc.Uid)
			}
		} else {
			// RobTop is retarded
			cl := core.CLevelList{DB: db, ID: lvl_id}
			if cl.IsOwnedBy(xacc.Uid) {
				cc.DeleteOwnerLevelComment(id, lvl_id)
			} else {
				cc.DeleteLevelComment(id, xacc.Uid)
			}
		}
		core.OnComment(db, conf, config)
		connector.Success("Comment deleted")
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func CommentGet(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
		page := 0
		core.TryInt(&page, Post.Get("page"))
		mode := false
		if Post.Get("mode") != "0" {
			mode = true
		}
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		var lvlId int
		core.TryInt(&lvlId, Post.Get("levelID"))
		cc := core.CComment{DB: db}
		comments := cc.GetAllLevelComments(lvlId, page, mode)
		connector.Comment_LevelGet(comments, cc.CountLevelComments(lvlId), page)
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func CommentGetHistory(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if Post.Get("userID") != "" {
		page := 0
		core.TryInt(&page, Post.Get("page"))
		mode := false
		if Post.Get("mode") != "0" {
			mode = true
		}
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		acc := core.CAccount{DB: db}
		core.TryInt(&acc.Uid, Post.Get("userID"))
		if !acc.Exists(acc.Uid) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		acc.LoadAuth(core.CAUTH_UID)
		acc.LoadStats()
		acc.LoadVessels()
		role := acc.GetRoleObj(false)
		cc := core.CComment{DB: db}
		comments := cc.GetAllCommentsHistory(acc.Uid, page, mode)
		connector.Comment_HistoryGet(comments, acc, role, cc.CountCommentHistory(acc.Uid), page)
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func CommentUpload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if core.CheckGDAuth(Post) && Post.Get("comment") != "" && Post.Get("levelID") != "" {
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
		comment := core.ClearGDRequest(Post.Get("comment"))
		percent := 0
		core.TryInt(&percent, Post.Get("percent"))
		percent %= 101
		cl := core.CLevel{DB: db}
		core.TryInt(&cl.Id, Post.Get("levelID"))
		if cl.Id < 0 {
			// CLevelList
			list := core.CLevelList{DB: db, ID: cl.Id * -1}
			if !list.Exists(list.ID) {
				connector.Error("-1", "Invalid LevelList ID")
				return
			}
			cc := core.CComment{DB: db, Uid: xacc.Uid, LvlId: cl.Id, Comment: comment, Percent: percent}
			protect := core.CProtect{DB: db, DisableProtection: config.SecurityConfig.DisableProtection}
			if !core.OnComment(db, conf, config) {
				connector.Error("-1", "Comment limit exceeded")
				return
			}
			if protect.DetectComments(xacc.Uid) && cc.PostLevelComment() {
				connector.Success("Comment posted")
			} else {
				serverError(connector)
			}
			return
		}

		// Next
		if !cl.Exists(cl.Id) {
			connector.Error("-1", "Invalid Level ID")
			return
		}
		acc := core.CAccount{DB: db, Uid: xacc.Uid}
		acc.LoadAuth(core.CAUTH_UID)
		role := acc.GetRoleObj(true)
		isOwned := cl.IsOwnedBy(xacc.Uid)
		if len(role.Privs) > 0 || isOwned {
			modCommentByte, err := base64.StdEncoding.DecodeString(comment)
			modComment := string(modCommentByte)
			if err == nil && modComment[0] == '!' {
				cl.LoadMain()
				cl.LoadParams()
				cmdRes := core.InvokeCommands(db, cl, acc, modComment, isOwned, role)
				switch cmdRes {
				case "err":
					connector.Error("-1", cmdRes)
				case "ok":
					connector.Success("Command executed")
				default:
					connector.Error("-1", cmdRes)
				}
			} else {
				cc := core.CComment{DB: db, Uid: xacc.Uid, LvlId: cl.Id, Comment: comment, Percent: percent}
				protect := core.CProtect{DB: db, DisableProtection: config.SecurityConfig.DisableProtection}
				if !core.OnComment(db, conf, config) {
					connector.Error("-1", "Comment limit exceeded")
					return
				}
				if protect.DetectComments(xacc.Uid) && cc.PostLevelComment() {
					connector.Success("Comment posted")
				} else {
					serverError(connector)
				}
			}
		} else {
			cc := core.CComment{DB: db, Uid: xacc.Uid, LvlId: cl.Id, Comment: comment, Percent: percent}
			protect := core.CProtect{DB: db, DisableProtection: config.SecurityConfig.DisableProtection}
			if !core.OnComment(db, conf, config) {
				connector.Error("-1", "Comment limit exceeded")
				return
			}
			if protect.DetectComments(xacc.Uid) && cc.PostLevelComment() {
				connector.Success("Comment posted")
			} else {
				serverError(connector)
			}
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}
