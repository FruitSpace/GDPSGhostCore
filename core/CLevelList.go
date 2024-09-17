package core

import "github.com/jmoiron/sqlx"

// Gauntlets
//1 	Fire
//2 	Ice
//3 	Poison
//4 	Shadow
//5 	Lava
//6 	Bonus
//7 	Chaos
//8 	Demon
//9 	Time
//10 	Crystal
//11 	Magic
//12 	spike
//13 	Monster
//14 	Doom
//15 	Death
//16 	Forest
//17 	Rune
//18 	Force
//19 	Spooky
//20 	Dragon
//21 	Water
//22 	Haunted
//23 	Acid
//24 	Witch
//25 	Power
//26 	Potion
//27 	Snake
//28 	Toxic
//29 	Halloween
//30 	Treasure
//31 	Ghost
//32 	Gem
//33 	Inferno
//34 	Portal
//35 	Strange
//36 	Fantasy
//37 	Christmas
//38 	Surprise
//39 	Mystery
//40 	Cursed
//41 	Cyborg
//42 	Castle
//43 	Grave
//44 	Temple

type CLevelList struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Version       int    `json:"version"`
	Difficulty    int    `json:"difficulty"`
	Downloads     int    `json:"downloads"`
	Likes         int    `json:"likes"`
	IsFeatured    bool   `json:"is_featured"`
	UID           int    `json:"uid"`
	Levels        string `json:"levels"`
	Diamonds      int    `json:"diamonds"`
	LevelDiamonds int    `json:"level_diamonds"`
	UploadDate    string `json:"upload_date"`
	UpdateDate    string `json:"update_date"`
	Unlisted      int    `json:"unlisted"`

	SideloadUname *string `json:"sideload_uname,omitempty"`

	DB *MySQLConn `json:"-"`
}

func (cll *CLevelList) Load(id int) {
	cll.DB.MustQueryRow("SELECT id,name,description,version,difficulty,downloads,likes,isFeatured,isUnlisted,uid,levels,diamonds,lvlDiamonds,uploadDate,updateDate FROM #DB#.lists WHERE id=?", id).
		Scan(&cll.ID, &cll.Name, &cll.Description, &cll.Version, &cll.Difficulty, &cll.Downloads, &cll.Likes,
			&cll.IsFeatured, &cll.Unlisted, &cll.UID, &cll.Levels, &cll.Diamonds, &cll.LevelDiamonds, &cll.UploadDate, &cll.UpdateDate)
}

func (cll *CLevelList) Exists(lid int) bool {
	var count int
	cll.DB.MustQueryRow("SELECT COUNT(*) FROM #DB#.lists WHERE id=?", lid).Scan(&count)
	return count > 0
}

func (cll *CLevelList) UpdateList() int {
	if !cll.CheckParams() {
		return -1
	}
	cll.DB.ShouldExec("UPDATE #DB#.lists SET name=?,description=?,version=?,difficulty=?,downloads=?,likes=?,isFeatured=?,isUnlisted=?,uid=?,levels=?,diamonds=?,lvlDiamonds=?,uploadDate=?,updateDate=NOW() WHERE id=?",
		cll.Name, cll.Description, cll.Version, cll.Difficulty, cll.Downloads, cll.Likes, cll.IsFeatured, cll.Unlisted, cll.UID, cll.Levels, cll.Diamonds, cll.LevelDiamonds, cll.UploadDate, cll.ID)

	return cll.ID
}

func (cll *CLevelList) UploadList() int {
	if !cll.CheckParams() {
		return -1
	}
	tx := cll.DB.ShouldPrepareExec("INSERT INTO #DB#.lists (name,description,version,difficulty,isUnlisted,uid,levels) VALUES (?,?,?,?,?,?,?)",
		cll.Name, cll.Description, cll.Version, cll.Difficulty, cll.Unlisted, cll.UID, cll.Levels)
	id, _ := tx.LastInsertId()
	cll.ID = int(id)
	return cll.ID
}

func (cll *CLevelList) CheckParams() bool {
	if len(cll.Name) > 32 || len(cll.Description) > 256 || len(cll.Levels) == 0 {
		return false
	}
	return true
}

func (cll *CLevelList) OnDownloadList() {
	cll.DB.ShouldExec("UPDATE #DB#.lists SET downloads=downloads+1 WHERE id=?", cll.ID)
}

func (cll *CLevelList) LikeList(lid int, uid int, action int) bool {
	if IsLiked(ITEMTYPE_LEVEL, uid, lid, cll.DB) {
		return false
	}
	actionv := "+"
	actions := "Like"
	if action == CLEVEL_ACTION_DISLIKE {
		actionv = "-"
		actions = "Dislike"
	}
	cll.DB.ShouldExec("UPDATE #DB#.lists SET likes=likes"+actionv+"1 WHERE id=?", lid)
	RegisterAction(ACTION_LIST_LIKE, uid, lid, map[string]string{"type": actions}, cll.DB)
	return true
}

func (cll *CLevelList) DeleteList() {
	cll.DB.ShouldExec("DELETE FROM #DB#.lists WHERE id=?", cll.ID)
}

func (cll *CLevelList) IsOwnedBy(uid int) bool {
	cll.Load(cll.ID)
	if cll.ID == 0 {
		return false
	}
	return uid == cll.UID
}

func (cll *CLevelList) LoadBulkSearch(ids []int) []CLevelList {
	var res []CLevelList
	query := "SELECT id,name,description,version,difficulty,downloads,likes,isFeatured,isUnlisted,#DB#.lists.uid,levels,#DB#.lists.diamonds," +
		"lvlDiamonds,uploadDate,updateDate, #DB#.users.uname FROM #DB#.lists LEFT JOIN #DB#.users on #DB#.lists.uid=#DB#.users.uid WHERE id IN(?)"
	q, args, _ := sqlx.In(query, ids)
	rows := cll.DB.MustQuery(q, args...)
	defer rows.Close()
	for rows.Next() {
		levl := CLevelList{DB: cll.DB}
		e := rows.Scan(&levl.ID, &levl.Name, &levl.Description, &levl.Version, &levl.Difficulty, &levl.Downloads, &levl.Likes,
			&levl.IsFeatured, &levl.Unlisted, &levl.UID, &levl.Levels, &levl.Diamonds, &levl.LevelDiamonds, &levl.UploadDate, &levl.UpdateDate,
			&levl.SideloadUname)
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

// !TEMP
func (cll *CLevelList) Preload() {
	cll.DB.ShouldExec(`
CREATE TABLE IF NOT EXISTS #DB#.lists
(
    id                   int(11)          NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name                 varchar(32)      NOT NULL DEFAULT 'Unnamed',
    description          varchar(256)     NOT NULL DEFAULT '',
    uid                  int(11)          NOT NULL DEFAULT 0,
    version              tinyint          NOT NULL DEFAULT 1,
    difficulty           tinyint          NOT NULL DEFAULT -1,
    downloads            int              NOT NULL DEFAULT 0,
    likes                int              NOT NULL DEFAULT 0,
    isFeatured           tinyint(1)       NOT NULL DEFAULT 0,
    isUnlisted           tinyint(1)       NOT NULL DEFAULT 0,
    levels               mediumtext       NOT NULL DEFAULT '',
    diamonds             int              NOT NULL DEFAULT 0,
    lvlDiamonds          int              NOT NULL DEFAULT 0,
    uploadDate           DATETIME         NOT NULL DEFAULT NOW(),
    updateDate           DATETIME         NOT NULL DEFAULT NOW()
)
`)
}
