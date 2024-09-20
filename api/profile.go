package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
)

func GetUserInfo(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if Post.Get("targetAccountID") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		acc := core.CAccount{DB: db}
		var uidSelf int
		core.TryInt(&acc.Uid, Post.Get("targetAccountID"))
		if core.CheckGDAuth(Post) {
			xacc := core.CAccount{DB: db}
			if !xacc.PerformGJPAuth(Post, IPAddr) {
				uidSelf = 0
			} else {
				uidSelf = xacc.Uid
			}
		}
		if !acc.Exists(acc.Uid) {
			connector.Error("-1", "User not Found")
			return
		}
		acc.LoadAll()
		blacklist := strings.Split(acc.Blacklist, ",")
		if uidSelf > 0 && slices.Contains(blacklist, strconv.Itoa(uidSelf)) {
			connector.Error("-1", "User has blacklisted you")
			return
		}

		connector.Profile_GetUserProfile(acc, uidSelf)
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func GetUserList(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
		var cType int
		core.TryInt(&cType, Post.Get("type"))
		xacc := core.CAccount{DB: db}
		if !xacc.PerformGJPAuth(Post, IPAddr) {
			connector.Error("-1", "Invalid Credentials")
			return
		}
		xacc.LoadSocial()
		var usersToDump []core.CAccount
		switch cType {
		case 0:
			if xacc.FriendsCount == 0 {
				connector.Error("-2", "No Friends")
				return
			} else {
				flist := strings.Split(xacc.FriendshipIds, ",")
				for _, fid := range flist {
					var Xfid int
					core.TryInt(&Xfid, fid)
					cf := core.CFriendship{DB: db}
					uid1, uid2 := cf.GetFriendByFID(Xfid)
					acc := core.CAccount{DB: db, Uid: uid1}
					if uid1 == xacc.Uid {
						acc.Uid = uid2
					}
					acc.LoadAuth(core.CAUTH_UID)
					acc.LoadVessels()
					acc.LoadStats()
					usersToDump = append(usersToDump, acc)
				}
			}
		case 1:
			blacklist := strings.Split(xacc.Blacklist, ",")
			if len(xacc.Blacklist) == 0 || len(blacklist) == 0 {
				connector.Error("-2", "No blacklisted users")
				return
			} else {
				for _, buid := range blacklist {
					acc := core.CAccount{DB: db}
					core.TryInt(&acc.Uid, buid)
					if acc.Uid == 0 {
						continue
					}
					acc.LoadAuth(core.CAUTH_UID)
					acc.LoadVessels()
					acc.LoadStats()
					usersToDump = append(usersToDump, acc)
				}
			}
		default:
			connector.Error("-1", "Bad Request")
			return
		}

		connector.Profile_ListUserProfiles(usersToDump)
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func GetUsers(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if Post.Get("str") != "" {
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		xacc := core.CAccount{DB: db}
		uids := xacc.SearchUsers(core.ClearGDRequest(Post.Get("str")))
		if len(uids) == 0 {
			connector.Error("-1", "No users found")
		} else {
			var accs []core.CAccount
			for _, uid := range uids {
				accs = append(accs, core.CAccount{DB: db, Uid: uid})
			}
			connector.Profile_ListUserProfiles(accs)
		}
	} else {
		connector.Error("-1", "Bad Request")
	}
}

func UpdateAccountSettings(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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

		core.TryInt(&xacc.MS, Post.Get("mS"))
		core.TryInt(&xacc.FrS, Post.Get("frS"))
		core.TryInt(&xacc.CS, Post.Get("cS"))
		xacc.Youtube = core.ClearGDRequest(Post.Get("yt"))
		xacc.Twitter = core.ClearGDRequest(Post.Get("twitter"))
		xacc.Twitch = core.ClearGDRequest(Post.Get("twitch"))
		xacc.PushSettings()
		connector.Success("Account updated")
	} else {
		connector.Error("-1", "Bad Request")
	}
}
