package api

import (
	"HalogenGhostCore/core"
	"fmt"
	gorilla "github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"
)

type GhostServer struct {
	Log    core.Logger
	Config core.GlobalConfig
}

var RouteMap = map[string]func(http.ResponseWriter, *http.Request, *core.GlobalConfig){
	"/": Shield,

	// Geometry Dash
	"/db/accounts/accountManagement.php": AccountManagement,
	"/db/accounts/backupGJAccount.php":   AccountBackup,
	"/db/accounts/loginGJAccount.php":    AccountLogin,
	"/db/accounts/registerGJAccount.php": AccountRegister,
	"/db/accounts/syncGJAccount.php":     AccountSync,
	"/db/accounts/syncGJAccount20.php":   AccountSync,

	"/db/database/accounts/backupGJAccountNew.php": AccountBackup,
	"/db/database/accounts/syncGJAccountNew.php":   AccountSync,

	"/db/acceptGJFriendRequest20.php":  FriendAcceptRequest,
	"/db/blockGJUser20.php":            BlockUser,
	"/db/deleteGJAccComment20.php":     AccountCommentDelete,
	"/db/deleteGJComment20.php":        CommentDelete,
	"/db/deleteGJFriendRequests20.php": FriendRejectRequest,
	"/db/deleteGJLevelUser20.php":      LevelDelete,
	"/db/deleteGJMessages20.php":       MessageDelete,
	"/db/downloadGJLevel.php":          LevelDownload,
	"/db/downloadGJLevel19.php":        LevelDownload,
	"/db/downloadGJLevel20.php":        LevelDownload,
	"/db/downloadGJLevel21.php":        LevelDownload,
	"/db/downloadGJLevel22.php":        LevelDownload,
	"/db/downloadGJMessage20.php":      MessageGet,
	"/db/getAccountURL.php":            GetAccountUrl,
	"/db/getGJAccountComments20.php":   AccountCommentGet,
	"/db/getGJChallenges.php":          GetChallenges,
	"/db/getGJCommentHistory.php":      CommentGetHistory,
	"/db/getGJComments.php":            CommentGet,
	"/db/getGJComments19.php":          CommentGet,
	"/db/getGJComments20.php":          CommentGet,
	"/db/getGJComments21.php":          CommentGet,
	"/db/getGJCreators.php":            GetCreators,
	"/db/getGJCreators19.php":          GetCreators,
	"/db/getGJDailyLevel.php":          LevelGetDaily,
	"/db/getGJFriendRequests20.php":    FriendGetRequests,
	"/db/getGJGauntlets.php":           GetGauntlets,
	"/db/getGJGauntlets21.php":         GetGauntlets,
	"/db/getGJLevels.php":              LevelGetLevels,
	"/db/getGJLevels19.php":            LevelGetLevels,
	"/db/getGJLevels20.php":            LevelGetLevels,
	"/db/getGJLevels21.php":            LevelGetLevels,
	"/db/getGJLevelScores.php":         GetLevelScores,
	"/db/getGJLevelScores211.php":      GetLevelScores,
	"/db/getGJMapPacks.php":            GetMapPacks,
	"/db/getGJMapPacks20.php":          GetMapPacks,
	"/db/getGJMapPacks21.php":          GetMapPacks,
	"/db/getGJMessages20.php":          MessageGetAll,
	"/db/getGJRewards.php":             GetRewards,
	"/db/getGJScores.php":              GetScores,
	"/db/getGJScores19.php":            GetScores,
	"/db/getGJScores20.php":            GetScores,
	"/db/getGJSongInfo.php":            GetSongInfo,
	"/db/getGJTopArtists.php":          GetTopArtists,
	"/db/getGJUserInfo20.php":          GetUserInfo,
	"/db/getGJUserList20.php":          GetUserList,
	"/db/getGJUsers20.php":             GetUsers,
	"/db/likeGJItem.php":               LikeItem,
	"/db/likeGJItem19.php":             LikeItem,
	"/db/likeGJItem20.php":             LikeItem,
	"/db/likeGJItem21.php":             LikeItem,
	"/db/likeGJItem211.php":            LikeItem,
	"/db/rateGJDemon21.php":            RateDemon,
	"/db/rateGJStars20.php":            RateStar,
	"/db/rateGJStars211.php":           RateStar,
	"/db/readGJFriendRequest20.php":    FriendReadRequest,
	"/db/removeGJFriend20.php":         FriendRemove,
	"/db/reportGJLevel.php":            LevelReport,
	"/db/requestUserAccess.php":        RequestMod,
	"/db/suggestGJStars20.php":         SuggestStars,
	"/db/unblockGJUser20.php":          UnblockUser,
	"/db/updateGJAccSettings20.php":    UpdateAccountSettings,
	"/db/updateGJDesc20.php":           LevelUpdateDescription,
	"/db/updateGJUserScore.php":        UpdateUserScore,
	"/db/updateGJUserScore19.php":      UpdateUserScore,
	"/db/updateGJUserScore20.php":      UpdateUserScore,
	"/db/updateGJUserScore21.php":      UpdateUserScore,
	"/db/updateGJUserScore22.php":      UpdateUserScore,
	"/db/uploadFriendRequest20.php":    FriendRequest,
	"/db/uploadGJAccComment20.php":     AccountCommentUpload,
	"/db/uploadGJComment.php":          CommentUpload,
	"/db/uploadGJComment19.php":        CommentUpload,
	"/db/uploadGJComment20.php":        CommentUpload,
	"/db/uploadGJComment21.php":        CommentUpload,
	"/db/uploadGJLevel.php":            LevelUpload,
	"/db/uploadGJLevel19.php":          LevelUpload,
	"/db/uploadGJLevel20.php":          LevelUpload,
	"/db/uploadGJLevel21.php":          LevelUpload,
	"/db/uploadGJMessage20.php":        MessageUpload,

	"/db/getCustomContentURL.php": GetContentURL,
	"/db/content/sfx/{sfxid}":     RelaySFX,

	//"/db/content/sfx/sfxlibrary.dat":         GetSFXLibrary,
	//"/db/content/sfx/sfxlibrary_version.txt": GetSFXLibraryVersion,
	//"/db/content/sfx/s{sfxid}.ogg":           GetSFXTrack,
}

var RouteIntegraMap = map[string]func(http.ResponseWriter, *http.Request, *core.GlobalConfig){
	// PRIVATE API
	"/integra/maintenance": TriggerMaintenance,
	//"/integra/killskew":    EventAction,
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (ghost *GhostServer) StartServer(Host string) {
	BallisticsCache = make(map[string]int64)
	BadRepIP = make(map[string]int)
	mux := gorilla.NewRouter()
	var nfh NotFoundHandler
	mux.NotFoundHandler = nfh
	mux.HandleFunc("/", Redirector)
	for route := range RouteMap {
		mux.HandleFunc("/{gdps:[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]}"+route,
			func(resp http.ResponseWriter, req *http.Request) {
				vars := gorilla.Vars(req)
				pref := strings.Replace(req.URL.Path, "/"+vars["gdps"], "", 1)
				handler := RouteMap[pref]
				if strings.HasPrefix(pref, "/db/content/sfx/") {
					handler = RelaySFX
				}
				IPAddr := req.Header.Get("CF-Connecting-IP")
				if IPAddr == "" {
					IPAddr = req.Header.Get("X-Real-IP")
				}
				if IPAddr == "" {
					IPAddr = strings.Split(req.RemoteAddr, ":")[0]
				}
				log.Println("["+IPAddr+"] "+req.URL.Path, " Got ", GetFunctionName(handler))
				handler(resp, req, &ghost.Config)
			})
	}
	for route, handler := range RouteIntegraMap {
		mux.HandleFunc(route,
			func(resp http.ResponseWriter, req *http.Request) {
				defer func() {
					if err := recover(); err != nil {
						log.Println(err)
						core.SendMessageDiscord(fmt.Sprintf("Panic: %s", err))
					}
				}()
				IPAddr := req.Header.Get("CF-Connecting-IP")
				if IPAddr == "" {
					IPAddr = req.Header.Get("X-Real-IP")
				}
				if IPAddr == "" {
					IPAddr = strings.Split(req.RemoteAddr, ":")[0]
				}
				log.Println("["+IPAddr+"] "+req.URL.Path, " Got ", GetFunctionName(handler))
				handler(resp, req, &ghost.Config)
			})
	}
	log.Println("Server is up and running on http://" + Host)
	err := http.ListenAndServe(Host, mux)
	if err != nil {
		ghost.Log.LogErr(ghost, err.Error())
	}

}

func ReadPost(req *http.Request) url.Values {
	if req.Body == nil {
		return url.Values{}
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err.Error())
		return url.Values{}
	}
	if len(body) == 0 || strings.Count(string(body), "=") == 0 {
		return url.Values{}
	}
	vals := make(url.Values)
	pairs := strings.Split(string(body), "&")
	for _, val := range pairs {
		if !strings.Contains(val, "=") {
			continue
		}
		m := strings.SplitN(val, "=", 2)
		//fmt.Println(m)
		rval, _ := url.QueryUnescape(m[1])
		rkey, _ := url.QueryUnescape(m[0])
		vals[rkey] = append(vals[rkey], rval)
	}
	return vals
}
