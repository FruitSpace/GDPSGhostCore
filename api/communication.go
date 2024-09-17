package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
)

func BlockUser(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	se := func() {
		connector.Error("-1", "Server Error")
	}
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
	if core.CheckGDAuth(Post) && Post.Get("targetAccountID") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			se()
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		var uidTarget int
		core.TryInt(&uidTarget, Post.Get("targetAccountID"))
		if uidTarget > 0 {
			xacc.UpdateBlacklist(core.CBLACKLIST_BLOCK, uidTarget)
		}
		connector.Success("Blocked user successfully")
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func UnblockUser(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	se := func() {
		connector.Error("-1", "Server Error")
	}
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
	if core.CheckGDAuth(Post) && Post.Get("targetAccountID") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			se()
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		var uidTarget int
		core.TryInt(&uidTarget, Post.Get("targetAccountID"))
		if uidTarget > 0 {
			xacc.UpdateBlacklist(core.CBLACKLIST_UNBLOCK, uidTarget)
		}
		connector.Success("Unblocked user successfully")
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func FriendAcceptRequest(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	se := func() {
		connector.Error("-1", "Server Error")
	}
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
	if core.CheckGDAuth(Post) && Post.Get("requestID") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			se()
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		var requestId int
		core.TryInt(&requestId, Post.Get("requestID"))
		if requestId > 0 {
			cf := core.CFriendship{DB: db}
			if stat := cf.AcceptFriendRequest(requestId, xacc.Uid); stat > 0 {
				connector.Success("Accepted friend request successfully")
			} else {
				connector.Error("-1", "Friend request not found")
			}
		} else {
			connector.Error("-1", "Bad Request")
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func FriendRejectRequest(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	se := func() {
		connector.Error("-1", "Server Error")
	}
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
	if core.CheckGDAuth(Post) && Post.Get("targetAccountID") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			se()
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		var targetId int
		core.TryInt(&targetId, Post.Get("targetAccountID"))
		if targetId > 0 {
			cf := core.CFriendship{DB: db}
			issender := Post.Get("isSender") == "1"
			cf.RejectFriendRequestByUid(xacc.Uid, targetId, issender)
			connector.Success("Rejected friend request successfully")
		} else {
			connector.Error("-1", "Bad Request")
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func FriendGetRequests(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	se := func() {
		connector.Error("-1", "Server Error")
	}
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
			se()
			return
		}
		xacc := core.CAccount{DB: db}
		page := 0
		core.TryInt(&page, Post.Get("page"))
		getSent := Post.Get("getSent") == "1"
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		cf := core.CFriendship{DB: db}
		count, frqs := cf.GetFriendRequests(xacc.Uid, page, getSent)
		if len(frqs) == 0 {
			connector.Error("-2", "No requests found")
		} else {
			connector.Communication_FriendGetRequests(frqs, count, page)
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func FriendReadRequest(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	se := func() {
		connector.Error("-1", "Server Error")
	}
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
	if core.CheckGDAuth(Post) && Post.Get("requestID") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			se()
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		var requestId int
		core.TryInt(&requestId, Post.Get("requestID"))
		if requestId > 0 {
			cf := core.CFriendship{DB: db}
			cf.ReadFriendRequest(requestId)
			connector.Success("Read friend request successfully")
		} else {
			connector.Error("-1", "Bad Request")
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func FriendRemove(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	se := func() {
		connector.Error("-1", "Server Error")
	}
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
	if core.CheckGDAuth(Post) && Post.Get("targetAccountID") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			se()
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		var targetId int
		core.TryInt(&targetId, Post.Get("targetAccountID"))
		if targetId > 0 {
			cf := core.CFriendship{DB: db}
			cf.DeleteFriendship(xacc.Uid, targetId)
			connector.Success("Removed friend successfully")
		} else {
			connector.Error("-1", "Bad Request")
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func FriendRequest(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	se := func() {
		connector.Error("-1", "Server Error")
	}
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
	if core.CheckGDAuth(Post) && Post.Get("toAccountID") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			se()
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		var targetId int
		core.TryInt(&targetId, Post.Get("toAccountID"))
		if targetId > 0 {
			cf := core.CFriendship{DB: db}
			comment := Post.Get("comment")
			comment = core.ClearGDRequest(comment)
			if stat := cf.RequestFriend(xacc.Uid, targetId, comment); stat > 0 {
				connector.Success("Request sent successfully")
			} else {
				connector.Error("-1", "Failed to send request")
			}
		} else {
			connector.Error("-1", "Bad Request")
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func MessageDelete(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	se := func() {
		connector.Error("-1", "Server Error")
	}
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
	if core.CheckGDAuth(Post) && Post.Get("messageID") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			se()
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		var msgId int
		core.TryInt(&msgId, Post.Get("messageID"))
		if msgId > 0 {
			cm := core.CMessage{DB: db, Id: msgId}
			cm.DeleteMessage(xacc.Uid)
		}
		connector.Success("Message deleted successfully")
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func MessageGet(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	se := func() {
		connector.Error("-1", "Server Error")
	}
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
	if core.CheckGDAuth(Post) && Post.Get("messageID") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			se()
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		var msgId int
		core.TryInt(&msgId, Post.Get("messageID"))
		cm := core.CMessage{DB: db}
		if cm.Exists(msgId) {
			cm.LoadMessageById(msgId)
			if xacc.Uid == cm.UidSrc || xacc.Uid == cm.UidDest {
				connector.Communication_MessageGet(cm, xacc.Uid)
			} else {
				connector.Error("-1", "Message not found")
			}
		} else {
			connector.Error("-1", "Message not found")
		}

	} else {
		connector.Error("-1", "Bad Request")
	}
}

func MessageGetAll(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	se := func() {
		connector.Error("-1", "Server Error")
	}
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
			se()
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		page := 0
		core.TryInt(&page, Post.Get("page"))
		getSent := Post.Get("getSent") == "1"
		cm := core.CMessage{DB: db}
		count, msgs := cm.GetMessageForUid(xacc.Uid, page, getSent)
		if len(msgs) == 0 {
			connector.Error("-2", "No messages found")
		} else {
			connector.Communication_MessageGetAll(msgs, getSent, count, page)
		}

	} else {
		connector.Error("-1", "Bad Request")
	}
}

func MessageUpload(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	se := func() {
		connector.Error("-1", "Server Error")
	}
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
	if core.CheckGDAuth(Post) && Post.Get("toAccountID") != "" && Post.Get("body") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			se()
			return
		}
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		var uidDest int
		core.TryInt(&uidDest, Post.Get("toAccountID"))
		body := core.ClearGDRequest(Post.Get("body"))
		subject := core.ClearGDRequest(Post.Get("subject"))
		cm := core.CMessage{
			DB:      db,
			UidSrc:  xacc.Uid,
			UidDest: uidDest,
			Subject: subject,
			Message: body,
		}
		protect := core.CProtect{DB: db, DisableProtection: config.SecurityConfig.DisableProtection}
		if protect.DetectMessages(xacc.Uid) && cm.SendMessageObj() {
			connector.Success("Message sent successfully")
		} else {
			connector.Error("-1", "Failed to send message")
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}
