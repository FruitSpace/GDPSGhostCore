package api

import (
	"io"
	"net/http"
	gorilla "github.com/gorilla/mux"
)

func GetAccountUrl(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func GetSongInfo(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func GetTopArtists(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func LikeItem(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func RequestMod(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}