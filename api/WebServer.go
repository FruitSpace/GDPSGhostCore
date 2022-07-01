package api

import (
	"HalogenGhostCore/core"
	"log"
	"net/http"
	gorilla "github.com/gorilla/mux"
)


type GhostServer struct {
	log core.Logger
	db core.MySQLConn
	rdb core.RedisConn
	config core.ConfigBlob
}

var RouteMap = map[string]func(resp http.ResponseWriter, req *http.Request){
	"/": Shield,

	"/database/accounts/accountManagement.php": AccountManagement,
	"/database/accounts/backupGJAccount": AccountBackup,
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


func (ghost *GhostServer) StartServer(Host string) {
	mux:=gorilla.NewRouter()
	mux.HandleFunc("/",Redirector)
	for route,handler:= range RouteMap {
		mux.HandleFunc("/{gdps:[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]}"+route,handler)
	}
	log.Println("Server is up and running on http://"+Host)
	err:=http.ListenAndServe(Host,mux)
	if err!=nil {
		ghost.log.LogErr(ghost,err.Error())
	}

}