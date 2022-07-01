package api

import (
	"io"
	"net/http"
	gorilla "github.com/gorilla/mux"
)

func GetAccountUrl(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func GetSongInfo(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func GetTopArtists(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func LikeItem(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func RequestMod(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}