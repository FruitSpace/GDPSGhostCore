package core

import (
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	CLEVEL_ACTION_LIKE    int = 300
	CLEVEL_ACTION_DISLIKE int = 301
)

type CLevel struct {
	//Main
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Uid             int    `json:"uid"`
	Password        string `json:"password"`
	Version         int    `json:"version"`
	Length          int    `json:"length"`
	Difficulty      int    `json:"difficulty"`
	DemonDifficulty int    `json:"demon_difficulty"`

	SuggestDifficulty    float64 `json:"-"`
	SuggestDifficultyCnt int     `json:"-"`

	//Level
	TrackId         int    `json:"track_id"`
	SongId          int    `json:"song_id"`
	VersionGame     int    `json:"version_game"`
	VersionBinary   int    `json:"version_binary"`
	StringExtra     string `json:"string_extra"`
	StringSettings  string `json:"string_settings"` // FUCKING SONG IDS ; SFX IDS (song;sfx)
	StringLevel     string `json:"string_level"`
	StringLevelInfo string `json:"string_level_info"`
	OrigId          int    `json:"original_id"`

	//Stats
	Objects        int `json:"objects"`
	StarsRequested int `json:"stars_requested"`
	StarsGot       int `json:"stars_got"`
	Ucoins         int `json:"ucoins"`
	Coins          int `json:"coins"`
	Downloads      int `json:"downloads"`
	Likes          int `json:"likes"`
	Reports        int `json:"reports"`

	//Params
	Is2p       bool `json:"is_2p"`
	IsVerified bool `json:"is_verified"`
	IsFeatured int  `json:"is_featured"`
	ISHall     bool `json:"is_hall"`
	IsEpic     int  `json:"is_epic"`
	IsUnlisted int  `json:"is_unlisted"`
	IsLDM      bool `json:"is_ldm"`

	//Dates
	UploadDate string `json:"upload_date"`
	UpdateDate string `json:"update_date"`

	UnlockLevelObject bool    `json:"-"`
	SideloadUname     *string `json:"username,omitempty"`

	DB     *MySQLConn `json:"-"`
	Logger Logger     `json:"-"`
}

func (lvl *CLevel) Exists(lvlid int) bool {
	var v int
	lvl.DB.ShouldQueryRow("SELECT uid FROM #DB#.levels WHERE id=?", lvlid).Scan(&v)
	return v > 0
}

func (lvl *CLevel) CountLevels() int {
	var cnt int
	lvl.DB.MustQueryRow("SELECT count(*) as cnt FROM #DB#.levels").Scan(&cnt)
	return cnt
}

func (lvl *CLevel) LoadParams() {
	lvl.DB.MustQueryRow("SELECT is2p, isVerified, isFeatured, isHall, isEpic, isUnlisted, isLDM FROM #DB#.levels WHERE id=?", lvl.Id).Scan(
		&lvl.Is2p, &lvl.IsVerified, &lvl.IsFeatured, &lvl.ISHall, &lvl.IsEpic, &lvl.IsUnlisted, &lvl.IsLDM)
}

func (lvl *CLevel) PushParams() {
	lvl.DB.ShouldExec("UPDATE #DB#.levels SET is2p=?,isVerified=?,isFeatured=?,isHall=?,isEpic=?,isUnlisted=?,isLDM=? WHERE id=?",
		lvl.Is2p, lvl.IsVerified, lvl.IsFeatured, lvl.ISHall, lvl.IsEpic, lvl.IsUnlisted, lvl.IsLDM, lvl.Id)
}

func (lvl *CLevel) LoadDates() {
	lvl.DB.MustQueryRow("SELECT uploadDate, updateDate FROM #DB#.levels WHERE id=?", lvl.Id).Scan(&lvl.UploadDate, &lvl.UpdateDate)
}

func (lvl *CLevel) LoadLevel() {
	lvl.DB.MustQueryRow("SELECT track_id, song_id,versionGame,versionBinary,stringExtra,stringSettings,stringLevel,stringLevelInfo,original_id FROM #DB#.levels WHERE id=?", lvl.Id).Scan(
		&lvl.TrackId, &lvl.SongId, &lvl.VersionGame, &lvl.VersionBinary, &lvl.StringExtra, &lvl.StringSettings, &lvl.StringLevel, &lvl.StringLevelInfo, &lvl.OrigId)
}

func (lvl *CLevel) LoadStats() {
	lvl.DB.MustQueryRow("SELECT objects,starsRequested,starsGot,ucoins,coins,downloads,likes,reports FROM #DB#.levels WHERE id=?", lvl.Id).Scan(
		&lvl.Objects, &lvl.StarsRequested, &lvl.StarsGot, &lvl.Ucoins, &lvl.Coins, &lvl.Downloads, &lvl.Likes, &lvl.Reports)
}

func (lvl *CLevel) OnDownloadLevel() {
	lvl.DB.ShouldExec("UPDATE #DB#.levels SET downloads=downloads+1 WHERE id=?", lvl.Id)
}

func (lvl *CLevel) LoadMain() {
	lvl.DB.MustQueryRow("SELECT name,description,uid,password,version,length,difficulty,demonDifficulty,suggestDifficulty,suggestDifficultyCnt FROM #DB#.levels WHERE id=?", lvl.Id).Scan(
		&lvl.Name, &lvl.Description, &lvl.Uid, &lvl.Password, &lvl.Version, &lvl.Length, &lvl.Difficulty, &lvl.DemonDifficulty, &lvl.SuggestDifficulty, &lvl.SuggestDifficultyCnt)
}

func (lvl *CLevel) LoadAll() {
	query := "SELECT name,description,uid,password,version,length,difficulty,demonDifficulty,suggestDifficulty,suggestDifficultyCnt," +
		"track_id,song_id,versionGame,versionBinary,stringExtra,stringSettings,stringLevel,stringLevelInfo,original_id,objects," +
		"starsRequested,starsGot,ucoins,coins,downloads,likes,reports,is2p,isVerified,isFeatured,isHall,isEpic,isUnlisted,isLDM," +
		"uploadDate,updateDate FROM #DB#.levels WHERE id=?"
	lvl.DB.MustQueryRow(query, lvl.Id).Scan(&lvl.Name, &lvl.Description, &lvl.Uid, &lvl.Password, &lvl.Version, &lvl.Length, &lvl.Difficulty,
		&lvl.DemonDifficulty, &lvl.SuggestDifficulty, &lvl.SuggestDifficultyCnt, &lvl.TrackId, &lvl.SongId, &lvl.VersionGame, &lvl.VersionBinary,
		&lvl.StringExtra, &lvl.StringSettings, &lvl.StringLevel, &lvl.StringLevelInfo, &lvl.OrigId, &lvl.Objects, &lvl.StarsRequested,
		&lvl.StarsGot, &lvl.Ucoins, &lvl.Coins, &lvl.Downloads, &lvl.Likes, &lvl.Reports, &lvl.Is2p, &lvl.IsVerified, &lvl.IsFeatured,
		&lvl.ISHall, &lvl.IsEpic, &lvl.IsUnlisted, &lvl.IsLDM, &lvl.UploadDate, &lvl.UpdateDate)
	//lvl.LoadMain()
	//lvl.LoadLevel()
	//lvl.LoadStats()
	//lvl.LoadParams()
	//lvl.LoadDates()
}

func (lvl *CLevel) LoadBase() {
	lvl.DB.MustQueryRow("SELECT uid,name FROM #DB#.levels WHERE id=?", lvl.Id).Scan(&lvl.Uid, &lvl.Name)
}

func (lvl *CLevel) IsOwnedBy(uid int) bool {
	if !lvl.Exists(lvl.Id) {
		return false
	}
	lvl.LoadBase()
	return uid == lvl.Uid
}

func (lvl *CLevel) CheckParams() bool {
	if len(lvl.Name) > 32 || len(lvl.Description) > 256 || len(lvl.Password) > 8 || lvl.Version < 1 || lvl.Version > 127 || lvl.TrackId < 0 || lvl.SongId < 0 || lvl.VersionGame < 0 {
		return false
	}
	if lvl.VersionBinary < 0 || len(lvl.StringLevel) < 16 || lvl.OrigId < 0 || (lvl.Objects < 100 && !lvl.UnlockLevelObject) || lvl.StarsRequested < 0 || lvl.StarsRequested > 10 || lvl.Ucoins < 0 || lvl.Ucoins > 3 {
		return false
	}
	return true
}

func (lvl *CLevel) DeleteLevel() {
	lvl.DB.ShouldExec("DELETE FROM #DB#.levels WHERE id=?", lvl.Id)
	lvl.RecalculateCPoints(lvl.Uid)
}

func (lvl *CLevel) UploadLevel() int {
	if !lvl.CheckParams() {
		return -1
	}
	date := time.Now().Format("2006-01-02 15:04:05")
	q := "INSERT INTO #DB#.levels (name, description, uid, password, version, length, track_id, song_id, versionGame, versionBinary, stringExtra, stringSettings, stringLevel, stringLevelInfo, original_id, objects, starsRequested, ucoins, is2p, isVerified, isUnlisted, isLDM, uploadDate, updateDate) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	r := lvl.DB.ShouldPrepareExec(q, lvl.Name, lvl.Description, lvl.Uid, lvl.Password, lvl.Version, lvl.Length, lvl.TrackId, lvl.SongId, lvl.VersionGame, lvl.VersionBinary, lvl.StringExtra, lvl.StringSettings, lvl.StringLevel, lvl.StringLevelInfo, lvl.OrigId, lvl.Objects, lvl.StarsRequested, lvl.Ucoins, lvl.Is2p, lvl.IsVerified, lvl.IsUnlisted, lvl.IsLDM, date, date)
	id, _ := r.LastInsertId()
	return int(id)
}

func (lvl *CLevel) UpdateLevel() int {
	if !lvl.CheckParams() {
		return -1
	}
	date := time.Now().Format("2006-01-02 15:04:05")
	q := "UPDATE #DB#.levels SET name=?, description=?, password=?, version=?, length=?, track_id=?, song_id=?, versionGame=?, versionBinary=?, stringExtra=?, stringSettings=?, stringLevel=?, stringLevelInfo=?, original_id=?, objects=?, starsRequested=?, ucoins=?, is2p=?, isVerified=?, isUnlisted=?, isLDM=?, updateDate=? WHERE id=?"
	lvl.DB.ShouldExec(q, lvl.Name, lvl.Description, lvl.Password, lvl.Version, lvl.Length, lvl.TrackId, lvl.SongId, lvl.VersionGame, lvl.VersionBinary, lvl.StringExtra, lvl.StringSettings, lvl.StringLevel, lvl.StringLevelInfo, lvl.OrigId, lvl.Objects, lvl.StarsRequested, lvl.Ucoins, lvl.Is2p, lvl.IsVerified, lvl.IsUnlisted, lvl.IsLDM, date, lvl.Id)
	return lvl.Id
}

func (lvl *CLevel) UpdateDescription(desc string) bool {
	if len(desc) > 256 {
		return false
	}
	lvl.DB.ShouldExec("UPDATE #DB#.levels SET description=? WHERE id=?", desc, lvl.Id)
	return true
}

func (lvl *CLevel) DoSuggestDifficulty(diffx int) {
	diff := float64(diffx)
	lvl.SuggestDifficulty = (lvl.SuggestDifficulty*float64(lvl.SuggestDifficultyCnt) + diff) / float64(lvl.SuggestDifficultyCnt+1)
	lvl.DB.ShouldExec("UPDATE #DB#.levels SET suggestDifficulty=?,suggestDifficultyCnt=? WHERE id=?", lvl.SuggestDifficulty, lvl.SuggestDifficultyCnt, lvl.Id)
}

func (lvl *CLevel) RateLevel(stars int) {
	lvl.StarsGot = stars
	postfix := ",demonDifficulty=-1"
	var diff int
	switch stars {
	case 1:
		diff = -1 //Auto
	case 2:
		diff = 10 //Easy
	case 3:
		diff = 20 //Normal
	case 4:
		fallthrough
	case 5:
		diff = 30 //Hard
	case 6:
		fallthrough
	case 7:
		diff = 40 //Harder
	case 8:
		fallthrough
	case 9:
		diff = 50 //Insane
	case 10:
		diff = 50 //Demon
		postfix = ",demonDifficulty=3"
	default:
		diff = 0 //N/A Unrated
	}
	lvl.DB.ShouldExec("UPDATE #DB#.levels SET difficulty=?,starsGot=?"+postfix+" WHERE id=?", diff, stars, lvl.Id)
	lvl.RecalculateCPoints(lvl.Uid)
}

func (lvl *CLevel) RateDemon(diff int) {
	var xdiff int
	switch diff {
	case 5:
		xdiff = 6
	case 4:
		xdiff = 5
	case 3:
		xdiff = 0
	case 2:
		xdiff = 4
	default:
		xdiff = 3
	}
	lvl.DemonDifficulty = xdiff
	lvl.DB.ShouldExec("UPDATE #DB#.levels SET demonDifficulty=? WHERE id=?", xdiff, lvl.Id)
}

func (lvl *CLevel) FeatureLevel(featured int) bool {
	if featured == 0 && lvl.IsEpic > 0 {
		return false
	}
	lvl.DB.ShouldExec("UPDATE #DB#.levels SET isFeatured=? WHERE id=?", featured, lvl.Id)
	lvl.RecalculateCPoints(lvl.Uid)
	return true
}

func (lvl *CLevel) EpicLevel(epic bool) {
	var epicd int
	if epic {
		epicd = 1
	}
	lvl.DB.ShouldExec("UPDATE #DB#.levels SET isEpic=? WHERE id=?", epicd, lvl.Id)
	lvl.RecalculateCPoints(lvl.Uid)
}
func (lvl *CLevel) LegendaryLevel(legend bool) {
	var epicd int
	if legend {
		epicd = 2
	}
	lvl.DB.ShouldExec("UPDATE #DB#.levels SET isEpic=? WHERE id=?", epicd, lvl.Id)
	lvl.RecalculateCPoints(lvl.Uid)
}
func (lvl *CLevel) MythicLevel(legend bool) {
	var epicd int
	if legend {
		epicd = 3
	}
	lvl.DB.ShouldExec("UPDATE #DB#.levels SET isEpic=? WHERE id=?", epicd, lvl.Id)
	lvl.RecalculateCPoints(lvl.Uid)
}

func (lvl *CLevel) LikeLevel(lvlid int, uid int, action int) bool {
	if IsLiked(ITEMTYPE_LEVEL, uid, lvlid, lvl.DB) {
		return false
	}
	actionv := "+"
	actions := "Like"
	if action == CLEVEL_ACTION_DISLIKE {
		actionv = "-"
		actions = "Dislike"
	}
	lvl.DB.ShouldExec("UPDATE #DB#.levels SET likes=likes"+actionv+"1 WHERE id=?", lvlid)
	RegisterAction(ACTION_LEVEL_LIKE, uid, lvlid, map[string]string{"type": actions}, lvl.DB)
	return true
}

func (lvl *CLevel) VerifyCoins(verify bool) {
	cc := "0"
	if verify {
		cc = "ucoins"
	}
	lvl.DB.ShouldExec("UPDATE #DB#.levels SET coins="+cc+" WHERE id=?", lvl.Id)
}

func (lvl *CLevel) ReportLevel() {
	lvl.DB.ShouldExec("UPDATE #DB#.levels SET reports=reports+1 WHERE id=?", lvl.Id)
}

func (lvl *CLevel) RecalculateCPoints(uid int) {
	req := lvl.DB.MustQuery("SELECT starsGot,isFeatured,isEpic,collab FROM #DB#.levels WHERE uid=?", uid)
	defer req.Close()
	totalCP := 0
	for req.Next() {
		var (
			starsGot   int
			isFeatured int
			isEpic     int
			collab     string
			cpoints    int = 0
		)
		req.Scan(&starsGot, &isFeatured, &isEpic, &collab)
		if starsGot > 0 {
			cpoints++
		}
		if isFeatured > 0 {
			cpoints++
		}
		if isEpic > 0 {
			cpoints++
		}
		//! COLLABS DISABLED
		//collablist := strings.Split(collab, ",")
		//for _, colid := range collablist {
		//	if colid == "" {
		//		continue
		//	}
		//	lvl.DB.ShouldExec("UPDATE #DB#.users SET cpoints=cpoints+? WHERE uid=?", cpoints, colid)
		//}
		totalCP += cpoints
	}
	lvl.DB.ShouldExec("UPDATE #DB#.users SET cpoints=? WHERE uid=?", totalCP, uid)
}

func (lvl *CLevel) SendReq(modUid int, stars int, isFeatured int) bool {
	var cnt int
	lvl.DB.MustQueryRow("SELECT count(*) AS cnt FROM #DB#.rateQueue WHERE lvl_id=? AND mod_uid=?", lvl.Id, modUid).Scan(&cnt)
	if cnt != 0 {
		return false
	}
	lvl.DB.ShouldExec("INSERT INTO #DB#.rateQueue (lvl_id,name,uid,mod_uid,stars,isFeatured) VALUES(?,?,?,?,?,?)",
		lvl.Id, lvl.Name, lvl.Uid, modUid, stars, isFeatured)
	return true
}

func (lvl *CLevel) LoadBulkSearch(ids []int) []CLevel {
	var res []CLevel
	query := "SELECT id,name,description,#DB#.levels.uid,password,version,length,difficulty,demonDifficulty,suggestDifficulty,suggestDifficultyCnt," +
		"track_id,song_id,versionGame,versionBinary,stringExtra,stringSettings,stringLevelInfo,original_id,objects," +
		"starsRequested,starsGot,#DB#.levels.ucoins,#DB#.levels.coins,downloads,likes,reports,is2p,isVerified,isFeatured,isHall,isEpic,isUnlisted,isLDM," +
		"uploadDate,updateDate, #DB#.users.uname FROM #DB#.levels LEFT JOIN #DB#.users on #DB#.levels.uid=#DB#.users.uid WHERE id IN(?)"
	q, args, _ := sqlx.In(query, ids)
	rows := lvl.DB.MustQuery(q, args...)
	defer rows.Close()
	for rows.Next() {
		levl := CLevel{DB: lvl.DB}
		e := rows.Scan(&levl.Id, &levl.Name, &levl.Description, &levl.Uid, &levl.Password, &levl.Version, &levl.Length, &levl.Difficulty,
			&levl.DemonDifficulty, &levl.SuggestDifficulty, &levl.SuggestDifficultyCnt, &levl.TrackId, &levl.SongId, &levl.VersionGame, &levl.VersionBinary,
			&levl.StringExtra, &levl.StringSettings, &levl.StringLevelInfo, &levl.OrigId, &levl.Objects, &levl.StarsRequested,
			&levl.StarsGot, &levl.Ucoins, &levl.Coins, &levl.Downloads, &levl.Likes, &levl.Reports, &levl.Is2p, &levl.IsVerified, &levl.IsFeatured,
			&levl.ISHall, &levl.IsEpic, &levl.IsUnlisted, &levl.IsLDM, &levl.UploadDate, &levl.UpdateDate, &levl.SideloadUname)
		if levl.SideloadUname == nil {
			s := "[DELETED]"
			levl.SideloadUname = &s
		}
		if e != nil {
			SendMessageDiscord(e.Error())
		}
		res = append(res, levl)
	}

	return res
}
