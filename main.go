//+build linux,amd64 windows,amd64

package main

import (
	"HalogenGhostCore/api"
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	"encoding/json"
	"fmt"
)

func main() {
	//config:=GenConfig()
	//gconfig:=GenGConfig()

	//logger:=core.Logger{
	//	connectors.GetWriter("",""),
	//}
	//db:=core.MySQLConn{}
	//logger.Must(db,db.ConnectBlob(config))
	//
	//cc:=core.ThunderAES{}
	//logger.Must(cc,cc.GenKey("Amongus"))
	//logger.Must(cc,cc.Init())
	//fmt.Println(logger.Must(cc.Encrypt("Artyom is Gay")))
	//
	//acc:= core.CAccount{DB: db, Logger: logger}
	//acc.Uid=acc.SearchUsers("Nanorun")
	//acc.LoadAll()
	//
	//PrintCAccount(acc)
	//
	x:=core.ConfigBlob{}
	err:=json.Unmarshal([]byte(`{"DBConfig":{"Host":"localhost","Port":3306,"User":"halogen","Password":"D0wn_Th3_r4BB1t_H0lE_731","DBName":"gdps_0002"},"LogConfig":{"LogEnable":true,"LogDB":false,"LogEndpoints":false,"LogRequests":false},"ChestConfig":{"ChestSmallOrbsMin":200,"ChestSmallOrbsMax":400,"ChestSmallDiamondsMin":2,"ChestSmallDiamondsMax":10,"ChestSmallShardsMin":1,"ChestSmallShardsMax":6,"ChestSmallKeysMin":1,"ChestSmallKeysMax":6,"ChestSmallWait":3600,"ChestBigOrbsMin":2000,"ChestBigOrbsMax":4000,"ChestBigDiamondsMin":20,"ChestBigDiamondsMax":100,"ChestBigShardsMin":1,"ChestBigShardsMax":6,"ChestBigKeysMin":1,"ChestBigKeysMax":6,"ChestBigWait":14400},"ServerConfig":{"SrvID":"0002","SrvKey":"SRV_KEY","MaxUsers":100,"MaxLevels":500,"MaxComments":1000,"MaxPosts":1000,"HalMusic":true,"Locked":false}} `),&x)
	fmt.Println(x,err)
	ghostServer:= api.GhostServer{
		Log: core.Logger{
			Output: connectors.GetWriter("",""),
		},
		Config: GenGConfig(),
	}
	ghostServer.StartServer("0.0.0.0:1997")

}

func GenConfig() core.ConfigBlob {
	return core.ConfigBlob{
		DBConfig: core.MysqlConfig{
			"localhost",
			3306,
			"halogen",
			"D0wn_Th3_r4BB1t_H0lE_731",
			"gdps_0002",
		},
		LogConfig: core.LogConfig{
			true,
			false,
			false,
			false,
		},
		ChestConfig: core.ChestConfig{
			200,400,2,10,
			1,6,1,6, 3600,
			2000,4000,20,100,
			1,6,1,6, 14400,
		},
		ServerConfig: core.ServerConfig{
			"0002",
			"SRV_KEY",
			100,
			500,
			1000,
			1000,
			true,
			false,
		},
	}
}

func GenGConfig() core.GlobalConfig {
	return core.GlobalConfig{
		"https://halhost.cc/app/api/gdps_api.php",
		"stdout",
		"null",
		false,
		"localhost",
		"6379",
		"",
		7,
		"./",

		map[string]string{},
	}
}

func PrintCAccount(acc core.CAccount){
	fmt.Printf("[%d] %s (%s) [Role:%d] %d\nPass: %s\nS:%d | D:%d | C:%d | UC:%d | CP:%d | O:%d\nDemons: %d | Special: %d | Lvls: %d\nReg: %s | Acc: %s | IP: %s | Ver: %s\n",
		acc.Uid,acc.Uname,acc.Email,acc.Role_id,acc.IsBanned,acc.Passhash,acc.Stars,acc.Diamonds,acc.Coins,acc.UCoins,acc.CPoints,acc.Orbs,acc.Demons,acc.Special,acc.LvlsCompleted,
		acc.RegDate,acc.AccessDate,acc.LastIP,acc.GameVer)
	fmt.Printf("Blacklist: %s | FrCnt: %d | Friends: %s\nfrS: %d | cS: %d | mS: %d | YT: %s | Tch: %s | Twr: %s\n",
		acc.Blacklist,acc.FriendsCount,acc.FriendshipIds,acc.FrS,acc.CS,acc.MS,acc.Youtube,acc.Twitch,acc.Twitter)
	fmt.Printf("Icon: %d | Clr1: %d | Clr2: %d | Trace: %d | Death: %d\nCube:%d | Ship: %d | Ball: %d | UFO: %d | Wave: %d | Robot: %d | Spider: %d\n",
		acc.IconType,acc.ColorPrimary,acc.ColorSecondary,acc.Trace,acc.Death,acc.Cube,acc.Ship,acc.Ball,acc.Ufo,acc.Wave,acc.Robot,acc.Spider)
}
