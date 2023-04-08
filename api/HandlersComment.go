package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	"encoding/base64"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func AccountCommentDelete(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if core.CheckGDAuth(Post) && Post.Get("commentID") != "" {
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
		cc := core.CComment{DB: db}
		var id int
		core.TryInt(&id, Post.Get("commentID"))
		cc.DeleteAccComment(id, xacc.Uid)
		core.OnPost(db, conf, config)
		io.WriteString(resp, "1")
	} else {
		io.WriteString(resp, "-1")
	}
}

func AccountCommentGet(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if Post.Get("accountID") != "" {
		page := 0
		core.TryInt(&page, Post.Get("page"))
		db := &core.MySQLConn{}
		defer db.CloseDB()
		if logger.Should(db.ConnectBlob(config)) != nil {
			return
		}
		var uid int
		core.TryInt(&uid, Post.Get("accountID"))
		cc := core.CComment{DB: db}
		comments := cc.GetAllAccComments(uid, page)
		if len(comments) == 0 {
			io.WriteString(resp, "#0:0:0")
		} else {
			output := ""
			for _, comm := range comments {
				output += connectors.GetAccountComment(comm)
			}
			io.WriteString(resp, output[:len(output)-1]+"#"+strconv.Itoa(cc.CountAccComments(uid))+":"+strconv.Itoa(page*10)+":10")
		}

	} else {
		io.WriteString(resp, "-1")
	}
}

func AccountCommentUpload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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

	if conf.MaintenanceMode {
		config.SecurityConfig.DisableProtection = false
	}

	if core.CheckIPBan(IPAddr, config) {
		return
	}
	//Get:=req.URL.Query()
	Post := ReadPost(req)
	if core.CheckGDAuth(Post) && Post.Get("comment") != "" {
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
		comment := core.ClearGDRequest(Post.Get("comment"))
		cc := core.CComment{DB: db, Uid: xacc.Uid, Comment: comment}
		protect := core.CProtect{DB: db, DisableProtection: config.SecurityConfig.DisableProtection}
		c := "-1"
		if !core.OnPost(db, conf, config) {
			io.WriteString(resp, "-1")
			return
		}
		if protect.DetectPosts(xacc.Uid) && cc.PostAccComment() {
			c = "1"
		}
		io.WriteString(resp, c)
	} else {
		io.WriteString(resp, "-1")
	}
}

func CommentDelete(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if core.CheckGDAuth(Post) && Post.Get("commentID") != "" && Post.Get("levelID") != "" {
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
		cc := core.CComment{DB: db}
		var id, lvl_id int
		core.TryInt(&id, Post.Get("commentID"))
		core.TryInt(&lvl_id, Post.Get("levelID"))
		cl := core.CLevel{DB: db, Id: lvl_id}
		if cl.IsOwnedBy(xacc.Uid) {
			cc.DeleteOwnerLevelComment(id, lvl_id)
		} else {
			cc.DeleteLevelComment(id, xacc.Uid)
		}
		core.OnComment(db, conf, config)
		io.WriteString(resp, "1")
	} else {
		io.WriteString(resp, "-1")
	}
}

func CommentGet(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if Post.Get("levelID") != "" {
		page := 0
		core.TryInt(&page, Post.Get("page"))
		mode := false
		if Post.Get("mode") != "0" {
			mode = true
		}
		db := &core.MySQLConn{}
		defer db.CloseDB()
		if logger.Should(db.ConnectBlob(config)) != nil {
			return
		}
		var lvlId int
		core.TryInt(&lvlId, Post.Get("levelID"))
		cc := core.CComment{DB: db}
		comments := cc.GetAllLevelComments(lvlId, page, mode)
		if len(comments) == 0 {
			io.WriteString(resp, "#0:0:0")
		} else {
			output := ""
			for _, comm := range comments {
				output += connectors.GetLevelComment(comm)
			}
			io.WriteString(resp, output[:len(output)-1]+"#"+strconv.Itoa(cc.CountLevelComments(lvlId))+":"+strconv.Itoa(page*10)+":10")
		}

	} else {
		io.WriteString(resp, "-1")
	}
}

func CommentGetHistory(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if Post.Get("userID") != "" {
		page := 0
		core.TryInt(&page, Post.Get("page"))
		mode := false
		if Post.Get("mode") != "0" {
			mode = true
		}
		db := &core.MySQLConn{}
		defer db.CloseDB()
		if logger.Should(db.ConnectBlob(config)) != nil {
			return
		}
		acc := core.CAccount{DB: db}
		core.TryInt(&acc.Uid, Post.Get("userID"))
		if !acc.Exists(acc.Uid) {
			io.WriteString(resp, "-1")
			return
		}
		acc.LoadAuth(core.CAUTH_UID)
		acc.LoadStats()
		acc.LoadVessels()
		role := acc.GetRoleObj(false)
		cc := core.CComment{DB: db}
		comments := cc.GetAllCommentsHistory(acc.Uid, page, mode)
		if len(comments) == 0 {
			io.WriteString(resp, "#0:0:0")
		} else {
			output := ""
			for _, comm := range comments {
				output += connectors.GetCommentHistory(comm, acc, role)
			}
			io.WriteString(resp, output[:len(output)-1]+"#"+strconv.Itoa(cc.CountCommentHistory(acc.Uid))+":"+strconv.Itoa(page*10)+":10")
		}

	} else {
		io.WriteString(resp, "-1")
	}
}

func CommentUpload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if core.CheckGDAuth(Post) && Post.Get("comment") != "" && Post.Get("levelID") != "" {
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
		comment := core.ClearGDRequest(Post.Get("comment"))
		percent := 0
		core.TryInt(&percent, Post.Get("percent"))
		percent %= 101
		cl := core.CLevel{DB: db}
		core.TryInt(&cl.Id, Post.Get("levelID"))
		if !cl.Exists(cl.Id) {
			io.WriteString(resp, "-1")
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
					io.WriteString(resp, "-1")
				case "ok":
					io.WriteString(resp, "1")
				default:
					io.WriteString(resp, "temp_1_"+cmdRes)
				}
			} else {
				cc := core.CComment{DB: db, Uid: xacc.Uid, LvlId: cl.Id, Comment: comment, Percent: percent}
				protect := core.CProtect{DB: db, DisableProtection: config.SecurityConfig.DisableProtection}
				if !core.OnComment(db, conf, config) {
					io.WriteString(resp, "-1")
					return
				}
				if protect.DetectComments(xacc.Uid) && cc.PostLevelComment() {
					io.WriteString(resp, "1")
				} else {
					io.WriteString(resp, "-1")
				}
			}

		} else {
			cc := core.CComment{DB: db, Uid: xacc.Uid, LvlId: cl.Id, Comment: comment, Percent: percent}
			protect := core.CProtect{DB: db, DisableProtection: config.SecurityConfig.DisableProtection}
			if !core.OnComment(db, conf, config) {
				io.WriteString(resp, "-1")
				return
			}
			if protect.DetectComments(xacc.Uid) && cc.PostLevelComment() {
				io.WriteString(resp, "1")
			} else {
				io.WriteString(resp, "-1")
			}
		}
	} else {
		io.WriteString(resp, "-1")
	}
}
