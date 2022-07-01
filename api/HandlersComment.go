package api

import (
	"io"
	"net/http"
	gorilla "github.com/gorilla/mux"
)

func AccountCommentDelete(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func AccountCommentGet(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func AccountCommentUpload(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}


func CommentDelete(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func CommentGet(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func CommentGetHistory(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func CommentUpload(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}