package api

import (
	"io"
	"net/http"
	gorilla "github.com/gorilla/mux"
)

func GetChallenges(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func GetRewards(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}