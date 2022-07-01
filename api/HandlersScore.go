package api

import (
	"io"
	"net/http"
	gorilla "github.com/gorilla/mux"
)

func GetCreators(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func GetLevelScores(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func GetScores(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func UpdateUserScore(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}