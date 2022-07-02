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
	Log core.Logger
	Config core.GlobalConfig
}

var RouteMap = map[string]func(http.ResponseWriter, *http.Request, *core.GlobalConfig){
	"/": Shield,

	"/database/accounts/accountManagement.php": AccountManagement,
	"/database/accounts/backupGJAccount.php": AccountBackup,
	"/database/accounts/loginGJAccount.php": AccountLogin,
	"/database/accounts/registerGJAccount.php": AccountRegister,
	"/database/accounts/syncGJAccount.php": AccountSync,
	"/database/accounts/syncGJAccount20.php": AccountSync,

	"/database/database/accounts/backupGJAccountNew.php": AccountBackup,
	"/database/database/accounts/syncGJAccountNew.php": AccountSync,

	"/database/acceptGJFriendRequest20.php": FriendAcceptRequest,
	"/database/blockGJUser20.php": BlockUser,
	"/database/deleteGJAccComment20.php": AccountCommentDelete,
	"/database/deleteGJComment20.php": CommentDelete,
	"/database/deleteGJFriendRequests20.php": FriendRejectRequest,
	"/database/deleteGJLevelUser20.php": LevelDelete,
	"/database/deleteGJMessages20.php": MessageDelete,
	"/database/downloadGJLevel.php": LevelDownload,
	"/database/downloadGJLevel19.php": LevelDownload,
	"/database/downloadGJLevel20.php": LevelDownload,
	"/database/downloadGJLevel21.php": LevelDownload,
	"/database/downloadGJLevel22.php": LevelDownload,
	"/database/downloadGJMessage20.php": MessageGet,
	"/database/getAccountURL.php": GetAccountUrl,
	"/database/getGJAccountComments20.php": AccountCommentGet,
	"/database/getGJChallenges.php": GetChallenges,
	"/database/getGJCommentHistory.php": CommentGetHistory,
	"/database/getGJComments.php": CommentGet,
	"/database/getGJComments19.php": CommentGet,
	"/database/getGJComments20.php": CommentGet,
	"/database/getGJComments21.php": CommentGet,
	"/database/getGJCreators.php": GetCreators,
	"/database/getGJCreators19.php": GetCreators,
	"/database/getGJDailyLevel.php": LevelGetDaily,
	"/database/getGJFriendRequests20.php": FriendGetRequests,
	"/database/getGJGauntlets.php": GetGauntlets,
	"/database/getGJGauntlets21.php": GetGauntlets,
	"/database/getGJLevels.php": LevelGetLevels,
	"/database/getGJLevels19.php": LevelGetLevels,
	"/database/getGJLevels20.php": LevelGetLevels,
	"/database/getGJLevels21.php": LevelGetLevels,
	"/database/getGJLevelScores.php": GetLevelScores,
	"/database/getGJLevelScores211.php": GetLevelScores,
	"/database/getGJMapPacks.php": GetMapPacks,
	"/database/getGJMapPacks20.php": GetMapPacks,
	"/database/getGJMapPacks21.php": GetMapPacks,
	"/database/getGJMessages20.php": MessageGet,
	"/database/getGJRewards.php": GetRewards,
	"/database/getGJScores.php": GetScores,
	"/database/getGJScores19.php": GetScores,
	"/database/getGJScores20.php": GetScores,
	"/database/getGJSongInfo.php": GetSongInfo,
	"/database/getGJTopArtists.php": GetTopArtists,
	"/database/getGJUserInfo20.php": GetUserInfo,
	"/database/getGJUserList20.php": GetUserList,
	"/database/getGJUsers20.php": GetUsers,
	"/database/likeGJItem.php": LikeItem,
	"/database/likeGJItem19.php": LikeItem,
	"/database/likeGJItem20.php": LikeItem,
	"/database/likeGJItem21.php": LikeItem,
	"/database/likeGJItem211.php": LikeItem,
	"/database/rateGJDemon21.php": RateDemon,
	"/database/rateGJStars20.php": RateStar,
	"/database/rateGJStars211.php": RateStar,
	"/database/readGJFriendRequest20.php": FriendReadRequest,
	"/database/removeGJFriend20.php": FriendRemove,
	"/database/reportGJLevel.php": LevelReport,
	"/database/requestUserAccess.php": RequestMod,
	"/database/suggestGJStars20.php": SuggestStars,
	"/database/unblockGJUser20.php": UnblockUser,
	"/database/updateGJAccSettings20.php": UpdateAccountSettings,
	"/database/updateGJDesc20.php": LevelUpdateDescription,
	"/database/updateGJUserScore.php": UpdateUserScore,
	"/database/updateGJUserScore19.php": UpdateUserScore,
	"/database/updateGJUserScore20.php": UpdateUserScore,
	"/database/updateGJUserScore21.php": UpdateUserScore,
	"/database/updateGJUserScore22.php": UpdateUserScore,
	"/database/uploadFriendRequest20.php": FriendRequest,
	"/database/uploadGJAccComment20.php": AccountCommentUpload,
	"/database/uploadGJComment.php": CommentUpload,
	"/database/uploadGJComment19.php": CommentUpload,
	"/database/uploadGJComment20.php": CommentUpload,
	"/database/uploadGJComment21.php": CommentUpload,
	"/database/uploadGJLevel.php": LevelUpload,
	"/database/uploadGJLevel19.php": LevelUpload,
	"/database/uploadGJLevel20.php": LevelUpload,
	"/database/uploadGJLevel21.php": LevelUpload,
	"/database/uploadGJMessage20.php": MessageUpload,


}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (ghost *GhostServer) StartServer(Host string) {
	mux:=gorilla.NewRouter()
	var nfh NotFoundHandler
	mux.NotFoundHandler=nfh
	mux.HandleFunc("/",Redirector)
	for route:= range RouteMap {
		mux.HandleFunc("/{gdps:[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]}"+route,
			func(resp http.ResponseWriter,req *http.Request){
				vars:=gorilla.Vars(req)
				handler:=RouteMap[strings.Replace(req.URL.Path,"/"+vars["gdps"],"",1)]
				log.Println(req.URL.Path," Got ", GetFunctionName(handler))
				handler(resp,req,&ghost.Config)
			})
	}
	log.Println("Server is up and running on http://"+Host)
	err:=http.ListenAndServe(Host,mux)
	if err!=nil {
		ghost.Log.LogErr(ghost,err.Error())
	}

}

func ReadPost(req *http.Request) url.Values {
	if req.Body==nil { return url.Values{}}
	body,err:=io.ReadAll(req.Body)
	if err!=nil {
		log.Println(err.Error())
		return url.Values{}
	}
	vals:=make(url.Values)
	pairs:=strings.Split(string(body),"&")
	for _,val:= range pairs {

		m:=strings.SplitN(val,"=",2)
		fmt.Println(m)
		rval,_:=url.QueryUnescape(m[1])
		rkey,_:=url.QueryUnescape(m[0])
		vals[rkey]=append(vals[rkey],rval)
	}
	return vals

}