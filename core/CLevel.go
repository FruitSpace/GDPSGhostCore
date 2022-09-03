package core

import (
	"strings"
	"time"
)

const (
	CLEVEL_ACTION_LIKE int = 300
	CLEVEL_ACTION_DISLIKE int = 301
)

type CLevel struct {
	//Main
	Id int
	Name string
	Description string
	Uid int
	Password string
	Version int
	Length int
	Difficulty int
	DemonDifficulty int

	SuggestDifficulty float64
	SuggestDifficultyCnt int

	//Level
	TrackId int
	SongId int
	VersionGame int
	VersionBinary int
	StringExtra string
	StringSettings string
	StringLevel string
	StringLevelInfo string
	OrigId int

	//Stats
	Objects int
	StarsRequested int
	StarsGot int
	Ucoins int
	Coins int
	Downloads int
	Likes int
	Reports int

	//Params
	Is2p bool
	IsVerified bool
	IsFeatured bool
	ISHall bool
	IsEpic int
	IsUnlisted int
	IsLDM bool

	//Dates
	UploadDate string
	UpdateDate string

	UnlockLevelObject bool

	DB MySQLConn
	Logger Logger
}

func (lvl *CLevel) Exists(lvlid int) bool {
	var v int
	lvl.DB.DB.QueryRow("SELECT uid FROM levels WHERE id=?",lvlid).Scan(&v)
	return v>0
}

func (lvl *CLevel) CountLevels() int {
	var cnt int
	lvl.DB.MustQueryRow("SELECT count(*) as cnt FROM levels").Scan(&cnt)
	return cnt
}

func (lvl *CLevel) LoadParams() {
	lvl.DB.MustQueryRow("SELECT is2p, isVerified, isFeatured, isHall, isEpic, isUnlisted, isLDM FROM levels WHERE id=?",lvl.Id).Scan(
		&lvl.Is2p,&lvl.IsVerified,&lvl.IsFeatured,&lvl.ISHall,&lvl.IsEpic,&lvl.IsUnlisted,&lvl.IsLDM)
}

func (lvl *CLevel) PushParams() {
	lvl.DB.ShouldQuery("UPDATE levels SET is2p=?,isVerified=?,isFeatured=?,isHall=?,isEpic=?,isUnlisted=?,isLDM=? WHERE id=?",
		lvl.Is2p,lvl.IsVerified,lvl.IsFeatured,lvl.ISHall,lvl.IsEpic,lvl.IsUnlisted,lvl.IsLDM,lvl.Id)
}

func (lvl *CLevel) LoadDates() {
	lvl.DB.MustQueryRow("SELECT uploadDate, updateDate FROM levels WHERE id=?",lvl.Id).Scan(&lvl.UploadDate,&lvl.UpdateDate)
}

func (lvl *CLevel) LoadLevel() {
	lvl.DB.MustQueryRow("SELECT track_id, song_id,versionGame,versionBinary,stringExtra,stringSettings,stringLevel,stringLevelInfo,original_id FROM levels WHERE id=?",lvl.Id).Scan(
		&lvl.TrackId,&lvl.SongId,&lvl.VersionGame,&lvl.VersionBinary,&lvl.StringExtra,&lvl.StringSettings,&lvl.StringLevel,&lvl.StringLevelInfo,&lvl.OrigId)
}

func (lvl *CLevel) LoadStats() {
	lvl.DB.MustQueryRow("SELECT objects,starsRequested,starsGot,ucoins,coins,downloads,likes,reports FROM levels WHERE id=?",lvl.Id).Scan(
		&lvl.Objects,&lvl.StarsRequested,&lvl.StarsGot,&lvl.Ucoins,&lvl.Coins,&lvl.Downloads,&lvl.Likes,&lvl.Reports)
}

func (lvl *CLevel) OnDownloadLevel() {
	lvl.DB.ShouldQuery("UPDATE levels SET downloads=downloads+1 WHERE id=?",lvl.Id)
}

func (lvl *CLevel) LoadMain() {
	lvl.DB.MustQueryRow("SELECT name,description,uid,password,version,length,difficulty,demonDifficulty,suggestDifficulty,suggestDifficultyCnt FROM levels WHERE id=?",lvl.Id).Scan(
		&lvl.Name,&lvl.Description,&lvl.Uid,&lvl.Password,&lvl.Version,&lvl.Length,&lvl.Difficulty,&lvl.DemonDifficulty,&lvl.SuggestDifficulty,&lvl.SuggestDifficultyCnt)
}

func (lvl *CLevel) LoadAll() {
	query:="SELECT name,description,uid,password,version,length,difficulty,demonDifficulty,suggestDifficulty,suggestDifficultyCnt,"+
		"track_id,song_id,versionGame,versionBinary,stringExtra,stringSettings,stringLevel,stringLevelInfo,original_id,objects,"+
		"starsRequested,starsGot,ucoins,coins,downloads,likes,reports,is2p,isVerified,isFeatured,isHall,isEpic,isUnlisted,isLDM," +
		"uploadDate,updateDate FROM levels WHERE id=?"
	lvl.DB.MustQueryRow(query,lvl.Id).Scan(&lvl.Name,&lvl.Description,&lvl.Uid,&lvl.Password,&lvl.Version,&lvl.Length,&lvl.Difficulty,
		&lvl.DemonDifficulty,&lvl.SuggestDifficulty,&lvl.SuggestDifficultyCnt,&lvl.TrackId,&lvl.SongId,&lvl.VersionGame,&lvl.VersionBinary,
		&lvl.StringExtra,&lvl.StringSettings,&lvl.StringLevel,&lvl.StringLevelInfo,&lvl.OrigId,&lvl.Objects,&lvl.StarsRequested,
		&lvl.StarsGot,&lvl.Ucoins,&lvl.Coins,&lvl.Downloads,&lvl.Likes,&lvl.Reports,&lvl.Is2p,&lvl.IsVerified,&lvl.IsFeatured,
		&lvl.ISHall,&lvl.IsEpic,&lvl.IsUnlisted,&lvl.IsLDM,&lvl.UploadDate,&lvl.UpdateDate)
	//lvl.LoadMain()
	//lvl.LoadLevel()
	//lvl.LoadStats()
	//lvl.LoadParams()
	//lvl.LoadDates()
}

func (lvl *CLevel) LoadBase() {
	lvl.DB.MustQueryRow("SELECT uid,name FROM levels WHERE id=?",lvl.Id).Scan(&lvl.Uid,&lvl.Name)
}

func (lvl *CLevel) IsOwnedBy(uid int) bool {
	if !lvl.Exists(lvl.Id) {return false}
	lvl.LoadBase()
	return uid==lvl.Uid
}

func (lvl *CLevel) CheckParams() bool {
	if len(lvl.Name)>32 || len(lvl.Description)>256 || len(lvl.Password)>8 || lvl.Version<1 || lvl.Version>127 || lvl.TrackId<0 || lvl.SongId<0 || lvl.VersionGame<0 {return false}
	if  lvl.VersionBinary<0 || len(lvl.StringLevel)<16 ||lvl.OrigId<0 || (lvl.Objects<100 && !lvl.UnlockLevelObject) || lvl.StarsRequested<0 || lvl.StarsRequested>10 || lvl.Ucoins<0 || lvl.Ucoins>3 {return false}
	return true
}

func (lvl *CLevel) DeleteLevel() {
	lvl.DB.ShouldQuery("DELETE FROM levels WHERE id=?",lvl.Id)
}

func (lvl *CLevel) UploadLevel() int {
	if !lvl.CheckParams() {return -1}
	date:=time.Now().Format("2006-01-02 15:04:05")
	q:="INSERT INTO levels (name, description, uid, password, version, length, track_id, song_id, versionGame, versionBinary, stringExtra, stringSettings, stringLevel, stringLevelInfo, original_id, objects, starsRequested, ucoins, is2p, isVerified, isUnlisted, isLDM, uploadDate, updateDate) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	r:=lvl.DB.ShouldPrepareExec(q,lvl.Name,lvl.Description,lvl.Uid,lvl.Password,lvl.Version,lvl.Length,lvl.TrackId,lvl.SongId,lvl.VersionGame,lvl.VersionBinary,lvl.StringExtra,lvl.StringSettings,lvl.StringLevel,lvl.StringLevelInfo,lvl.OrigId,lvl.Objects,lvl.StarsRequested,lvl.Ucoins,lvl.Is2p,lvl.IsVerified,lvl.IsUnlisted,lvl.IsLDM,date,date)
	id,_:=r.LastInsertId()
	return int(id)
}

func (lvl *CLevel) UpdateLevel() int {
	if !lvl.CheckParams() {return -1}
	date:=time.Now().Format("2006-01-02 15:04:05")
	q:="UPDATE levels SET name=?, description=?, password=?, version=?, length=?, track_id=?, song_id=?, versionGame=?, versionBinary=?, stringExtra=?, stringSettings=?, stringLevel=?, stringLevelInfo=?, original_id=?, objects=?, starsRequested=?, ucoins=?, is2p=?, isVerified=?, isUnlisted=?, isLDM=?, updateDate=? WHERE id=?"
	lvl.DB.ShouldQuery(q,lvl.Name,lvl.Description,lvl.Password,lvl.Version,lvl.Length,lvl.TrackId,lvl.SongId,lvl.VersionGame,lvl.VersionBinary,lvl.StringExtra,lvl.StringSettings,lvl.StringLevel,lvl.StringLevelInfo,lvl.OrigId,lvl.Objects,lvl.StarsRequested,lvl.Ucoins,lvl.Is2p,lvl.IsVerified,lvl.IsUnlisted,lvl.IsLDM,date,lvl.Id)
	return lvl.Id
}

func (lvl *CLevel) UpdateDescription(desc string) bool {
	if len(desc)>256 {return false}
	lvl.DB.ShouldQuery("UPDATE levels SET description=? WHERE id=?",desc,lvl.Id)
	return true
}

func (lvl *CLevel) DoSuggestDifficulty(diffx int){
	diff:=float64(diffx)
	lvl.SuggestDifficulty=(lvl.SuggestDifficulty*float64(lvl.SuggestDifficultyCnt)+diff)/float64(lvl.SuggestDifficultyCnt+1)
	lvl.DB.ShouldQuery("UPDATE levels SET suggestDifficulty=?,suggestDifficultyCnt=? WHERE id=?",lvl.SuggestDifficulty,lvl.SuggestDifficultyCnt,lvl.Id)
}

func (lvl *CLevel) RateLevel(stars int) {
	lvl.StarsGot=stars
	postfix:=",demonDifficulty=-1"
	var diff int
	switch stars {
	case 1:
		diff=-1 //Auto
	case 2:
		diff=10 //Easy
	case 3:
		diff=20 //Normal
	case 4:
		fallthrough
	case 5:
		diff=30 //Hard
	case 6:
		fallthrough
	case 7:
		diff=40 //Harder
	case 8:
		fallthrough
	case 9:
		diff=50 //Insane
	case 10:
		diff=50 //Demon
		postfix=",demonDifficulty=3"
	default:
		diff=0 //N/A Unrated
	}
	lvl.DB.ShouldQuery("UPDATE levels SET difficulty=?,starsGot=?"+postfix+" WHERE id=?",diff,stars,lvl.Id)
	lvl.RecalculateCPoints(lvl.Uid)
}

func (lvl *CLevel) RateDemon(diff int) {
	var xdiff int
	switch diff {
	case 5:
		xdiff=6
	case 4:
		xdiff=5
	case 3:
		xdiff=0
	case 2:
		xdiff=4
	default:
		xdiff=3
	}
	lvl.DB.ShouldQuery("UPDATE levels SET demonDifficulty=? WHERE id=?",xdiff,lvl.Id)
}

func (lvl *CLevel) FeatureLevel(feature bool) {
	var featured int
	if feature {featured=1}
	lvl.DB.ShouldQuery("UPDATE levels SET isFeatured=? WHERE id=?",featured,lvl.Id)
	lvl.RecalculateCPoints(lvl.Uid)
}

func (lvl *CLevel) EpicLevel(epic bool) {
	var epicd int
	if epic {epicd=1}
	lvl.DB.ShouldQuery("UPDATE levels SET isEpic=? WHERE id=?",epicd,lvl.Id)
	lvl.RecalculateCPoints(lvl.Uid)
}

func (lvl *CLevel) LikeLevel(lvlid int, uid int, action int) bool {
	if IsLiked(ITEMTYPE_LEVEL, uid, lvlid, lvl.DB) {return false}
	actionv:="+"
	actions:="Like"
	if action == CLEVEL_ACTION_DISLIKE {
		actionv="-"
		actions="Dislike"
	}
	lvl.DB.ShouldQuery("UPDATE levels SET likes=likes"+actionv+"1 WHERE id=?",lvlid)
	RegisterAction(ACTION_LEVEL_LIKE, uid, lvlid, map[string]string{"type":actions},lvl.DB)
	return true
}

func (lvl *CLevel) VerifyCoins(verify bool) {
	cc:="0"
	if verify {cc="ucoins"}
	lvl.DB.ShouldQuery("UPDATE levels SET coins="+cc+" WHERE id=?",lvl.Id)
}

func (lvl *CLevel) ReportLevel() {
	lvl.DB.ShouldQuery("UPDATE levels SET reports=reports+1 WHERE id=?",lvl.Id)
}

func (lvl *CLevel) RecalculateCPoints(uid int) {
	req:=lvl.DB.MustQuery("SELECT starsGot,isFeatured,isEpic,collab FROM levels WHERE uid=?",uid)
	defer req.Close()
	for req.Next() {
		var (
			starsGot int
			isFeatured int
			isEpic int
			collab string
			cpoints int
		)
		req.Scan(&starsGot,&isFeatured,&isEpic,&collab)
		if starsGot>0 {cpoints++}
		if isFeatured>0 {cpoints++}
		if isEpic>0 {cpoints++}
		collablist:=strings.Split(collab,",")
		for _,colid:=range collablist {
			if colid=="" {continue}
			lvl.DB.ShouldQuery("UPDATE users SET cpoints=cpoints+? WHERE uid=?",cpoints,colid)
		}
		lvl.DB.ShouldQuery("UPDATE users SET cpoints=cpoints+? WHERE uid=?",cpoints,uid)
	}
}


func (lvl *CLevel) SendReq(modUid int, stars int, isFeatured int) bool {
	var cnt int
	lvl.DB.MustQueryRow("SELECT count(*) AS cnt FROM rateQueue WHERE lvl_id=? AND mod_uid=?",lvl.Id,modUid).Scan(&cnt)
	if cnt!=0 {return false}
	lvl.DB.ShouldQuery("INSERT INTO rateQueue (lvl_id,name,uid,mod_uid,stars,isFeatured) VALUES(?,?,?,?,?,?)",
		lvl.Id,lvl.Name,lvl.Uid,modUid,stars,isFeatured)
	return true
}