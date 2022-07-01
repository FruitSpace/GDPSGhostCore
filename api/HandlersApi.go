package api

import (
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
)

func ApiIntegra(resp http.ResponseWriter, req *http.Request){
	vars:= gorilla.Vars(req)
	io.WriteString(resp,"[GhostCore] Serving //"+vars["gdps"]+"//")
}
