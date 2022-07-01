package api

import (
	"io"
	"net/http"
	gorilla "github.com/gorilla/mux"
)

func BlockUser(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func UnblockUser(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}


func FriendAcceptRequest(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func FriendRejectRequest(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func FriendGetRequests(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func FriendReadRequest(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func FriendRemove(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func FriendRequest(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}


func MessageDelete(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func MessageGet(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func MessageGetAll(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}

func MessageUpload(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
}