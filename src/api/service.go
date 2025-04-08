package api

import (
	"HalogenGhostCore/core"
	"encoding/json"
	gorilla "github.com/gorilla/mux"
	"html"
	"io"
	"net/http"
	"strings"
)

var NotFoundTemplate = `
<html>
	<head>
		<title>404 | Not Found</title>
		<style>
			body{background-color:#212125;color:white;text-align:center;}
			h3{margin:5rem 0;}
		</style>
	</head>
	<body>
		<h3>You asked for [PATH], but found another 404 page</h3>
		<p>The chances are there is something, but it's not exactly here</p>
	</body>
</html>`

func Shield(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	vars := gorilla.Vars(req)
	io.WriteString(resp, "[GhostCore] Serving //"+vars["gdps"]+"//")
}

func Redirector(resp http.ResponseWriter, req *http.Request) {
	http.Redirect(resp, req, "https://fruitspace.ru/", http.StatusMovedPermanently)
}

type NotFoundHandler int

func (n NotFoundHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	io.WriteString(resp, strings.ReplaceAll(NotFoundTemplate, "[PATH]", html.EscapeString(req.URL.Path)))
}

// Private API

func TriggerMaintenance(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	if req.URL.Query().Get("key") != conf.MasterKey {
		return
	}
	io.WriteString(resp, "OK")
	//go core.MaintainTasks()
}

func EventAction(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	q := req.URL.Query()
	if q.Get("key") != conf.MasterKey {
		io.WriteString(resp, "KEY")
		return
	}
	mk := "off"
	if conf.MaintenanceMode {
		mk = "on"
	}
	switch q.Get("a") {
	case "on":
		conf.MaintenanceMode = true
	case "off":
		conf.MaintenanceMode = false
	default:
		io.WriteString(resp, mk)
	}
	core.SendMessageDiscord("Touched killswitch: status: " + mk)
}

func SendJson(resp http.ResponseWriter, jsonData map[string]string) {
	data, _ := json.Marshal(jsonData)
	io.WriteString(resp, string(data))
}
