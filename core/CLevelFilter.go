package core

import (
	"strconv"
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
	if completed, ok := params["star"]; ok {
		whereq+=" AND starsGot"
		if completed=="0" {whereq+="="}else{whereq+=">"}
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