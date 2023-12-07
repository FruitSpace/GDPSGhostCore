package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type CMusic struct {
	Status    string
	Id        int
	Name      string
	Artist    string
	Size      float64
	Url       string
	IsBanned  bool
	Downloads int

	DB       *MySQLConn
	Logger   Logger
	Config   *GlobalConfig
	ConfBlob ConfigBlob
}

func (mus *CMusic) Exists(id int) bool {
	var cnt int
	mus.DB.MustQueryRow("SELECT count(*) as cnt FROM #DB#.songs WHERE id=?", id).Scan(&cnt)
	return cnt > 0
}

func (mus *CMusic) RequestNGOuter(id int) bool {
	mus.Id = id
	resp, err := http.Get(mus.Config.ApiEndpoint + "?srvid=" + mus.ConfBlob.ServerConfig.SrvID + "&key=" +
		mus.ConfBlob.ServerConfig.SrvKey + "&action=requestSong&id=" + strconv.Itoa(id))
	if err != nil {
		return false
	}
	rsp, _ := io.ReadAll(resp.Body)
	rsp = bytes.ReplaceAll(rsp, []byte("#"), []byte(""))
	if err = json.Unmarshal(rsp, mus); err != nil {
		fmt.Println(err)
	}

	return mus.Status == "ok"
}

func (mus *CMusic) TransformHalResource() bool {
	arn := strings.Split(mus.Url, ":")
	if !mus.isArnValid(arn) {
		return false
	}
	resp, err := http.Get(mus.Config.ApiEndpoint + "?srvid=" + mus.ConfBlob.ServerConfig.SrvID + "&key=" +
		mus.ConfBlob.ServerConfig.SrvKey + "&action=requestSongARN&type=" + arn[1] + "&id=" + arn[2])
	if err != nil {
		return false
	}
	rsp, _ := io.ReadAll(resp.Body)
	rsp = bytes.ReplaceAll(rsp, []byte("#"), []byte(""))
	bufArtist := mus.Artist
	bufTitle := mus.Name
	if err = json.Unmarshal(rsp, mus); err != nil {
		fmt.Println(err)
	}
	mus.Artist = bufArtist
	mus.Name = bufTitle
	return mus.Status == "ok"
}

func (mus *CMusic) isArnValid(arn []string) bool {
	if len(arn) != 3 {
		return false
	}
	switch arn[1] {
	case "ng":
		if f, _ := regexp.MatchString(`[0-9]`, arn[2]); !f {
			return false
		}
	case "dz":
		if f, _ := regexp.MatchString(`[0-9]`, arn[2]); !f {
			return false
		}
	case "sc":
		if f, _ := regexp.MatchString(`(?i)([a-z\d\-\_])+[\\\\\/]([a-z\d\-\_])+$`, arn[2]); !f {
			return false
		}
	case "yt":
		if f, _ := regexp.MatchString(`(?i)^([a-z\d\-\_])+$`, arn[2]); !f {
			return false
		}
	case "vk":
		if f, _ := regexp.MatchString(`^(\d)+\_(\d)+$`, arn[2]); !f {
			return false
		}
	default:
		return false
	}
	return true
}

func (mus *CMusic) TransformBulkHalResources(links []string) []CMusic {
	var musics []CMusic
	var linksValid []string
	for _, link := range links {
		arn := strings.Split(link, ":")
		if !mus.isArnValid(arn) {
			continue
		}
		linksValid = append(linksValid, arn[1]+":"+arn[2])
	}
	resp, err := http.Get(mus.Config.ApiEndpoint + "?srvid=" + mus.ConfBlob.ServerConfig.SrvID + "&key=" +
		mus.ConfBlob.ServerConfig.SrvKey + "&action=requestSongARNList&ids=" + strings.Join(linksValid, ","))
	if err != nil {
		return musics
	}
	rsp, _ := io.ReadAll(resp.Body)
	rsp = bytes.ReplaceAll(rsp, []byte("#"), []byte(""))
	if err = json.Unmarshal(rsp, &musics); err != nil {
		fmt.Println(err)
	}
	return musics
}

func (mus *CMusic) GetSong(id int) bool {
	if !mus.ConfBlob.ServerConfig.HalMusic {
		return mus.RequestNGOuter(id)
	}
	if !mus.Exists(id) {
		return false
	}
	mus.DB.MustQueryRow("SELECT id,name,artist,size,url,isBanned,downloads FROM #DB#.songs WHERE id=?", id).Scan(
		&mus.Id, &mus.Name, &mus.Artist, &mus.Size, &mus.Url, &mus.IsBanned, &mus.Downloads)
	if mus.IsBanned {
		return false
	}
	if len(mus.Url) > 4 && mus.Url[0:4] == "hal:" {
		return mus.TransformHalResource()
	}
	return true
}

func (mus *CMusic) GetBulkSongs(ids []int) []CMusic {
	q, args, _ := sqlx.In("SELECT id,name,artist,size,url,isBanned,downloads FROM #DB#.songs WHERE id IN (?)", ids)
	rows := mus.DB.ShouldQuery(q, args...)
	defer rows.Close()
	var musics []CMusic
	var queryMusics []CMusic
	var externalQueries []string
	for rows.Next() {
		var xmus CMusic
		rows.Scan(&xmus.Id, &xmus.Name, &xmus.Artist, &xmus.Size, &xmus.Url, &xmus.IsBanned, &xmus.Downloads)
		if xmus.IsBanned {
			continue
		}
		if len(xmus.Url) > 4 && xmus.Url[0:4] == "hal:" {
			externalQueries = append(externalQueries, xmus.Url)
			queryMusics = append(queryMusics, xmus)
		} else {
			musics = append(musics, xmus)
		}
	}
	if len(externalQueries) > 0 {
		deltas := mus.TransformBulkHalResources(externalQueries)
		//musics = append(musics, deltas...)
		for i := 0; i < len(deltas); i++ {
			delta := deltas[i]
			delta.Id = queryMusics[i].Id
			if delta.Status != "ok" {
				musics = append(musics, queryMusics[i])
			} else {
				musics = append(musics, delta)
			}
		}
		//for _, delta := range deltas {
		//	for _, muss := range queryMusics {
		//		log.Println("Checking", muss.Id, "against", delta.Id)
		//		if delta.Id == muss.Id {
		//			delta.Name = muss.Name
		//			delta.Artist = muss.Artist
		//			musics = append(musics, delta)
		//			log.Println("Added", delta.Id)
		//		}
		//	}
		//}
	}
	return musics
}

func (mus *CMusic) UploadSong() int {
	c := mus.DB.ShouldPrepareExec("INSERT INTO #DB#.songs (name,artist,size,url) VALUES (?,?,?,?)", mus.Name, mus.Artist, mus.Size, mus.Url)
	id, _ := c.LastInsertId()
	return int(id)
}

func (mus *CMusic) BanMusic(id int, ban bool) {
	mus.DB.ShouldExec("UPDATE #DB#.songs SET isBanned=? WHERE id=?", ban, id)
}

func (mus *CMusic) CountDownloads() {
	req := mus.DB.MustQuery("SELECT id FROM #DB#.songs")
	defer req.Close()
	for req.Next() {
		var id int
		req.Scan(&id)
		var cnt int
		mus.DB.ShouldQueryRow("SELECT SUM(downloads) FROM #DB#.levels WHERE song_id=?", id).Scan(&cnt)
		mus.DB.ShouldExec("UPDATE #DB#.songs SET downloads=? WHERE id=?", cnt, strconv.Itoa(id))
	}
}

// ! Implement normal API
func (mus *CMusic) GetTopArtists() map[string]string {
	//	resp,err:=http.Get(mus.Config.ApiEndpoint+"?srvid="+mus.ConfBlob.ServerConfig.SrvID+"&key="+mus.ConfBlob.ServerConfig.SrvKey+"&action=requestSong&id="+strconv.Itoa(id))
	//	if err!=nil {return []map[string]string}
	//	rsp,_:=io.ReadAll(resp.Body)
	//	json.Unmarshal(rsp,mus)
	return map[string]string{
		"Riot [Monstercat]": "Monstercat",
		"Noisestorm":        "noisestorm",
		"Nitro Fun":         "NitroFunOfficial",
		"Throttle":          "officialThrottle",
		"Tokyo Machine":     "TokyoMachine",
		"BadComputer":       "BadComputer",
		"Liquid Soul":       "liquidsouliboga",
		"Xtrullor":          "xtrullor",
		"Creo":              "CreoMusic",
		"Dirty Paws":        "DirtyPaws",
		"Have better candidates? Join FruitSpace Discord": "",
	}
}
