package api

import (
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func AccountBackup(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if Post.Get("saveData") != "" {
		uname := core.ClearGDRequest(Post.Get("userName"))
		pass := core.ClearGDRequest(Post.Get("password"))
		saveData := core.ClearGDRequest(Post.Get("saveData"))
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		acc := core.CAccount{DB: db}
		var res int
		if Post.Get("gameVersion") == "22" {
			res = core.ToInt(acc.PerformGJPAuth(Post, IPAddr))
		} else {
			res = acc.LogIn(uname, pass, IPAddr, 0)
		}
		if res > 0 {
			savepath := "/gdps_savedata/" + vars["gdps"] + "/"
			taes := core.ThunderAES{}
			if logger.Should(taes.GenKey(config.ServerConfig.SrvKey)) != nil {
				serverError(connector)
				return
			}
			if logger.Should(taes.Init()) != nil {
				serverError(connector)
				return
			}
			datax, err := taes.EncryptRaw(saveData)
			if logger.Should(err) != nil {
				serverError(connector)
				return
			}

			s3 := core.NewS3FS()
			if logger.Should(s3.PutFile(savepath+strconv.Itoa(acc.Uid)+".hsv", datax)) != nil {
				serverError(connector)
				return
			}

			saveData = strings.ReplaceAll(strings.ReplaceAll(strings.Split(saveData, ";")[0], "_", "/"), "-", "+")
			b, err := base64.StdEncoding.DecodeString(saveData)
			if logger.Should(err) != nil {
				serverError(connector)
				return
			}
			r, err := gzip.NewReader(bytes.NewBuffer(b))
			if logger.Should(err) != nil {
				serverError(connector)
				return
			}
			d, err := io.ReadAll(r)
			if logger.Should(err) != nil {
				serverError(connector)
				return
			}
			saveData = string(d)
			acc.LoadStats()
			acc.Orbs, _ = strconv.Atoi(strings.Split(strings.Split(saveData, "</s><k>14</k><s>")[1], "</s>")[0])
			acc.LvlsCompleted, _ = strconv.Atoi(strings.Split(strings.Split(strings.Split(saveData, "<k>GS_value</k>")[1], "</s><k>4</k><s>")[1], "</s>")[0])
			acc.PushStatsAndExtra()
			//! Temp
			s3.DeleteFile("/savedata_old/" + vars["gdps"] + "/files/savedata/" + strconv.Itoa(acc.Uid) + ".hal")
			connector.Success("Backup successful")
		} else {
			connector.Error("-2", "Invalid credentials")
		}
	} else {
		connector.Error("-1", "Bad request")
	}
}

func AccountSync(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	connector := connectors.NewConnector(req.URL.Query().Has("json"))
	defer func() { _, _ = io.WriteString(resp, connector.Output()) }()
	config, err := conf.LoadById(vars["gdps"])
	if logger.Should(err) != nil {
		serverError(connector)
		return
	}
	if core.CheckIPBan(IPAddr, config) {
		connector.Error("-1", "Banned")
		return
	}
	Post := ReadPost(req)
	if (Post.Get("userName") != "" && Post.Get("password") != "") || Post.Get("gjp2") != "" {
		uname := core.ClearGDRequest(Post.Get("userName"))
		pass := core.ClearGDRequest(Post.Get("password"))
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			return
		}
		acc := core.CAccount{DB: db}
		var res int
		if Post.Get("gameVersion") == "22" {
			res = core.ToInt(acc.PerformGJPAuth(Post, IPAddr))
		} else {
			res = acc.LogIn(uname, pass, IPAddr, 0)
		}
		if res > 0 {
			savepath := "/gdps_savedata/" + vars["gdps"] + "/" + strconv.Itoa(acc.Uid) + ".hsv"
			s3 := core.NewS3FS()
			if d, err := s3.GetFile(savepath); err == nil {
				taes := core.ThunderAES{}
				if logger.Should(taes.GenKey(config.ServerConfig.SrvKey)) != nil {
					serverError(connector)
					return
				}
				if logger.Should(taes.Init()) != nil {
					serverError(connector)
					return
				}
				data, err := taes.DecryptRaw(d)
				if err != nil {
					serverError(connector)
					core.ReportFail(fmt.Sprintf("[%s] NG savedata decrypt error for `%s`", vars["gdps"], uname))
					return
				}
				connector.Account_Sync(data)
				//! Temp transitional
			} else if d, err := s3.GetFile("/savedata_old/" + vars["gdps"] + "/files/savedata/" + strconv.Itoa(acc.Uid) + ".hal"); err == nil {
				taes := core.ThunderAES{}
				if logger.Should(taes.GenKey(pass)) != nil {
					serverError(connector)
					return
				}
				if logger.Should(taes.Init()) != nil {
					serverError(connector)
					return
				}
				data, err := taes.DecryptLegacy(string(d))
				if err != nil {
					serverError(connector)
					core.ReportFail(fmt.Sprintf("[%s] HAL savedata decrypt error for `%s`", vars["gdps"], uname))
					return
				}
				connector.Account_Sync(data)
			} else {
				connector.Error("-1", "No savedata found")
			}
		} else {
			connector.Error("-2", "Invalid credentials")
		}
	} else {
		connector.Error("-1", "Bad request")
	}
}

func AccountManagement(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	vars := gorilla.Vars(req)
	http.Redirect(resp, req, "https://gofruit.space/gdps/"+vars["gdps"], http.StatusMovedPermanently)
}

func AccountLogin(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if Post.Get("userName") != "" && Post.Get("password")+Post.Get("gjp2") != "" {
		gjp2 := core.ClearGDRequest(Post.Get("gjp2"))
		uname := core.ClearGDRequest(Post.Get("userName"))
		pass := core.ClearGDRequest(Post.Get("password"))
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		acc := core.CAccount{DB: db}
		var uid int
		if len(gjp2) != 0 {
			uid = acc.LogIn22(uname, gjp2, IPAddr, 0)
		} else {
			uid = acc.LogIn(uname, pass, IPAddr, 0)
		}

		if uid < 0 {
			connector.Error("-1", "Invalid credentials")
		} else {
			connector.Account_Login(uid)
			core.RegisterAction(core.ACTION_USER_LOGIN, 0, uid, map[string]string{"uname": uname}, db)
		}
	} else {
		connector.Error("-1", "Bad request")
	}
}

func AccountRegister(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := ipOf(req)
	vars := gorilla.Vars(req)
	if conf.MaintenanceMode {
		resp.WriteHeader(403)
		core.SendMessageDiscord(fmt.Sprintf("[%s] %s reached registration killswitch", vars["gdps"], IPAddr))
		return
	}

	//Ballistics
	if PrepareBallistics(req) {
		return
	}

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
	if Post.Get("userName") != "" && Post.Get("password") != "" && Post.Get("email") != "" {
		uname := core.ClearGDRequest(Post.Get("userName"))
		pass := core.ClearGDRequest(Post.Get("password"))
		email := core.ClearGDRequest(Post.Get("email"))
		db := &core.MySQLConn{}

		if logger.Should(db.ConnectBlob(config)) != nil {
			serverError(connector)
			return
		}
		acc := core.CAccount{DB: db}
		if core.OnRegister(db, conf, config) {
			uid := acc.Register(uname, pass, email, IPAddr, config.SecurityConfig.AutoActivate)
			if uid > 0 {
				core.RegisterAction(core.ACTION_USER_REGISTER, 0, uid, map[string]string{"uname": uname, "email": email}, db)
				connector.Success("Registered")
			} else {
				connector.Error(strconv.Itoa(uid), "Refer to https://github.com/gd-programming/gd.docs/blob/main/docs/topics/status_codes.md#registergjaccount")
			}
		} else {
			connector.Error("-1", "Player limits exceeded")
		}
	} else {
		connector.Error("-1", "Bad request")
	}
}
