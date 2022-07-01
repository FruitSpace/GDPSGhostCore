package api

import (
	"io"
	"net/http"
	gorilla "github.com/gorilla/mux"
)

func GetChallenges(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func GetRewards(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}