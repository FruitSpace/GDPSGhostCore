package api

import (
	"HalogenGhostCore/core"
	"fmt"
	gorilla "github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
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

func GetMusicLibraryVersion(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	t := time.Now()
	v := fmt.Sprintf("%s%d", strconv.Itoa(t.Year())[2:], t.YearDay())
	io.WriteString(resp, v)
}

func GetMusicLibrary(resp http.ResponseWriter, req *http.Request, conf *core.GlobalConfig) {
	vars := gorilla.Vars(req)
	logger := core.Logger{Output: os.Stderr}
	config, err := conf.LoadById(vars["gdps"])
	if logger.Should(err) != nil {
		return
	}
	db := &core.MySQLConn{}

	if logger.Should(db.ConnectBlob(config)) != nil {
		return
	}

	url := fmt.Sprintf("https://cdn.fruitspace.one/gdps_sfx/%s_library.dat", vars["gdps"])
	mdata, err := http.Head(url)
	if err != nil {
		fmt.Println(err)
		resp.WriteHeader(500)
		return
	}
	lmod := mdata.Header.Get("last-modified")
	if mdata.StatusCode != 200 {
		lmod = "Mon, 02 Jan 2006 15:04:05 GMT"
	}
	date, err := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", lmod)
	if err != nil {
		fmt.Println(err)
		resp.WriteHeader(500)
		return
	}
	if date.YearDay() != time.Now().YearDay() {
		core.GenerateMusicLibraryFile(db, core.NewS3FS(), vars["gdps"])
	}
	resp.Header().Set("Location", url)
	resp.WriteHeader(301)
}
