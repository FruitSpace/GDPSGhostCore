package core

import "strings"

type CAccountJSON struct {
	Uid             int    `json:"uid"`
	Uname           string `json:"uname"`
	Email           string `json:"email,omitempty"`
	ModBadge        int    `json:"mod_badge,omitempty"`
	Role            Role   `json:"role,omitempty"`
	IsBanned        bool   `json:"is_banned"`
	LeaderboardRank int    `json:"leaderboard_rank,omitempty"`

	WeAreFriends      bool `json:"we_are_friends,omitempty"`
	NewMessagesCount  int  `json:"new_messages_count,omitempty"`
	NewFriendRequests int  `json:"new_friend_requests_count,omitempty"`

	Stats     CAccountStats     `json:"stats,omitempty"`
	Technical CAccountTechnical `json:"technical,omitempty"`
	Social    CAccountSocial    `json:"social,omitempty"`
	Vessels   CAccountVessels   `json:"vessels,omitempty"`
	Chests    CAccountChests    `json:"chests,omitempty"`
	Settings  CAccountSettings  `json:"settings,omitempty"`
}

func NewCAccountJSONFromAccount(acc CAccount, role *Role, personal bool) CAccountJSON {
	// Basic stuff
	aj := CAccountJSON{
		Uid:             acc.Uid,
		Uname:           acc.Uname,
		IsBanned:        acc.IsBanned > 0,
		LeaderboardRank: acc.GetLeaderboardRank() + 1,
		Stats:           NewCAccountStatsFromAccount(acc),
		Vessels:         NewCAccountVesselsFromAccount(acc),
		Settings:        NewCAccountSettingsFromAccount(acc),
	}

	// Roles
	if role != nil {
		aj.ModBadge = role.ModLevel
		if personal {
			aj.Role = *role
		}
	}

	// Personal
	if personal {
		aj.Email = acc.Email
		aj.Social = NewCAccountSocialFromAccount(acc)
		aj.Technical = NewCAccountTechnicalFromAccount(acc)
		aj.Chests = NewCAccountChestsFromAccount(acc)
		cf := CFriendship{DB: acc.DB}
		cm := CMessage{DB: acc.DB}
		aj.NewFriendRequests = cf.CountFriendRequests(acc.Uid, true)
		aj.NewMessagesCount = cm.CountMessages(acc.Uid, true)
	}

	return aj
}

func NewCAccountJSONFromAccountLite(acc CAccount) CAccountJSON {
	return CAccountJSON{
		Uid:      acc.Uid,
		Uname:    acc.Uname,
		IsBanned: acc.IsBanned > 0,
		Stats:    NewCAccountStatsFromAccount(acc),
		Vessels:  NewCAccountVesselsFromAccount(acc),
	}
}

type CAccountStats struct {
	Stars         int `json:"stars"`
	Diamonds      int `json:"diamonds"`
	Coins         int `json:"coins"`
	UCoins        int `json:"ucoins"`
	Demons        int `json:"demons"`
	CPoints       int `json:"cpoints"`
	Orbs          int `json:"orbs"`
	Moons         int `json:"moons"`
	Special       int `json:"special"`
	LvlsCompleted int `json:"lvls_completed"`
}

func NewCAccountStatsFromAccount(acc CAccount) CAccountStats {
	return CAccountStats{
		Stars:         acc.Stars,
		Diamonds:      acc.Diamonds,
		Coins:         acc.Coins,
		UCoins:        acc.UCoins,
		Demons:        acc.Demons,
		CPoints:       acc.CPoints,
		Orbs:          acc.Orbs,
		Moons:         acc.Moons,
		Special:       acc.Special,
		LvlsCompleted: acc.LvlsCompleted,
	}
}

type CAccountTechnical struct {
	RegDate    string `json:"reg_date"`
	AccessDate string `json:"access_date"`
	LastIP     string `json:"last_ip"`
	GameVer    string `json:"game_ver"`
}

func NewCAccountTechnicalFromAccount(acc CAccount) CAccountTechnical {
	return CAccountTechnical{
		RegDate:    acc.RegDate,
		AccessDate: acc.AccessDate,
		LastIP:     acc.LastIP,
		GameVer:    acc.GameVer,
	}
}

type CAccountSocial struct {
	BlacklistIds  []int `json:"blacklist_ids"`
	FriendsCount  int   `json:"friends_count"`
	FriendshipIds []int `json:"friendship_ids"`
}

func NewCAccountSocialFromAccount(acc CAccount) CAccountSocial {
	blacklist := strings.Split(acc.Blacklist, ",")
	friendships := strings.Split(acc.FriendshipIds, ",")
	return CAccountSocial{
		BlacklistIds:  ArrTranslateToInt(blacklist),
		FriendsCount:  acc.FriendsCount,
		FriendshipIds: ArrTranslateToInt(friendships),
	}
}

type CAccountVessels struct {
	ShownIcon      int `json:"shown_icon"`
	IconType       int `json:"icon_type"`
	ColorPrimary   int `json:"color_primary"`
	ColorSecondary int `json:"color_secondary"`
	ColorGlow      int `json:"color_glow"`
	Cube           int `json:"cube"`
	Ship           int `json:"ship"`
	Ball           int `json:"ball"`
	Ufo            int `json:"ufo"`
	Wave           int `json:"wave"`
	Robot          int `json:"robot"`
	Spider         int `json:"spider"`
	Swing          int `json:"swing"`
	Jetpack        int `json:"jetpack"`
	Trace          int `json:"trace"`
	Death          int `json:"death"`
}

func NewCAccountVesselsFromAccount(acc CAccount) CAccountVessels {
	return CAccountVessels{
		ShownIcon:      acc.GetShownIcon(),
		IconType:       acc.IconType,
		ColorPrimary:   acc.ColorPrimary,
		ColorSecondary: acc.ColorSecondary,
		ColorGlow:      acc.ColorGlow,
		Cube:           acc.Cube,
		Ship:           acc.Ship,
		Ball:           acc.Ball,
		Ufo:            acc.Ufo,
		Wave:           acc.Wave,
		Robot:          acc.Robot,
		Spider:         acc.Spider,
		Swing:          acc.Swing,
		Jetpack:        acc.Jetpack,
		Trace:          acc.Trace,
		Death:          acc.Death,
	}
}

type CAccountChests struct {
	ChestSmallCount int `json:"chest_small_count"`
	ChestSmallTime  int `json:"chest_small_time_left"`
	ChestBigCount   int `json:"chest_big_count"`
	ChestBigTime    int `json:"chest_big_time_left"`
}

func NewCAccountChestsFromAccount(acc CAccount) CAccountChests {
	return CAccountChests{
		ChestSmallCount: acc.ChestSmallCount,
		ChestSmallTime:  acc.ChestSmallTime,
		ChestBigCount:   acc.ChestBigCount,
		ChestBigTime:    acc.ChestBigTime,
	}
}

type CAccountSettings struct {
	AllowFriendReq bool   `json:"allow_friend_requests"`
	AllowComments  string `json:"allow_view_comments"`
	AllowMessages  string `json:"allow_messages"`
	Youtube        string `json:"youtube"`
	Twitch         string `json:"twitch"`
	Twitter        string `json:"twitter"`
}

func NewCAccountSettingsFromAccount(acc CAccount) CAccountSettings {
	s := CAccountSettings{
		Youtube:        acc.Youtube,
		Twitch:         acc.Twitch,
		Twitter:        acc.Twitter,
		AllowFriendReq: true,
		AllowComments:  "everybody",
		AllowMessages:  "everybody",
	}
	allowances := []string{"everybody", "friends", "nobody"}

	if acc.FrS > 0 {
		s.AllowFriendReq = false
	}
	if acc.CS > 0 {
		s.AllowComments = allowances[Clamp(acc.CS, 0, 2)]
	}
	if acc.MS > 0 {
		s.AllowMessages = allowances[Clamp(acc.MS, 0, 2)]
	}
	return s
}
