package api

import (
	"HalogenGhostCore/core"
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
	if Post.Get("userName") != "" && Post.Get("password") != "" && Post.Get("saveData") != "" {
		uname := core.ClearGDRequest(Post.Get("userName"))
		pass := core.ClearGDRequest(Post.Get("password"))
		saveData := core.ClearGDRequest(Post.Get("saveData"))
		db := &core.MySQLConn{}
		defer db.CloseDB()
		if logger.Should(db.ConnectBlob(config)) != nil {
			return
		}
		acc := core.CAccount{DB: db}
		if acc.LogIn(uname, pass, IPAddr, 0) > 0 {
			savepath := "/gdps_savedata/" + vars["gdps"] + "/"
			taes := core.ThunderAES{}
			if logger.Should(taes.GenKey(config.ServerConfig.SrvKey)) != nil {
				return
			}
			if logger.Should(taes.Init()) != nil {
				return
			}
			datax, err := taes.EncryptRaw(saveData)
			if logger.Should(err) != nil {
				return
			}

			s3 := core.NewS3FS()
			if logger.Should(s3.PutFile(savepath+strconv.Itoa(acc.Uid)+".hsv", datax)) != nil {
				return
			}

			saveData = strings.ReplaceAll(strings.ReplaceAll(strings.Split(saveData, ";")[0], "_", "/"), "-", "+")
			b, err := base64.StdEncoding.DecodeString(saveData)
			if logger.Should(err) != nil {
				return
			}
			r, err := gzip.NewReader(bytes.NewBuffer(b))
			if logger.Should(err) != nil {
				return
			}
			d, err := io.ReadAll(r)
			if logger.Should(err) != nil {
				return
			}
			saveData = string(d)
			acc.LoadStats()
			acc.Orbs, _ = strconv.Atoi(strings.Split(strings.Split(saveData, "</s><k>14</k><s>")[1], "</s>")[0])
			acc.LvlsCompleted, _ = strconv.Atoi(strings.Split(strings.Split(strings.Split(saveData, "<k>GS_value</k>")[1], "</s><k>4</k><s>")[1], "</s>")[0])
			acc.PushStats()
			//! Temp
			s3.DeleteFile("/savedata_old/" + vars["gdps"] + "/files/savedata/" + strconv.Itoa(acc.Uid) + ".hal")
			io.WriteString(resp, "1")
		} else {
			io.WriteString(resp, "-2")
		}
	} else {
		io.WriteString(resp, "-1")
	}
}

func AccountSync(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if Post.Get("userName") != "" && Post.Get("password") != "" {
		uname := core.ClearGDRequest(Post.Get("userName"))
		pass := core.ClearGDRequest(Post.Get("password"))
		db := &core.MySQLConn{}
		defer db.CloseDB()
		if logger.Should(db.ConnectBlob(config)) != nil {
			return
		}
		acc := core.CAccount{DB: db}
		if acc.LogIn(uname, pass, IPAddr, 0) > 0 {
			savepath := "/gdps_savedata/" + vars["gdps"] + "/" + strconv.Itoa(acc.Uid) + ".hsv"
			s3 := core.NewS3FS()
			if d, err := s3.GetFile(savepath); err == nil {
				taes := core.ThunderAES{}
				if logger.Should(taes.GenKey(config.ServerConfig.SrvKey)) != nil {
					return
				}
				if logger.Should(taes.Init()) != nil {
					return
				}
				data, err := taes.DecryptRaw(d)
				if err != nil {
					core.ReportFail(fmt.Sprintf("[%s] NG savedata decrypt error for `%s`", vars["gdps"], uname))
					return
				}
				io.WriteString(resp, data+";21;30;a;a")
				//! Temp transitional
			} else if d, err := s3.GetFile("/savedata_old/" + vars["gdps"] + "/files/savedata/" + strconv.Itoa(acc.Uid) + ".hal"); err == nil {
				taes := core.ThunderAES{}
				if logger.Should(taes.GenKey(pass)) != nil {
					return
				}
				if logger.Should(taes.Init()) != nil {
					return
				}
				data, err := taes.DecryptLegacy(string(d))
				if err != nil {
					core.ReportFail(fmt.Sprintf("[%s] HAL savedata decrypt error for `%s`", vars["gdps"], uname))
					return
				}
				io.WriteString(resp, data+";21;30;a;a")
			} else {
				io.WriteString(resp, "-1")
			}
		} else {
			io.WriteString(resp, "-2")
		}
	} else {
		io.WriteString(resp, "-1")
	}
}

func AccountManagement(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	vars := gorilla.Vars(req)
	http.Redirect(resp, req, "https://gofruit.space/gdps/"+vars["gdps"], http.StatusMovedPermanently)
}

func AccountLogin(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
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
	if Post.Get("userName") != "" && Post.Get("password") != "" {
		uname := core.ClearGDRequest(Post.Get("userName"))
		pass := core.ClearGDRequest(Post.Get("password"))
		db := &core.MySQLConn{}
		defer db.CloseDB()
		if logger.Should(db.ConnectBlob(config)) != nil {
			return
		}
		acc := core.CAccount{DB: db}
		uid := acc.LogIn(uname, pass, IPAddr, 0)
		if uid < 0 {
			io.WriteString(resp, strconv.Itoa(uid))
		} else {
			io.WriteString(resp, strconv.Itoa(uid)+","+strconv.Itoa(uid))
			core.RegisterAction(core.ACTION_USER_LOGIN, 0, uid, map[string]string{"uname": uname}, db)
		}
	} else {
		io.WriteString(resp, "-1")
	}
}

func AccountRegister(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := req.Header.Get("CF-Connecting-IP")
	if IPAddr == "" {
		IPAddr = req.Header.Get("X-Real-IP")
	}
	if IPAddr == "" {
		IPAddr = strings.Split(req.RemoteAddr, ":")[0]
	}
	vars := gorilla.Vars(req)
	if conf.MaintenanceMode {
		resp.WriteHeader(403)
		core.SendMessageDiscord(fmt.Sprintf("[%s] %s reached registration killswitch", vars["gdps"], IPAddr))
		return
	}
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
	if Post.Get("userName") != "" && Post.Get("password") != "" && Post.Get("email") != "" {
		uname := core.ClearGDRequest(Post.Get("userName"))
		pass := core.ClearGDRequest(Post.Get("password"))
		email := core.ClearGDRequest(Post.Get("email"))
		db := &core.MySQLConn{}
		defer db.CloseDB()
		if logger.Should(db.ConnectBlob(config)) != nil {
			return
		}
		acc := core.CAccount{DB: db}
		if core.OnRegister(db, conf, config) {
			uid := acc.Register(uname, pass, email, IPAddr, config.SecurityConfig.AutoActivate)
			io.WriteString(resp, strconv.Itoa(uid))
			if uid > 0 {
				core.RegisterAction(core.ACTION_USER_REGISTER, 0, uid, map[string]string{"uname": uname, "email": email}, db)
			}
		} else {
			io.WriteString(resp, "-1")
		}
	} else {
		io.WriteString(resp, "-1")
	}
}
