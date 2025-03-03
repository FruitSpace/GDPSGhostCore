package connectors

import (
	"HalogenGhostCore/core"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
)

type JSONConnector struct {
	output map[string]interface{}
}

func (c *JSONConnector) Output() string {
	d, _ := json.Marshal(c.output)
	return string(d)
}

func (c *JSONConnector) Error(code string, reason string) {
	c.output["status"] = "error"
	c.output["message"] = reason
	c.output["code"] = code
}

func (c *JSONConnector) Success(message string) {
	c.output["status"] = "success"
	c.output["message"] = message
}

func (c *JSONConnector) NumberedSuccess(id int) {
	c.output["code"] = id
	c.Success("Success indeed")
}

func (c *JSONConnector) Account_Sync(savedata string) {
	c.output["savedata"] = savedata
	c.Success("Savedata present")
}

func (c *JSONConnector) Account_Login(uid int) {
	c.output["uid"] = strconv.Itoa(uid)
	c.Success("Logged in")
}

func (c *JSONConnector) Comment_AccountGet(comments []core.CComment, count int, page int) {
	if len(comments) == 0 {
		c.output["comments"] = []string{}
		c.output["count"] = 0
		c.output["page"] = page
	} else {
		cms := make([]core.CComment, 0)
		for _, comm := range comments {
			if r, err := base64.StdEncoding.DecodeString(comm.Comment); err == nil {
				comm.Comment = string(r)
			}
			cms = append(cms, comm)
		}

		c.output["comments"] = cms
		c.output["count"] = count
		c.output["page"] = page
	}
	c.Success("Comments retrieved")
}

func (c *JSONConnector) Comment_LevelGet(comments []core.CComment, count int, page int) {
	c.Comment_AccountGet(comments, count, page)
}

func (c *JSONConnector) Comment_HistoryGet(comments []core.CComment, acc core.CAccount, role core.Role, count int, page int) {
	c.Comment_AccountGet(comments, count, page)
	c.output["user"] = struct {
		ModBadge       int    `json:"mod_badge"`
		CommentColor   string `json:"comment_color"`
		Uname          string `json:"uname"`
		IconId         int    `json:"icon_id"`
		IconType       int    `json:"icon_type"`
		ColorPrimary   int    `json:"color_primary"`
		ColorSecondary int    `json:"color_secondary"`
		Special        int    `json:"special"`
	}{
		role.ModLevel,
		role.CommentColor,
		acc.Uname,
		acc.GetShownIcon(),
		acc.IconType,
		acc.ColorPrimary,
		acc.ColorSecondary,
		acc.Special,
	}
}

func (c *JSONConnector) Communication_FriendGetRequests(reqs []map[string]string, count int, page int) {
	c.output["requests"] = reqs
	c.output["count"] = count
	c.output["page"] = page
	c.Success("Friend requests retrieved")
}

func (c *JSONConnector) Communication_MessageGet(msg core.CMessage, uid int) {
	if content, err := base64.StdEncoding.DecodeString(msg.Message); err == nil {
		msg.Message = string(content)
	}
	uidx := msg.UidDest
	if uid == msg.UidDest {
		uidx = msg.UidSrc
	}
	xacc := core.CAccount{DB: msg.DB, Uid: uidx}
	xacc.LoadAuth(core.CAUTH_UID)
	c.output["content"] = struct {
		core.CMessage
		Uname string `json:"uname"`
	}{
		msg,
		xacc.Uname,
	}
	c.Success("Message retrieved")
}

func (c *JSONConnector) Communication_MessageGetAll(messages []map[string]string, getSent bool, count int, page int) {
	c.output["messages"] = messages
	c.output["count"] = count
	c.output["page"] = page
	c.output["sent"] = getSent
	c.Success("Messages retrieved")
}

func (c *JSONConnector) Essential_GetMusic(mus core.CMusic) {
	c.output["music"] = mus
	c.Success("Music retrieved")
}

func (c *JSONConnector) Essential_GetTopArtists(artists map[string]string) {
	c.output["artists"] = artists
	c.Success("Top artists retrieved")
}

func (c *JSONConnector) Level_GetGauntlets(gaus []map[string]string, hash string) {
	type r struct {
		PackName string   `json:"pack_name"`
		Levels   []string `json:"levels"`
	}
	var gaunts []r
	for _, gau := range gaus {
		gaunts = append(gaunts, r{
			PackName: gau["pack_name"],
			Levels:   strings.Split(gau["levels"], ","),
		})
	}
	c.output["gauntlets"] = gaunts
	c.output["hash"] = hash
	c.Success("Gauntlets retrieved")
}

func (c *JSONConnector) Level_SearchList(intlists []int, lists []core.CLevelList, count int, page int) {
	var llists []*core.CLevelList
	for _, lid := range intlists {
		for i, list := range lists {
			if list.ID == lid {
				list.DecoupledLevels = strings.Split(list.Levels, ",")
				llists = append(llists, &list)
				lists = append(lists[:i], lists[i+1:]...)
				break
			}
		}
	}
	c.output["lists"] = llists
	c.output["count"] = count
	c.output["page"] = page
	c.Success("Level list retrieved")
}

func (c *JSONConnector) Level_GetMapPacks(packs []core.LevelPack, count int, page int) {
	c.output["packs"] = packs
	c.output["count"] = count
	c.output["page"] = page
	c.Success("Map packs retrieved")
}

func (c *JSONConnector) Level_GetLevelFull(lvl core.CLevel, passwd string, phash string, quest_id int) {
	if txt, err := base64.StdEncoding.DecodeString(lvl.Description); err == nil {
		lvl.Description = string(txt)
	}
	c.output["level"] = lvl
	c.output["quest_id"] = quest_id
	c.Success("Level retrieved")
}

func (c *JSONConnector) Level_GetSpecials(id int, left int) {
	c.output["id"] = id
	c.output["seconds_left"] = left
	c.Success("Specials retrieved")
}

func (c *JSONConnector) Level_SearchLevels(
	intlevels []int, levels []core.CLevel, mus *core.CMusic,
	count int, page int, gdVersion int, gauntlet int,
) {
	var musQueue []int
	musMap := make(map[int]core.CMusic)

	// To keep in order
	var lvls []*core.CLevel
	for _, lvlid := range intlevels {
		for i, lvl := range levels {
			if lvl.Id == lvlid {
				if lvl.SongId != 0 {
					musQueue = append(musQueue, lvl.SongId)
				}
				if ns, err := base64.StdEncoding.DecodeString(lvl.Description); err == nil {
					lvl.Description = string(ns)
				}
				lvls = append(lvls, &lvl)
				levels = append(levels[:i], levels[i+1:]...)
				break
			}
		}
	}

	if len(musQueue) > 0 {
		songs := mus.GetBulkSongs(musQueue)
		for _, sng := range songs {
			musMap[sng.Id] = sng
		}
	}

	c.output["levels"] = lvls
	c.output["music"] = musMap
	c.Success("Levels search completed")
}

func (c *JSONConnector) Rewards_ChallengesOutput(cq core.CQuests, uid int, chk string, udid string) {
	q := cq.GetQuests(uid)
	c.output["challenges"] = q
	c.Success("Challenges retrieved")
}

func (c *JSONConnector) Rewards_ChestOutput(acc core.CAccount, config core.ConfigBlob, udid string, chk string, smallLeft int, bigLeft int, chestType int) {
	c.output["readme"] = "You DO realise chests are imaginary and all the loot is managed by the game client itself? " +
		"I mean the server is just posing some time restrictions and spitting random numbers to make loot seem random. " +
		"Also querying this endpoint usually means that chests timer will be reset, so you won't be able to automate looting " +
		"or just \"peek\" inside chests."
	c.Success("Did absolutely nothing")
}

func (c *JSONConnector) Profile_GetUserProfile(acc core.CAccount, selfUid int) {
	role := acc.GetRoleObj(acc.Uid == selfUid)
	obj := core.NewCAccountJSONFromAccount(acc, &role, acc.Uid == selfUid)
	cf := core.CFriendship{DB: acc.DB}
	obj.WeAreFriends = cf.IsAlreadyFriend(acc.Uid, selfUid)
	c.output["user"] = obj
	c.Success("User retrieved successfully")
}

func (c *JSONConnector) Profile_ListUserProfiles(accs []core.CAccount) {
	var users []core.CAccountJSON

	for _, acc := range accs {
		users = append(users, core.NewCAccountJSONFromAccountLite(acc))
	}

	c.output["users"] = users
	c.Success("User list retrieved successfully")
}

func (c *JSONConnector) Profile_GetSearchableUsers(accs []core.CAccount, count int, page int) {
	var users []core.CAccountJSON
	for _, acc := range accs {
		acc.LoadAuth(core.CAUTH_UID)
		acc.LoadVessels()
		acc.LoadStats()
		users = append(users, core.NewCAccountJSONFromAccountLite(acc))
	}

	c.output["users"] = users
	c.output["count"] = count
	c.output["page"] = page
	c.Success("Users retrieved successfully")
}

func (c *JSONConnector) Score_GetLeaderboard(intaccs []int, xacc core.CAccount) {
	var accs []core.CAccountJSON

	lpos := 0
	for _, uid := range intaccs {
		lpos++
		acc := core.CAccount{DB: xacc.DB, Uid: uid}
		acc.LoadAll()
		ja := core.NewCAccountJSONFromAccountLite(acc)
		ja.LeaderboardRank = lpos
		accs = append(accs, ja)
	}
	c.output["leaderboard"] = accs
	c.Success("Leaderboard retrieved")
}

func (c *JSONConnector) Score_GetScores(scores []core.CScores, mode string) {
	type CoupledScore struct {
		core.CScores
		User core.CAccountJSON `json:"user"`
	}

	var coupledScores []CoupledScore
	for _, score := range scores {
		switch mode {
		case "coins":
			score.Percent = score.Coins
		case "attempts":
			score.Percent = score.Attempts
		}
		acc := core.CAccount{DB: score.DB, Uid: score.Uid}
		acc.LoadAll()
		coupledScores = append(coupledScores, CoupledScore{
			CScores: score,
			User:    core.NewCAccountJSONFromAccountLite(acc),
		})
	}
	c.output["scores"] = coupledScores
	c.Success("Scores retrieved")
}
