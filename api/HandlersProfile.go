package api

import (
	"io"
	"net/http"
	gorilla "github.com/gorilla/mux"
)

func GetUserInfo(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func GetUserList(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func GetUsers(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func UpdateAccountSettings(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}