package core

import (
	"encoding/json"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"os"
)

func RunSingleTask(Srvid string, rdb RedisConn, log Logger, config GlobalConfig) {
	t,err:=rdb.DB.Get(rdb.context,Srvid).Result()
	if err!=nil {
		log.LogWarn(rdb,err.Error())
		return
	}
	fmt.Println(Srvid)
	conf:=ConfigBlob{}
	err=json.Unmarshal([]byte(t),&conf)
	if err!=nil{
		log.LogWarn(rdb,err.Error())
		return
	}
	db:=MySQLConn{}
	defer db.CloseDB()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Failed. Dequeuing...")
		}

	}()
	if log.Should(db.ConnectBlob(conf))!=nil {return}
	//Start real stuff
	os.MkdirAll(config.SavePath+"/"+Srvid+"/savedata",0777)
	mus:=CMusic{DB: db}
	mus.CountDownloads()
	protect:=CProtect{DB: db, Savepath: config.SavePath+"/"+Srvid}
	protect.FillLevelModel()
	protect.ResetUserLimits()
}

func MaintainTasks(config GlobalConfig) {
	rdb:=RedisConn{}
	log:=Logger{}
	if err:=rdb.ConnectBlob(config); err!=nil {
		log.LogWarn(rdb,err.Error())
		return
	}
	strsl,err:=rdb.DB.Keys(rdb.context,"*").Result()
	if err!=nil {
		log.LogWarn(rdb,err.Error())
		return
	}

	for _,Srvid := range strsl {
		RunSingleTask(Srvid, rdb, log, config)
	}

}

func MaintainRoutines(config GlobalConfig) {
	gocron.Every(1).Day().At("00:00").Do(MaintainTasks, config)
}

