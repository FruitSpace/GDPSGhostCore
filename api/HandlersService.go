package api

import (
	"HalogenGhostCore/core"
	"encoding/json"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strings"
)

var NotFoundTemplate = `
<html>
	<head>
		<title>404 | Not Found</title>
		<style>
			body{background-color:#212125;color:white;text-align:center;}
			h3{margin:5rem 0;}
		</style>
	</head>
	<body>
		<h3>You asked for [PATH], but found another 404 page</h3>
	</body>
</html>`

func Shield(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
	io.WriteString(resp,"[GhostCore] Serving //"+vars["gdps"]+"//")
}

func Redirector(resp http.ResponseWriter, req *http.Request){
	http.Redirect(resp,req,"https://halhost.cc",http.StatusMovedPermanently)
}

type NotFoundHandler int

func (n NotFoundHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	io.WriteString(resp,strings.ReplaceAll(NotFoundTemplate,"[PATH]",req.URL.Path))
}

// Private API

func ModifyGDPS(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	//vars:= gorilla.Vars(req)
	//Post:=ReadPost(req)
	//response:=map[string]string{"status":"ok"}
	//if Post.Get("key")!=conf.MasterKey {
	//	response["status"]="error"
	//	response["error"]="Unauthenticated"
	//	SendJson(resp, response)
	//	return
	//}
	//logger:=core.Logger{Output: os.Stderr}
	//config,err:=conf.LoadById(vars["gdps"])
	//if logger.Should(err)!=nil {return}
	//switch req.Method {
	//case "GET":
	//case "POST":
	//
	//}
}

func CreateGDPS(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	vars:= gorilla.Vars(req)
	Post:=ReadPost(req)
	response:=map[string]string{"status":"ok"}
	if Post.Get("key")!=conf.MasterKey {
		response["status"]="error"
		response["error"]="Unauthenticated"
		SendJson(resp, response)
		return
	}
	logger:=core.Logger{Output: os.Stderr}
	config,err:=conf.LoadById(vars["gdps"])
	if logger.Should(err)!=nil {return}
	core.HalInitialize(config, conf)
}

func EventAction(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	//vars:= gorilla.Vars(req)
	Post:=ReadPost(req)
	response:=map[string]string{"status":"ok"}
	if Post.Get("key")!=conf.MasterKey {
		response["status"]="error"
		response["error"]="Unauthenticated"
		SendJson(resp, response)
		return
	}
	switch req.Method {

	}
}

func SendJson(resp http.ResponseWriter, jsonData map[string]string){
	data,_:=json.Marshal(jsonData)
	io.WriteString(resp,string(data))
}