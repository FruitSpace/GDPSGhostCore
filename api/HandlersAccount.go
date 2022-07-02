package api

import (
	"HalogenGhostCore/core"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"github.com/go-redis/redis/v8"
	gorilla "github.com/gorilla/mux"
	"io"
	"log"
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
	config,err:=conf.LoadById(vars["gdps"])
	if err!=nil{
		if err==redis.Nil {return}
		io.WriteString(resp,"There was an error")
		log.Panicln(err.Error())
	}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Has("userName") && Post.Has("password") && Post.Get("userName")!="" && Post.Get("password")!="" {
		uname:=core.ClearGDRequest(Post.Get("userName"))
		pass:=core.ClearGDRequest(Post.Get("password"))
		saveData:=core.ClearGDRequest(Post.Get("saveData"))
		db:=core.MySQLConn{}
		if err:=db.ConnectBlob(config); err!=nil {log.Fatalln(err.Error())}
		acc:=core.CAccount{DB: db}
		if acc.LogIn(uname,pass, IPAddr, 0)>0 {
			savepath:=conf.SavePath+"/"+vars["gdps"]+"/savedata/"
			taes:=core.ThunderAES{}
			logger:=core.Logger{Output: os.Stderr}
			logger.Must(taes.GenKey(config.ServerConfig.SrvKey))
			logger.Must(taes.Init())
			datax,err:=taes.EncryptRaw(saveData)
			if err!=nil{
				io.WriteString(resp,"There was an error")
				logger.LogErr(taes,err.Error())
				return
			}
			os.MkdirAll(savepath,os.ModePerm)
			logger.Must(os.WriteFile(savepath+strconv.Itoa(acc.Uid)+".hal",datax,0644))
			saveData=strings.ReplaceAll(strings.ReplaceAll(strings.Split(saveData,";")[0],"_","/"),"-","+")
			b,err:=base64.StdEncoding.DecodeString(saveData)
			logger.Must(err)
			r,err:=gzip.NewReader(bytes.NewBuffer(b))
			logger.Must(err)
			d,err:=io.ReadAll(r)
			logger.Must(err)
			saveData=string(d)
			acc.LoadStats()
			acc.Orbs,_=strconv.Atoi(strings.Split(strings.Split(saveData,"</s><k>14</k><s>")[1],"</s>")[0])
			acc.LvlsCompleted,_=strconv.Atoi(strings.Split(strings.Split(strings.Split(saveData,"<k>GS_value</k>")[1],"</s><k>4</k><s>")[1],"</s>")[0])
			acc.PushStats()
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
	config,err:=conf.LoadById(vars["gdps"])
	if err!=nil{
		if err==redis.Nil {return}
		io.WriteString(resp,"There was an error")
		log.Panicln(err.Error())
	}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Has("userName") && Post.Has("password") && Post.Get("userName")!="" && Post.Get("password")!="" {
		uname:=core.ClearGDRequest(Post.Get("userName"))
		pass:=core.ClearGDRequest(Post.Get("password"))
		db:=core.MySQLConn{}
		if err:=db.ConnectBlob(config); err!=nil {log.Fatalln(err.Error())}
		acc:=core.CAccount{DB: db}
		if acc.LogIn(uname,pass, IPAddr, 0)>0 {
			savepath:=conf.SavePath+"/"+vars["gdps"]+"/savedata/"+strconv.Itoa(acc.Uid)+".hal"
			if _, err := os.Stat(savepath); err==nil {
				logger:=core.Logger{Output: os.Stderr}
				taes := core.ThunderAES{}
				logger.Must(taes.GenKey(config.ServerConfig.SrvKey))
				logger.Must(taes.Init())
				d,err:=os.ReadFile(savepath)
				data,err:=taes.DecryptRaw(d)
				if err!=nil{
					io.WriteString(resp,"There was an error")
					log.Panicln(err.Error())
					return
				}
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
    io.WriteString(resp,vars["gdps"])
	http.Redirect(resp,req,"https://get.halhost.cc/"+vars["gdps"],http.StatusMovedPermanently)
}

func AccountLogin(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	IPAddr:=req.Header.Get("CF-Connecting-IP")
	if IPAddr=="" {IPAddr=req.Header.Get("X-Real-IP")}
	if IPAddr=="" {IPAddr=strings.Split(req.RemoteAddr,":")[0]}
	vars:= gorilla.Vars(req)
	config,err:=conf.LoadById(vars["gdps"])
	if err!=nil{
		if err==redis.Nil {return}
		io.WriteString(resp,"There was an error")
		log.Panicln(err.Error())
	}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Has("userName") && Post.Has("password") && Post.Get("userName")!="" && Post.Get("password")!="" {
		uname:=core.ClearGDRequest(Post.Get("userName"))
		pass:=core.ClearGDRequest(Post.Get("password"))
		db:=core.MySQLConn{}
		if err:=db.ConnectBlob(config); err!=nil {log.Fatalln(err.Error())}
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
	config, err := conf.LoadById(vars["gdps"])
	if err != nil {
		if err == redis.Nil {return}
		io.WriteString(resp, "There was an error")
		log.Panicln(err.Error())
	}
	//Get:=req.URL.Query()
	Post:=ReadPost(req)
	if Post.Has("userName") && Post.Has("password") && Post.Has("email") &&
		Post.Get("userName") != "" && Post.Get("password") != "" && Post.Get("email")!="" {
		uname := core.ClearGDRequest(Post.Get("userName"))
		pass := core.ClearGDRequest(Post.Get("password"))
		email := core.ClearGDRequest(Post.Get("email"))
		db := core.MySQLConn{}
		if err := db.ConnectBlob(config); err != nil {
			log.Fatalln(err.Error())
		}
		acc := core.CAccount{DB: db}
		uid:=acc.Register(uname,pass,email,IPAddr)
		io.WriteString(resp,strconv.Itoa(uid))
		if uid>0 {
			core.RegisterAction(core.ACTION_USER_REGISTER,0,uid, map[string]string{"uname":uname,"email":email},db)
		}else {
			io.WriteString(resp, "-1")
		}
	}else{
		io.WriteString(resp,"-1")
	}
}