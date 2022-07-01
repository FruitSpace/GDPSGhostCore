package api

import (
	"io"
	"net/http"
	gorilla "github.com/gorilla/mux"
)

func GetCreators(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func GetLevelScores(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func GetScores(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func UpdateUserScore(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}