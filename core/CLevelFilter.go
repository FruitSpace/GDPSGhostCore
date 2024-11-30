package core

import (
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	CLEVELFILTER_MOSTLIKED      int = 700
	CLEVELFILTER_MOSTDOWNLOADED int = 701
	CLEVELFILTER_TRENDING       int = 702
	CLEVELFILTER_LATEST         int = 703
	CLEVELFILTER_MAGIC          int = 704
	CLEVELFILTER_HALL           int = 705
	CLEVELFILTER_SAFE_DAILY     int = 706
	CLEVELFILTER_SAFE_WEEKLY    int = 707
	CLEVELFILTER_SAFE_EVENT     int = 708
	CLEVELFILTER_SENT           int = 709
)

type CLevelFilter struct {
	DB    *MySQLConn
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
	whereq := ""
	// Demon difficulty filter, else normal
	if demonDiff, ok := params["demonDiff"]; ok {
		var demonDiffI int
		TryInt(&demonDiffI, demonDiff)
		if demonDiff == "0" {
			whereq += " AND demonDifficulty>=0"
		} else {
			whereq += " AND demonDifficulty=" + strconv.Itoa(demonDiffI)
		}
	} else {
		if diff, ok := params["diff"]; ok {
			diff = QuickComma(diff)
			whereq += " AND difficulty IN (" + diff + ") AND demonDifficulty=-1"
		}
	}

	// Length filter
	if clen, ok := params["length"]; ok {
		clen = QuickComma(clen)
		whereq += " AND length IN (" + clen + ")"
	}

	//Completed/uncompleted stuff
	if completed, ok := params["completed"]; ok {
		whereq += " AND id"
		if completed == "0" {
			whereq += " NOT"
		}
		whereq += " IN (" + QuickComma(params["completedLevels"]) + ")"
	}

	if _, ok := params["isFeatured"]; ok {
		whereq += " AND isFeatured=1"
	}
	if _, ok := params["is2p"]; ok {
		whereq += " AND is2p=1"
	}
	if _, ok := params["isOrig"]; ok {
		whereq += " AND original_id=0"
	}
	// rateconfig
	{
		var ratestring []string
		if _, ok := params["isFeatured"]; ok {
			ratestring = append(ratestring, "isFeatured=1")
		}
		if _, ok := params["isEpic"]; ok {
			ratestring = append(ratestring, "isEpic=1")
		}
		if _, ok := params["isMythic"]; ok {
			ratestring = append(ratestring, "isEpic=2")
		}
		if _, ok := params["isLegendary"]; ok {
			ratestring = append(ratestring, "isEpic=3")
		}

		if len(ratestring) > 0 {
			whereq += " AND (" + strings.Join(ratestring, " OR ") + ")"
		}
	}

	if _, ok := params["coins"]; ok {
		whereq += " AND coins>0"
	}

	//Is starred
	if star, ok := params["star"]; ok {
		whereq += " AND starsGot"
		if star == "0" {
			whereq += "="
		} else {
			whereq += ">"
		}
		whereq += "0"
	}

	//Song Custom/Classic stuff
	if songid, ok := params["songid"]; ok {
		whereq += " AND song_id="
		var sid int
		TryInt(&sid, songid)
		if sid < 0 {
			whereq += "0 AND track_id=" + strconv.Itoa(-1*sid+1)
		} else {
			whereq += strconv.Itoa(sid)
		}
	}

	return whereq
}

// SearchLevels searches Levels with filters given
func (filter *CLevelFilter) SearchLevels(page int, params map[string]string, xtype int) []int {
	page = int(math.Abs(float64(page))) * 10
	suffix := filter.GenerateQueryString(params)
	query := " FROM #DB#.levels WHERE versionGame<=?"
	orderBy := ""

	switch xtype {
	case CLEVELFILTER_MOSTLIKED:
		orderBy = "likes DESC, downloads DESC"
	case CLEVELFILTER_MOSTDOWNLOADED:
		orderBy = "downloads DESC, likes DESC"
	case CLEVELFILTER_TRENDING:
		date := time.Now().AddDate(0, 0, -7).Format("2006-01-02 15:04:05")
		query += " AND uploadDate>'" + date + "'"
		orderBy = "likes DESC, downloads DESC"
	case CLEVELFILTER_LATEST:
		orderBy = "uploadDate DESC, downloads DESC"
	case CLEVELFILTER_MAGIC:
		orderBy = "uploadDate DESC, downloads DESC"
		query += " AND objects>9999 AND length>=3 AND original_id=0"
	case CLEVELFILTER_SENT:
		orderBy = "uploadDate DESC, downloads DESC"
		query += " AND EXISTS (SELECT id FROM #DB#.rateQueue WHERE #DB#.levels.id = #DB#.rateQueue.lvl_id)"
	case CLEVELFILTER_HALL:
		query += " AND isEpic>=1"
		orderBy = "likes DESC, downloads DESC"
	// Here be The Safe
	case CLEVELFILTER_SAFE_DAILY:
		query += " AND EXISTS (SELECT id FROM #DB#.quests WHERE #DB#.levels.id = #DB#.quests.lvl_id AND #DB#.quests.type=0)"
		orderBy = "uploadDate DESC, downloads DESC"
	case CLEVELFILTER_SAFE_WEEKLY:
		query += " AND EXISTS (SELECT id FROM #DB#.quests WHERE #DB#.levels.id = #DB#.quests.lvl_id AND #DB#.quests.type=1)"
		orderBy = "uploadDate DESC, downloads DESC"
	case CLEVELFILTER_SAFE_EVENT:
		query += " AND EXISTS (SELECT id FROM #DB#.quests WHERE #DB#.levels.id = #DB#.quests.lvl_id AND #DB#.quests.type=-1)"
		orderBy = "uploadDate DESC, downloads DESC"
	default:
		query += " AND 1=0" //Because I can
	}
	sortstr := " ORDER BY " + orderBy + " LIMIT 10 OFFSET " + strconv.Itoa(page)

	var levels []int

	//If we actually search for something
	if sterm, ok := params["sterm"]; ok {
		// If it's an ID
		if _, err := strconv.Atoi(sterm); err == nil {
			compq := query + " AND id=?" + suffix
			rows := filter.DB.ShouldQuery("SELECT id"+compq+sortstr, params["versionGame"], sterm)
			defer rows.Close()
			filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq, params["versionGame"], sterm).Scan(&filter.Count)
			for rows.Next() {
				var lvlid int
				rows.Scan(&lvlid)
				levels = append(levels, lvlid)
			}
		} else {
			// But if it's just text we search title
			//! To support unlisted2 aka friendList maybe we should use isUnlisted<>1 or "isUnlisted=ANY(0"+",2"+")
			compq := query + " AND name LIKE ? AND isUnlisted=0" + suffix
			rows := filter.DB.ShouldQuery("SELECT id"+compq+sortstr, params["versionGame"], "%"+sterm+"%")
			defer rows.Close()
			filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq, params["versionGame"], "%"+sterm+"%").Scan(&filter.Count)
			for rows.Next() {
				var lvlid int
				rows.Scan(&lvlid)
				levels = append(levels, lvlid)
			}
		}
	} else {
		// Or if we're just wandering and clicking buttons
		compq := query + " AND isUnlisted=0" + suffix
		rows := filter.DB.ShouldQuery("SELECT id"+compq+sortstr, params["versionGame"])
		defer rows.Close()
		filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq, params["versionGame"]).Scan(&filter.Count)
		for rows.Next() {
			var lvlid int
			rows.Scan(&lvlid)
			levels = append(levels, lvlid)
		}
	}

	return levels
}

// SearchUserLevels searches levels of Followed users or by UID
func (filter *CLevelFilter) SearchUserLevels(page int, params map[string]string, followMode bool) []int {
	page = int(math.Abs(float64(page))) * 10
	suffix := filter.GenerateQueryString(params)
	query := " FROM #DB#.levels WHERE versionGame<=?"
	sortstr := " ORDER BY downloads DESC LIMIT 10 OFFSET " + strconv.Itoa(page)

	var levels []int

	if sterm, ok := params["sterm"]; ok {
		if followMode {
			if _, err := strconv.Atoi(sterm); err != nil {
				query += " AND isUnlisted=0 AND name LIKE ?"
				sterm = "%" + sterm + "%"
			} else {
				query += " AND id=?"
			}
			compq := query + " AND uid IN (" + QuickComma(params["followList"]) + ")" + suffix
			rows := filter.DB.ShouldQuery("SELECT id"+compq+sortstr, params["versionGame"], sterm)
			defer rows.Close()
			filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq, params["versionGame"], sterm).Scan(&filter.Count)
			for rows.Next() {
				var lvlid int
				rows.Scan(&lvlid)
				levels = append(levels, lvlid)
			}
		} else {
			if stermi, err := strconv.Atoi(sterm); err == nil {
				compq := query + " AND uid=?" + suffix
				rows := filter.DB.ShouldQuery("SELECT id"+compq+sortstr, params["versionGame"], stermi)
				defer rows.Close()
				filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq, params["versionGame"], stermi).Scan(&filter.Count)
				for rows.Next() {
					var lvlid int
					rows.Scan(&lvlid)
					levels = append(levels, lvlid)
				}
			}
		}
	} else {
		if followMode {
			compq := query + " AND isUnlisted=0 AND uid IN (" + QuickComma(params["followList"]) + ")" + suffix
			rows := filter.DB.ShouldQuery("SELECT id"+compq+sortstr, params["versionGame"])
			defer rows.Close()
			filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq, params["versionGame"]).Scan(&filter.Count)
			for rows.Next() {
				var lvlid int
				rows.Scan(&lvlid)
				levels = append(levels, lvlid)
			}
		} else {
			compq := query + suffix
			rows := filter.DB.ShouldQuery("SELECT id"+compq+sortstr, params["versionGame"])
			defer rows.Close()
			filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq, params["versionGame"]).Scan(&filter.Count)
			for rows.Next() {
				var lvlid int
				rows.Scan(&lvlid)
				levels = append(levels, lvlid)
			}
		}
	}

	return levels
}

// SearchListLevels searches levels for Gauntlets/Mappacks via sterm list
func (filter *CLevelFilter) SearchListLevels(page int, params map[string]string) []int {
	page = int(math.Abs(float64(page))) * 10
	suffix := filter.GenerateQueryString(params)
	query := " FROM #DB#.levels WHERE versionGame<=?"
	sortstr := " LIMIT 10 OFFSET " + strconv.Itoa(page)

	var levels []int

	if sterm, ok := params["sterm"]; ok {
		compq := query + " AND id IN (" + QuickComma(sterm) + ")" + suffix
		rows := filter.DB.ShouldQuery("SELECT id"+compq+sortstr, params["versionGame"])
		defer rows.Close()
		filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq, params["versionGame"]).Scan(&filter.Count)
		for rows.Next() {
			var lvlid int
			rows.Scan(&lvlid)
			levels = append(levels, lvlid)
		}
	}

	return levels
}

// GetGauntlets retrieves gauntlet list and levels (w/ trailing hash)
func (filter *CLevelFilter) GetGauntlets() (gauntlets []map[string]string, hashString string) {
	rows := filter.DB.ShouldQuery("SELECT packName, levels FROM #DB#.levelpacks WHERE packType=1 ORDER BY CAST(packname as int)")
	defer rows.Close()
	for rows.Next() {
		var packName, levels string
		rows.Scan(&packName, &levels)
		if len(Decompose(CleanDoubles(levels, ","), ",")) != 5 {
			continue
		}
		if _, err := strconv.Atoi(packName); err != nil {
			continue
		}
		gauntlets = append(gauntlets, map[string]string{"pack_name": packName, "levels": levels})
		hashString += packName + levels
	}
	return gauntlets, HashSolo2(hashString)
}

// GetGauntletLevels returns gauntlet level IDs
func (filter *CLevelFilter) GetGauntletLevels(gau int) []int {
	var levels string
	filter.DB.ShouldQueryRow("SELECT levels FROM #DB#.levelpacks WHERE packType=1 AND packName=? LIMIT 1", gau).Scan(&levels)
	malevels := Decompose(CleanDoubles(levels, ","), ",")
	if len(malevels) < 5 {
		return []int{}
	}
	return []int{malevels[0], malevels[1], malevels[2], malevels[3], malevels[4]}
}

func (filter *CLevelFilter) CountMapPacks() int {
	var cnt int
	filter.DB.ShouldQueryRow("SELECT count(*) FROM #DB#.levelpacks WHERE packType=0").Scan(&cnt)
	return cnt
}

// GetMapPacks retrieves MapPacks list and levels (w/ trailing hash)
func (filter *CLevelFilter) GetMapPacks(page int) (packs []LevelPack, count int) {
	page = int(math.Abs(float64(page))) * 10
	rows := filter.DB.ShouldQuery("SELECT id,packName,levels,packStars,packCoins,packDifficulty,packColor FROM #DB#.levelpacks WHERE packType=0 LIMIT 10 OFFSET " + strconv.Itoa(page))
	defer rows.Close()

	for rows.Next() {
		var pack LevelPack
		rows.Scan(&pack.Id, &pack.PackName, &pack.Levels, &pack.PackStars, &pack.PackCoins, &pack.PackDifficulty, &pack.PackColor)
		packs = append(packs, pack)
	}
	if len(packs) > 0 {
		count = filter.CountMapPacks()
	}
	return
}

type LevelPack struct {
	Id             int    `json:"id"`
	PackName       string `json:"pack_name"`
	Levels         string `json:"levels"`
	PackStars      int    `json:"pack_stars"`
	PackCoins      int    `json:"pack_coins"`
	PackDifficulty int    `json:"pack_difficulty"`
	PackColor      string `json:"pack_color"`
}
