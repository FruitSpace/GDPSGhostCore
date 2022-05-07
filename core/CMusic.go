package core

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type CMusic struct {
	Status string
	Id int
	Name string
	Artist string
	Size string
	Url string
	IsBanned bool
	Downloads int

	DB MySQLConn
	Logger Logger
	Config GlobalConfig
	ConfBlob ConfigBlob
}

func (mus *CMusic) Exists(id int) bool {
	var cnt int
	mus.Logger.Must(mus.DB,mus.DB.DB.QueryRow("SELECT count(*) as cnt FROM songs WHERE id=?",id).Scan(&cnt))
	return cnt>0
}

func (mus *CMusic) RequestNGOuter(id int) bool {
	resp,err:=http.Get(mus.Config.ApiEndpoint+"?id="+mus.ConfBlob.ServerConfig.SrvID+"&key="+mus.ConfBlob.ServerConfig.SrvKey+"&action=requestSong&id="+strconv.Itoa(id))
	if err!=nil {return false}
	rsp,_:=io.ReadAll(resp.Body)
	json.Unmarshal(rsp,mus)

	return mus.Status=="ok"
}

func (mus *CMusic) TransformHalResource() bool {
	arn:=strings.Split(mus.Url,":")
	if len(arn)!=3 {return false}
	switch arn[1] {
	case "ng":
		if f,_:=regexp.MatchString(`[0-9]`,arn[2]); !f {return false}
		break
	case "sc":
		if f,_:=regexp.MatchString(`(?i)([a-z\d\-\_])+[\\\\\/]([a-z\d\-\_])+$`,arn[2]); !f {return false}
		break
	case "yt":
		if f,_:=regexp.MatchString(`(?i)^([a-z\d\-\_])+$`,arn[2]); !f {return false}
		break
	case "vk":
		if f,_:=regexp.MatchString(`^(\d)+\_(\d)+$`,arn[2]); !f {return false}
		break
	default:
		return false
	}
	resp,err:=http.Get(mus.Config.ApiEndpoint+"?id="+mus.ConfBlob.ServerConfig.SrvID+"&key="+mus.ConfBlob.ServerConfig.SrvKey+"&action=requestSongARN&type="+arn[1]+"&id="+arn[2])
	if err!=nil {return false}
	rsp,_:=io.ReadAll(resp.Body)
	json.Unmarshal(rsp,mus)
	return mus.Status=="ok"
}

func (mus *CMusic) GetSong(id int) bool {
	if !mus.ConfBlob.ServerConfig.HalMusic {return mus.RequestNGOuter(id)}
	if !mus.Exists(id) {return false}
	mus.Logger.Must(mus.DB,mus.DB.DB.QueryRow("SELECT id,name,artist,size,url,isBanned,downloads FROM songs WHERE id=?",id).Scan(
		&mus.Id,&mus.Name,&mus.Artist,&mus.Size,&mus.Url,&mus.IsBanned,&mus.Downloads))
	if mus.IsBanned {return false}
	if mus.Url[0:4]=="hal:" {return mus.TransformHalResource()}
	return true
}

func (mus *CMusic) UploadSong() int {
	req,_:=mus.DB.DB.Prepare("INSERT INTO songs (name,artist,size,url) VALUES (?,?,?,?)")
	c,_:=req.Exec(mus.Name,mus.Artist,mus.Size,mus.Url)
	id,_:=c.LastInsertId()
	return int(id)
}

func (mus *CMusic) BanMusic(id int, ban bool) {
	var banc int
	if ban{banc=1}
	mus.DB.DB.Query("UPDATE songs SET isBanned=? WHERE id=?",banc,id)
}

func (mus *CMusic) CountDownloads() {
	req,_:=mus.DB.DB.Query("SELECT id FROM songs")
	for req.Next() {
		var id int
		req.Scan(&id)
		creq,_:=mus.DB.DB.Query("SELECT downloads FROM levels WHERE song_id=?",id)
		var cnt int
		for creq.Next() {
			var downs int
			creq.Scan(&downs)
			cnt+=downs
		}
		mus.DB.DB.Query("UPDATE songs SET downloads=? WHERE id=?",cnt,id)
	}
}