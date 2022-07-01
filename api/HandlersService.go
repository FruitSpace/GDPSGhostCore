package api

import (
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
)

func Shield(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
	io.WriteString(resp,"[GhostCore] Serving //"+vars["gdps"]+"//")
}

func Redirector(resp http.ResponseWriter, req *http.Request){
	http.Redirect(resp,req,"https://halhost.cc",http.StatusMovedPermanently)
}