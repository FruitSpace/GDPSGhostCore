package api

import (
	"HalogenGhostCore/core"
	gorilla "github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"net/url"
)

func AccountBackup(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
	config,err:=conf.LoadById(vars["gdps"])
	if err!=nil{
		io.WriteString(resp,"There was an error")
		log.Println(err.Error())
	}
	Get:=req.URL.Query()
	Post,_:=url.ParseQuery(ReadPost(req))
	if Post.Has("userName")
}

func AccountSync(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func AccountManagement(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
	http.Redirect(resp,req,"https://get.halhost.cc/"+vars["gdps"],http.StatusMovedPermanently)
}

func AccountLogin(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func AccountRegister(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}