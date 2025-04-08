package connectors

import (
	"HalogenGhostCore/core"
	"time"
)

var loc, _ = time.LoadLocation("Europe/Moscow")

type Connector interface {
	Output() string

	Error(code string, reason string)
	Success(message string)
	NumberedSuccess(id int)
	Account_Sync(savedata string)
	Account_Login(uid int)
	Comment_AccountGet(comments []core.CComment, count int, page int)
	Comment_LevelGet(comments []core.CComment, count int, page int)
	Comment_HistoryGet(comments []core.CComment, acc core.CAccount, role core.Role, count int, page int)
	Communication_FriendGetRequests(reqs []map[string]string, count int, page int)
	Communication_MessageGet(message core.CMessage, uid int)
	Communication_MessageGetAll(messages []map[string]string, getSent bool, count int, page int)
	Essential_GetMusic(core.CMusic)
	Essential_GetTopArtists(artists map[string]string)
	Level_GetGauntlets(gaus []map[string]string, hash string)
	Level_SearchList(intlists []int, lists []core.CLevelList, count int, page int)
	Level_GetMapPacks(packs []core.LevelPack, count int, page int)
	Level_GetLevelFull(lvl core.CLevel, passwd string, phash string, quest_id int)
	Level_GetSpecials(id int, timeLeft int)
	Level_SearchLevels(intlevels []int, levels []core.CLevel, mus *core.CMusic, count int, page int, gdVersion int, gauntlet int)
	Rewards_ChallengesOutput(cq core.CQuests, uid int, chk string, udid string)
	Rewards_ChestOutput(acc core.CAccount, config core.ConfigBlob, udid string, chk string, smallLeft int, bigLeft int, chestType int)
	Profile_GetUserProfile(acc core.CAccount, selfUid int)
	Profile_ListUserProfiles(accs []core.CAccount)
	Profile_GetSearchableUsers(accs []core.CAccount, count int, page int)
	Score_GetLeaderboard(intaccs []int, xacc core.CAccount)
	Score_GetScores(scores []core.CScores, mode string)
}

func NewConnector(isJson bool) Connector {
	if isJson {
		return &JSONConnector{output: make(map[string]interface{})}
	} else {
		return &GDConnector{}
	}
}
