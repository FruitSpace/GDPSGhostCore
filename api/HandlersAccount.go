package api

import (
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
)

func AccountBackup(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func AccountSync(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func AccountManagement(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
	http.Redirect(resp,req,"https://get.halhost.cc/"+vars["gdps"],http.StatusMovedPermanently)
}

func AccountLogin(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func AccountRegister(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}