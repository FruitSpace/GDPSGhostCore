package core

import (
	"strconv"
	"strings"
	"time"
)

const (
	CLEVELFILTER_MOSTLIKED int = 700
	CLEVELFILTER_MOSTDOWNLOADED int = 701
	CLEVELFILTER_TRENDING int = 702
	CLEVELFILTER_LATEST int = 703
	CLEVELFILTER_MAGIC int = 704
	CLEVELFILTER_HALL int = 705
)

type CLevelFilter struct {
	DB MySQLConn
	Count int
}

/*
 * --- [PARAMS] Object ---
 * + (s) sterm - search term (used with other params to specify exact usage)
 * + (s) diff - difficulties array. If it doesn't exist then all diffs | Ex: array(10,20,30)
 * + (i) demonDiff - demon difficulty (auto sorted), if doesnt exist - ignore
 * + (i) length - level length array (if not specified -> all) | Ex: array(0,1,4)
 * + (b) completed - if not set than all, else only completed/uncompleted
 * + (s) completedLevels - if completed is set then list of comp/uncomp is sent
 * + (b) isFeatured - obv
 * + (b) isOrig - where origid=0
 * + (b) is2p - straightforward
 * + (b) coins - should we search coins only
 * + (b) isEpic - also obv
 * + (b) star - if not set then all, else star/nostar
 * + (i) songid - official song id or custom song id
 * + (b) songCustom - if set then songid is custom song
 *
 * !! Demon overrides diff
 */

// GenerateQueryString generates SQL string out of params
func (filter *CLevelFilter) GenerateQueryString(params map[string]string) string {
	whereq:=""
	// Demon difficulty filter, else normal
	if demonDiff, ok := params["demonDiff"]; ok {
		var demonDiffI int
		TryInt(&demonDiffI, demonDiff)
		if demonDiff=="0" {
			whereq+=" AND demonDifficulty>=0"
		}else{
			whereq+=" AND demonDifficulty="+strconv.Itoa(demonDiffI)
		}
	}else{
		if diff, ok := params["diff"]; ok {
			diff=QuickComma(diff)
			whereq+=" AND difficulty IN ("+diff+") AND demonDifficulty=-1"
		}
	}

	// Length filter
	if clen, ok := params["length"]; ok {
		clen=QuickComma(clen)
		whereq+=" AND length IN ("+clen+")"
	}

	//Completed/uncompleted stuff
	if completed, ok := params["completed"]; ok {
		whereq+=" AND id"
		if completed=="0" {whereq+=" NOT"}
		whereq+=" IN ("+QuickComma(params["completedLevels"])+")"
	}

	if _,ok:=params["isFeatured"]; ok {whereq+=" AND isFeatured=1"}
	if _,ok:=params["is2p"]; ok {whereq+=" AND is2p=1"}
	if _,ok:=params["isOrig"]; ok {whereq+=" AND original_id=0"}
	if _,ok:=params["isEpic"]; ok {whereq+=" AND isFeatured=1"}
	if _,ok:=params["isFeatured"]; ok {whereq+=" AND isEpic=1"}
	if _,ok:=params["coins"]; ok {whereq+=" AND coins>0"}

	//Is starred
	if star, ok := params["star"]; ok {
		whereq+=" AND starsGot"
		if star=="0" {whereq+="="}else{whereq+=">"}
		whereq+="0"
	}

	//Song Custom/Classic stuff
	if songid, ok := params["songid"]; ok {
		whereq+=" AND song_id="
		var sid int
		TryInt(&sid, songid)
		if sid<0 {whereq+="0 AND track_id="+strconv.Itoa(-1*sid+1)}else{whereq+=strconv.Itoa(sid)}
	}

	return whereq
}

func (filter *CLevelFilter) SearchLevels(page int, params map[string]string, xtype int) []int {
	suffix:=filter.GenerateQueryString(params)
	query:=" FROM levels WHERE versionGame<=?"
	orderBy:=""

	switch xtype {
	case CLEVELFILTER_MOSTLIKED:
		orderBy="likes DESC, downloads DESC"
		break
	case CLEVELFILTER_MOSTDOWNLOADED:
		orderBy="downloads DESC, likes DESC"
		break
	case CLEVELFILTER_TRENDING:
		date:=time.Now().AddDate(0,0,-7).Format("2006-01-02 15:04:05")
		query+=" AND uploadDate>'"+date+"'"
		orderBy="likes DESC, downloads DESC"
		break
	case CLEVELFILTER_LATEST:
		orderBy="uploadDate DESC, downloads DESC"
		break
	case CLEVELFILTER_MAGIC:
		orderBy="uploadDate DESC, downloads DESC"
		if strings.Contains(suffix,"starsGot>0") {
			// Old magic
			query+=" AND objects>9999 AND length>=3 AND original_id=0"
		}else{
			// New magic
			query+=" WHERE EXISTS (SELECT id FROM rateQueue WHERE levels.id = rateQueue.lvl_id)"
		}
		break
	case CLEVELFILTER_HALL:
		query+=" AND isEpic=1"
		orderBy="likes DESC, downloads DESC"
		break
	default:
		query+=" AND 1=0" //Because I can
	}
	sortstr:=" ORDER BY "+orderBy+" LIMIT 10 OFFSET"+strconv.Itoa(page)

	var levels []int

	//If we actually search for something
	if sterm, ok := params["sterm"]; ok {
		// If it's an ID
		if _,err := strconv.Atoi(sterm); err==nil {
			compq:=query+" AND id=?"+suffix
			rows:=filter.DB.ShouldQuery("SELECT id"+compq+sortstr,params["versionGame"],sterm)
			filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq,params["versionGame"],sterm).Scan(&filter.Count)
			for rows.Next() {
				var lvlid int
				rows.Scan(&lvlid)
				levels=append(levels,lvlid)
			}
		}else{
			// But if it's just text we search title
			compq:=query+" AND name LIKE ? AND isUnlisted=0"+suffix
			rows:=filter.DB.ShouldQuery("SELECT id"+compq+sortstr,params["versionGame"],sterm)
			filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq,params["versionGame"],sterm).Scan(&filter.Count)
			for rows.Next() {
				var lvlid int
				rows.Scan(&lvlid)
				levels=append(levels,lvlid)
			}
		}
	}else{
		// Or if we're just wandering and clicking buttons
		compq:=query+" AND isUnlisted=0"+suffix
		rows:=filter.DB.ShouldQuery("SELECT id"+compq+sortstr,params["versionGame"])
		filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq,params["versionGame"]).Scan(&filter.Count)
		for rows.Next() {
			var lvlid int
			rows.Scan(&lvlid)
			levels=append(levels,lvlid)
		}
	}

	return levels
}