package core

import (
	"encoding/json"
	"github.com/jasonlvhit/gocron"
)

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
		t,err:=rdb.DB.Get(rdb.context,Srvid).Result()
		if err!=nil {
			log.LogWarn(rdb,err.Error())
			continue
		}
		conf:=ConfigBlob{}
		err=json.Unmarshal([]byte(t),&conf)
		if err!=nil{
			log.LogWarn(rdb,err.Error())
			continue
		}
		db:=MySQLConn{}
		if log.Should(db.ConnectBlob(conf))!=nil {continue}
		//Start real stuff
		mus:=CMusic{DB: db}
		mus.CountDownloads()
		protect:=CProtect{DB: db, Savepath: config.SavePath+"/"+Srvid+"/levelModel.json"}
		protect.FillLevelModel()
		protect.ResetUserLimits()
	}

}

func MaintainRoutines(config GlobalConfig) {
	gocron.Every(1).Day().At("00:00").Do(MaintainTasks, config)
}

