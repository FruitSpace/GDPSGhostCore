package core

import (
	"encoding/json"
	"fmt"
	"github.com/go-co-op/gocron"
	consul "github.com/hashicorp/consul/api"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

var ucron = gocron.NewScheduler(loc)

var LEADER = false
var LEAD_CONFIG GlobalConfig
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
			SendMessageDiscord("[" + Srvid + "] Failed. Dequeuing...")
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
	protect.ResetUserLimits()
	protect.FillLevelModel()

}

func CleanModels() {
	dir, err := os.ReadDir(LEAD_CONFIG.SavePath)
	if err != nil {
		return
	}
	for _, p := range dir {
		os.RemoveAll(path.Join(LEAD_CONFIG.SavePath, p.Name()))
	}
}

func MaintainTasks() {
	if !LEADER {
		CleanModels()
		return
	}
	config := LEAD_CONFIG
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
	SendMessageDiscord("Starting Maintenance Routine by `" + EnvOrDefault("NOMAD_SHORT_ALLOC_ID", "default") + "`")
	for i, SrvId := range strsl {
		fmt.Println("["+strconv.Itoa(i+1)+"/"+strconv.Itoa(len(strsl))+"]", SrvId)
		RunSingleTask(SrvId, rdb, log, config)
	}

	// Update sfx library
	updateSFXLibrary()

}

func updateSFXLibrary() {
	s3 := NewS3FS()
	d, err := http.Get("https://geometrydashfiles.b-cdn.net/sfx/sfxlibrary.dat")
	if err != nil {
		SendMessageDiscord("⚠️ Failed to fetch SFX Library: " + err.Error())
		return
	}
	sfx, _ := io.ReadAll(d.Body)
	err = s3.PutFile("/gdps_sfx/library.dat", sfx)
	if err != nil {
		SendMessageDiscord("⚠️ Failed to save SFX Library: " + err.Error())
		return
	}

	v, err := http.Get("https://geometrydashfiles.b-cdn.net/sfx/sfxlibrary_version.txt")
	if err != nil {
		SendMessageDiscord("⚠️ Failed to fetch SFX Library Version: " + err.Error())
		return
	}
	ver, _ := io.ReadAll(v.Body)
	err = s3.PutFile("/gdps_sfx/library_version.txt", ver)
	if err != nil {
		SendMessageDiscord("⚠️ Failed to save SFX Library Version: " + err.Error())
		return
	}
}

func GetConsulKV() (consulKV *consul.KV, err error) {
	consulConf := consul.DefaultConfig()
	consulConf.Address = GetEnv("CONSUL_ADDR", "127.0.0.1")
	consulConf.Token = GetEnv("CONSUL_TOKEN", "")
	consulConf.Datacenter = GetEnv("CONSUL_DC", "m41")
	consulCli, err := consul.NewClient(consulConf)
	if err != nil {
		log.Println("Unable to connect to Consul cluster. Assuming self-leadership: " + err.Error())
		return nil, err
	}
	KvEngine = consulCli.KV()
	SessID, _, err := consulCli.Session().Create(&consul.SessionEntry{Name: "FiberAPI", TTL: "5m"}, nil)
	SessionID = SessID
	if err != nil {
		log.Println("Unable to connect to create Consul Session. Assuming self-leadership: " + err.Error())
		return nil, err
	}
	ucron.Every(30).Seconds().Do(func() {
		consulCli.Session().Renew(SessID, nil)
	})
	return KvEngine, nil
}

func PrepareElection(config GlobalConfig) {
	LEAD_CONFIG = config

	KvEngine, _ = GetConsulKV()

	if KvEngine == nil {
		log.Println("Unable to connect to create Session. Assuming self-leadership")
		LEADER = true
	} else {
		AquireLeadership()
		if !LEADER {
			log.Println("Couldn't acquire leadership. Dispatching 10sec watchdog")
			if _, err := ucron.Every(10).Seconds().Do(AquireLeadership); err != nil {
				log.Println(err)
			}
		}
	}
	_, err := ucron.Every(1).Day().At("00:00").Do(MaintainTasks)
	if err != nil {
		log.Println("CANNOT LAUNCH TASKS")
	}
	ucron.StartAsync()
}

func AquireLeadership() {
	kvData := &consul.KVPair{
		Key:     "sessions/ghostcore_lead",
		Value:   []byte(EnvOrDefault("NOMAD_SHORT_ALLOC_ID", "default")),
		Session: SessionID,
	}
	isAcq, _, err := KvEngine.Acquire(kvData, nil)
	if err == nil && isAcq {
		if LEADER {
			log.Println("Still leader (ensuring tasks)")
		} else {
			log.Println("Lock was successfully acquired. NOW LEADER")
			LEADER = true
		}
	} else {
		if LEADER {
			log.Println("Couldn't acquire leadership. Stepped down by force.")
			LEADER = false
		} else {
			log.Println("Couldn't aquire leadership. Still follower")
		}
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
