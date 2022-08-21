package api

import (
	"HalogenGhostCore/core"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func AccountBackup(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Get("userName")!="" && Post.Get("password")!="" && Post.Get("saveData")!="" {
		uname:=core.ClearGDRequest(Post.Get("userName"))
		pass:=core.ClearGDRequest(Post.Get("password"))
		saveData:=core.ClearGDRequest(Post.Get("saveData"))
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		acc:=core.CAccount{DB: db}
		if acc.LogIn(uname,pass, IPAddr, 0)>0 {
			savepath:=conf.SavePath+"/"+vars["gdps"]+"/savedata/"
			taes:=core.ThunderAES{}
			if logger.Should(taes.GenKey(config.ServerConfig.SrvKey))!=nil {return}
			if logger.Should(taes.Init())!=nil {return}
			datax,err:=taes.EncryptRaw(saveData)
			if logger.Should(err)!=nil {return}
			os.MkdirAll(savepath,os.ModePerm)
			if logger.Should(os.WriteFile(savepath+strconv.Itoa(acc.Uid)+".hsv",datax,0644))!=nil {return}
			saveData=strings.ReplaceAll(strings.ReplaceAll(strings.Split(saveData,";")[0],"_","/"),"-","+")
			b,err:=base64.StdEncoding.DecodeString(saveData)
			if logger.Should(err)!=nil {return}
			r,err:=gzip.NewReader(bytes.NewBuffer(b))
			if logger.Should(err)!=nil {return}
			d,err:=io.ReadAll(r)
			if logger.Should(err)!=nil {return}
			saveData=string(d)
			acc.LoadStats()
			acc.Orbs,_=strconv.Atoi(strings.Split(strings.Split(saveData,"</s><k>14</k><s>")[1],"</s>")[0])
			acc.LvlsCompleted,_=strconv.Atoi(strings.Split(strings.Split(strings.Split(saveData,"<k>GS_value</k>")[1],"</s><k>4</k><s>")[1],"</s>")[0])
			acc.PushStats()
			//! Temp
			os.Remove("/var/www/gdps/"+vars["gdps"]+"/files/savedata/"+strconv.Itoa(acc.Uid)+".hal")
			io.WriteString(resp,"1")
		}else{io.WriteString(resp,"-2")}
	}else{
		io.WriteString(resp,"-1")
	}
}

func AccountSync(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Get("userName")!="" && Post.Get("password")!="" {
		uname:=core.ClearGDRequest(Post.Get("userName"))
		pass:=core.ClearGDRequest(Post.Get("password"))
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		acc:=core.CAccount{DB: db}
		if acc.LogIn(uname,pass, IPAddr, 0)>0 {
			savepath:=conf.SavePath+"/"+vars["gdps"]+"/savedata/"+strconv.Itoa(acc.Uid)+".hsv"
			if _, err := os.Stat(savepath); err==nil {
				taes := core.ThunderAES{}
				if logger.Should(taes.GenKey(config.ServerConfig.SrvKey))!=nil {return}
				if logger.Should(taes.Init())!=nil {return}
				d,err:=os.ReadFile(savepath)
				data,err:=taes.DecryptRaw(d)
				if logger.Should(err)!=nil {return}
				io.WriteString(resp,data+";21;30;a;a")
				//! Temp transitional
			}else if  _, err := os.Stat("/var/www/gdps/"+vars["gdps"]+"/files/savedata/"+strconv.Itoa(acc.Uid)+".hal"); err==nil{
				taes := core.ThunderAES{}
				if logger.Should(taes.GenKey(pass))!=nil {return}
				if logger.Should(taes.Init())!=nil {return}
				d,err:=os.ReadFile(savepath)
				data,err:=taes.DecryptRaw(d)
				if logger.Should(err)!=nil {return}
				io.WriteString(resp,data+";21;30;a;a")
			}else{
				io.WriteString(resp,"-1")
			}
		}else{
			io.WriteString(resp,"-2")
		}
	}else{
		io.WriteString(resp,"-1")
	}
}

func AccountManagement(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
	http.Redirect(resp,req,"https://get.halhost.cc/"+vars["gdps"],http.StatusMovedPermanently)
}

func AccountLogin(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Get("userName")!="" && Post.Get("password")!="" {
		uname:=core.ClearGDRequest(Post.Get("userName"))
		pass:=core.ClearGDRequest(Post.Get("password"))
		db:=core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		acc:=core.CAccount{DB: db}
		uid:=acc.LogIn(uname,pass, IPAddr, 0)
		if uid<0 {
			io.WriteString(resp,strconv.Itoa(uid))
		}else{
			io.WriteString(resp,strconv.Itoa(uid)+","+strconv.Itoa(uid))
			core.RegisterAction(core.ACTION_USER_LOGIN,0,uid, map[string]string{"uname":uname},db)
		}
	}else{
		io.WriteString(resp,"-1")
	}
}

func AccountRegister(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	IPAddr := req.Header.Get("CF-Connecting-IP")
	if IPAddr == "" {IPAddr = req.Header.Get("X-Real-IP")}
	if IPAddr == "" {IPAddr = strings.Split(req.RemoteAddr, ":")[0]}
	vars := gorilla.Vars(req)
	logger:=core.Logger{Output: os.Stderr}
	config, err := conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Get("userName") != "" && Post.Get("password") != "" && Post.Get("email")!="" {
		uname := core.ClearGDRequest(Post.Get("userName"))
		pass := core.ClearGDRequest(Post.Get("password"))
		email := core.ClearGDRequest(Post.Get("email"))
		db := core.MySQLConn{}
		if logger.Should(db.ConnectBlob(config))!=nil {return}
		acc := core.CAccount{DB: db}
		uid:=acc.Register(uname,pass,email,IPAddr)
		io.WriteString(resp,strconv.Itoa(uid))
		if uid>0 {
			core.RegisterAction(core.ACTION_USER_REGISTER,0,uid, map[string]string{"uname":uname,"email":email},db)
		}
	}else{
		io.WriteString(resp,"-1")
	}
}