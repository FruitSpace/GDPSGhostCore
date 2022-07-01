package api

import (
	"HalogenGhostCore/core"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
)

func ApiIntegra(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig){
	vars:= gorilla.Vars(req)
	io.WriteString(resp,"[GhostCore] Serving //"+vars["gdps"]+"//")
}
