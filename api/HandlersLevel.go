package api

import (
	"io"
	"net/http"
	gorilla "github.com/gorilla/mux"
)

func GetGauntlets(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func GetMapPacks(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}



func LevelDelete(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func LevelDownload(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}


func LevelGetDaily(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func LevelGetLevels(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func LevelReport(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func LevelUpdateDescription(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func LevelUpload(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func RateDemon(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func RateStar(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}

func SuggestStars(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
    io.WriteString(resp,vars["gdps"])
}