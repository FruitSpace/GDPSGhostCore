package core

import (
	"math"
	"strconv"
	"time"
)

const (
	CLEVELLISTFILTER_MOSTLIKED      int = 800
	CLEVELLISTFILTER_MOSTDOWNLOADED int = 801
	CLEVELLISTFILTER_TRENDING       int = 802
	CLEVELLISTFILTER_LATEST         int = 803
	CLEVELLISTFILTER_MAGIC          int = 804
	CLEVELLISTFILTER_AWARDED        int = 805
	CLEVELLISTFILTER_SENT           int = 806
)

type CLevelListFilter struct {
	DB    *MySQLConn
	Count int
}

/*
 * --- [PARAMS] Object ---
 * + (s) sterm - search term (used with other params to specify exact usage)
 * + (s) diff - difficulties array. If it doesn't exist then all diffs | Ex: array(10,20,30)
 * + (b) star - if not set then all, else star/nostar
 * + (s) followList - who to follow
 */

// GenerateQueryString generates SQL string out of params
func (filter *CLevelListFilter) GenerateQueryString(params map[string]string) string {
	whereq := ""
	// Difficulty
	if diff, ok := params["diff"]; ok {
		diff = QuickComma(diff)
		whereq += " AND difficulty IN (" + diff + ")"
	}

	//Is starred
	if star, ok := params["star"]; ok {
		whereq += " AND diamonds"
		if star == "0" {
			whereq += "="
		} else {
			whereq += ">"
		}
		whereq += "0"
	}

	return whereq
}

// SearchLevels searches Levels with filters given
func (filter *CLevelListFilter) SearchLists(page int, params map[string]string, xtype int) []int {
	page = int(math.Abs(float64(page))) * 10
	suffix := filter.GenerateQueryString(params)
	query := " FROM #DB#.lists WHERE 1=1" //placeholder
	orderBy := ""

	switch xtype {
	case CLEVELLISTFILTER_MOSTLIKED:
		orderBy = "likes DESC, downloads DESC"
	case CLEVELLISTFILTER_MOSTDOWNLOADED:
		orderBy = "downloads DESC, likes DESC"
	case CLEVELLISTFILTER_TRENDING:
		date := time.Now().AddDate(0, 0, -7).Format("2006-01-02 15:04:05")
		query += " AND uploadDate>'" + date + "'"
		orderBy = "likes DESC, downloads DESC"
	case CLEVELLISTFILTER_LATEST:
		orderBy = "uploadDate DESC, downloads DESC"
	case CLEVELLISTFILTER_MAGIC:
		orderBy = "uploadDate DESC, downloads DESC"
		query += "" // robtop sniffed some coke
	case CLEVELLISTFILTER_AWARDED:
		orderBy = "uploadDate DESC, downloads DESC"
		query += " AND isFeatured>0 AND diamonds>0"
	case CLEVELLISTFILTER_SENT:
		orderBy = "uploadDate DESC, downloads DESC"
		query += " AND isFeatured=0 AND diamonds=0"
	default:
		query += " AND 1=0" //Because I can
	}
	sortstr := " ORDER BY " + orderBy + " LIMIT 10 OFFSET " + strconv.Itoa(page)

	var lists []int

	//If we actually search for something
	if sterm, ok := params["sterm"]; ok {
		// If it's an ID
		if _, err := strconv.Atoi(sterm); err == nil {
			compq := query + " AND id=?" + suffix
			rows := filter.DB.ShouldQuery("SELECT id"+compq+sortstr, sterm)
			defer rows.Close()
			filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq, sterm).Scan(&filter.Count)
			for rows.Next() {
				var lid int
				rows.Scan(&lid)
				lists = append(lists, lid)
			}
		} else {
			// But if it's just text we search title
			//! To support unlisted2 aka friendList maybe we should use isUnlisted<>1 or "isUnlisted=ANY(0"+",2"+")
			compq := query + " AND name LIKE ? AND isUnlisted=0" + suffix
			rows := filter.DB.ShouldQuery("SELECT id"+compq+sortstr, "%"+sterm+"%")
			defer rows.Close()
			filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq, "%"+sterm+"%").Scan(&filter.Count)
			for rows.Next() {
				var lid int
				rows.Scan(&lid)
				lists = append(lists, lid)
			}
		}
	} else {
		// Or if we're just wandering and clicking buttons
		compq := query + " AND isUnlisted=0" + suffix
		rows := filter.DB.ShouldQuery("SELECT id" + compq + sortstr)
		defer rows.Close()
		filter.DB.ShouldQueryRow("SELECT count(*) as cnt" + compq).Scan(&filter.Count)
		for rows.Next() {
			var lid int
			rows.Scan(&lid)
			lists = append(lists, lid)
		}
	}

	return lists
}

// SearchUserLevels searches levels of Followed users or by UID
func (filter *CLevelListFilter) SearchUserLists(page int, params map[string]string, followMode bool) []int {
	page = int(math.Abs(float64(page))) * 10
	suffix := filter.GenerateQueryString(params)
	query := " FROM #DB#.lists WHERE 1=1" //placehodler
	sortstr := " ORDER BY downloads DESC LIMIT 10 OFFSET " + strconv.Itoa(page)

	var lists []int

	if sterm, ok := params["sterm"]; ok {
		if followMode {
			if _, err := strconv.Atoi(sterm); err != nil {
				query += " AND isUnlisted=0 AND name LIKE ?"
				sterm = "%" + sterm + "%"
			} else {
				query += " AND id=?"
			}
			compq := query + " AND uid IN (" + QuickComma(params["followList"]) + ")" + suffix
			rows := filter.DB.ShouldQuery("SELECT id"+compq+sortstr, sterm)
			defer rows.Close()
			filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq, sterm).Scan(&filter.Count)
			for rows.Next() {
				var lid int
				rows.Scan(&lid)
				lists = append(lists, lid)
			}
		} else {
			if stermi, err := strconv.Atoi(sterm); err == nil {
				compq := query + " AND uid=?" + suffix
				rows := filter.DB.ShouldQuery("SELECT id"+compq+sortstr, stermi)
				defer rows.Close()
				filter.DB.ShouldQueryRow("SELECT count(*) as cnt"+compq, stermi).Scan(&filter.Count)
				for rows.Next() {
					var lid int
					rows.Scan(&lid)
					lists = append(lists, lid)
				}
			}
		}
	} else {
		if followMode {
			compq := query + " AND isUnlisted=0 AND uid IN (" + QuickComma(params["followList"]) + ")" + suffix
			rows := filter.DB.ShouldQuery("SELECT id" + compq + sortstr)
			defer rows.Close()
			filter.DB.ShouldQueryRow("SELECT count(*) as cnt" + compq).Scan(&filter.Count)
			for rows.Next() {
				var lid int
				rows.Scan(&lid)
				lists = append(lists, lid)
			}
		} else {
			compq := query + suffix
			rows := filter.DB.ShouldQuery("SELECT id" + compq + sortstr)
			defer rows.Close()
			filter.DB.ShouldQueryRow("SELECT count(*) as cnt" + compq).Scan(&filter.Count)
			for rows.Next() {
				var lid int
				rows.Scan(&lid)
				lists = append(lists, lid)
			}
		}
	}

	return lists
}
