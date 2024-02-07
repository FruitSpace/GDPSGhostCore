package core

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GenerateMusicLibraryFile(db *MySQLConn, fs *S3FS, srvid string) string {
	clean := func(s string) string {
		s = strings.ReplaceAll(s, ",", " ")
		s = strings.ReplaceAll(s, ";", " ")
		s = strings.ReplaceAll(s, "|", " ")
		s = strings.ReplaceAll(s, "#", " ")
		return s
	}

	t := time.Now()
	ver := fmt.Sprintf("%s%d", strconv.Itoa(t.Year())[2:], t.YearDay())

	tags := []string{
		"1,NewGrounds",
		"2,YouTube",
		"3,Deezer",
		"4,VK",
		"5,Dropbox",
	}

	var artists []string
	artistsLookup := make(map[string]int)
	rows := db.MustQuery("SELECT DISTINCT artist FROM #DB#.songs")
	defer rows.Close()
	for rows.Next() {
		var artist string
		rows.Scan(&artist)
		artists = append(artists, fmt.Sprintf("%d,%s, , ", len(artists)+1, clean(artist)))
		artistsLookup[artist] = len(artists)
	}

	var tracks []string

	songs := db.MustQuery("SELECT id,name,artist,size,url,isBanned,downloads FROM #DB#.songs")
	defer songs.Close()
	for songs.Next() {
		var song CMusic
		songs.Scan(&song.Id, &song.Name, &song.Artist, &song.Size, &song.Url, &song.IsBanned, &song.Downloads)
		if song.IsBanned {
			continue
		}
		songstr := fmt.Sprintf("%d,%s,%d,%.0f,69,.%s", song.Id, clean(song.Name), artistsLookup[song.Artist], song.Size*1024*1024, getTagByARN(song.Url))
		tracks = append(tracks, songstr)
	}

	preparedblock := strings.Join([]string{
		ver,
		strings.Join(artists, ";"),
		strings.Join(tracks, ";"),
		strings.Join(tags, ";"),
	}, "|")

	// zlib compress preparedblock
	var compressedBlock bytes.Buffer
	w := zlib.NewWriter(&compressedBlock)
	w.Write([]byte(preparedblock))
	w.Close()
	preparedblock = base64.StdEncoding.EncodeToString(compressedBlock.Bytes())
	preparedblock = strings.ReplaceAll(preparedblock, "/", "_")
	preparedblock = strings.ReplaceAll(preparedblock, "+", "-")

	fs.PutFile(fmt.Sprintf("/gdps_sfx/%s_library.dat", srvid), []byte(preparedblock))

	return fmt.Sprintf("/gdps_sfx/%s_library.dat", srvid)
}

func getTagByARN(arn string) string {
	if strings.HasPrefix(arn, "hal:") {
		arnType := strings.Split(arn, ":")[1]
		switch arnType {
		case "ng":
			return "1"
		case "yt":
			return "2"
		case "dz":
			return "3"
		case "vk":
			return "4"
		default:
			return "5"
		}
	}
	return "5"
}
