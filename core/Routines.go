package core

import (
	"encoding/json"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/jasonlvhit/gocron"
	"log"
	"os"
	"strconv"
)

var LEADER = false
var SessionID string
var KvEngine *consul.KV

func RunSingleTask(Srvid string, rdb RedisConn, log Logger, config GlobalConfig) {
	t, err := rdb.DB.Get(rdb.context, Srvid).Result()
	if err != nil {
		log.LogWarn(rdb, err.Error())
		return
	}
	conf := ConfigBlob{}
	err = json.Unmarshal([]byte(t), &conf)
	if err != nil {
		log.LogWarn(rdb, err.Error())
		return
	}
	db := MySQLConn{}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Failed. Dequeuing...")
		}

	}()
	if log.Should(db.ConnectBlob(conf)) != nil {
		return
	}
	//Start real stuff
	os.MkdirAll(config.SavePath+"/"+Srvid+"/savedata", 0777)
	mus := CMusic{DB: &db}
	mus.CountDownloads()
	protect := CProtect{DB: &db, Savepath: config.SavePath + "/" + Srvid}
	protect.FillLevelModel()
	protect.ResetUserLimits()
}

func MaintainTasks(config GlobalConfig) {
	rdb := RedisConn{}
	log := Logger{}
	if err := rdb.ConnectBlob(config); err != nil {
		log.LogWarn(rdb, err.Error())
		return
	}
	strsl, err := rdb.DB.Keys(rdb.context, "*").Result()
	if err != nil {
		log.LogWarn(rdb, err.Error())
		return
	}
	SendMessageDiscord("Starting maintenance routine")
	for i, Srvid := range strsl {
		fmt.Println("["+strconv.Itoa(i+1)+"/"+strconv.Itoa(len(strsl))+"]", Srvid)
		RunSingleTask(Srvid, rdb, log, config)
	}

}

func PrepareElection(config GlobalConfig) {
	consulConf := consul.DefaultConfig()
	consulConf.Address = EnvOrDefault("CONSUL_ADDR", "127.0.0.1")
	consulConf.Token = EnvOrDefault("CONSUL_TOKEN", "")
	consulConf.Datacenter = "hal"
	consulCli, err := consul.NewClient(consulConf)
	if err != nil {
		log.Println("Unable to connect to Consul cluster. Assuming self-leadership: " + err.Error())
		LEADER = true
	} else {
		sessEngine := consulCli.Session()
		KvEngine = consulCli.KV()
		SessID, _, err := sessEngine.Create(&consul.SessionEntry{Name: "GhostCore"}, nil)
		if err != nil {
			log.Println("Unable to connect to create Session. Assuming self-leadership: " + err.Error())
			LEADER = true
		} else {
			SessionID = SessID
			AquireLeadership()
			if !LEADER {
				log.Println("Couldn't acquire leadership. Dispatching 10min watchdog")
				if err = gocron.Every(10).Seconds().Do(AquireLeadership); err != nil {
					log.Println(err)
				}
			}
		}
	}

	if LEADER {
		gocron.Every(1).Day().At("03:00").Do(MaintainTasks, config)
	}

	go gocron.Start()

}

func AquireLeadership() {
	kvData := &consul.KVPair{
		Key:     "sessions/ghostcore_lead",
		Value:   []byte(EnvOrDefault("NOMAD_SHORT_ALLOC_ID", "default")),
		Session: SessionID,
	}
	isAcq, _, err := KvEngine.Acquire(kvData, nil)
	if err == nil && isAcq {
		log.Println("Lock was successfully acquired. NOW LEADER")
		LEADER = true
	} else {
		log.Println("Couldn't acquire leadership. Still a follower.")
	}
}

func StepDown() {
	kvData := &consul.KVPair{
		Key:     "sessions/ghostcore_lead",
		Value:   []byte(EnvOrDefault("NOMAD_SHORT_ALLOC_ID", "default")),
		Session: SessionID,
	}
	isRel, _, err := KvEngine.Release(kvData, nil)
	if err == nil && isRel {
		log.Println("Lock was successfully released. NOW FOLLOWER")
		LEADER = false
	} else {
		log.Println("[!!!] COULD NOT RELEASE LOCK [!!!]")
	}
}

func EnvOrDefault(key string, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
