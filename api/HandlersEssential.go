package api

import (
	"HalogenGhostCore/core"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
)

func GetAccountUrl(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,"http://s.halhost.cc/"+vars["gdps"]+"/database")
}

func GetSongInfo(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func GetTopArtists(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func LikeItem(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func RequestMod(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}