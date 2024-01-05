package api

import (
	"HalogenGhostCore/core"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
)

func GetContentURL(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	vars := gorilla.Vars(req)
	io.WriteString(resp, "https://rugd.gofruit.space/"+vars["gdps"]+"/db/content")
}

func GetSFXLibraryVersion(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	resp.Header().Set("Location", "https://cdn.fruitspace.one/gdps_sfx/sfxlibrary_version.txt")
	resp.WriteHeader(301)
}

func GetSFXLibrary(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	resp.Header().Set("Location", "https://cdn.fruitspace.one/gdps_sfx/sfxlibrary.dat")
}

func GetSFXTrack(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	vars := gorilla.Vars(req)
	resp.Header().Set("Location", "https://geometrydashfiles.b-cdn.net/sfx/s"+vars["sfxid"]+".ogg")
	resp.WriteHeader(301)
	resp.WriteHeader(301)
}

func RelaySFX(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	vars := gorilla.Vars(req)
	resp.Header().Set("Location", "https://geometrydashfiles.b-cdn.net/sfx/"+vars["sfxid"])
	resp.WriteHeader(301)
}
