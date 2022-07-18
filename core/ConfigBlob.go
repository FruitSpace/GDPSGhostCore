package core

import (
	"encoding/json"
)

type MysqlConfig struct {
	Host string
	Port int
	User string
	Password string
	DBName string
}

type LogConfig struct {
	LogEnable bool
	LogDB bool
	LogEndpoints bool
	LogRequests bool
}

type ChestConfig struct {
	ChestSmallOrbsMin int
	ChestSmallOrbsMax int
	ChestSmallDiamondsMin int
	ChestSmallDiamondsMax int
	ChestSmallShards []int
	ChestSmallKeysMin int
	ChestSmallKeysMax int
	ChestSmallWait int

	ChestBigOrbsMin int
	ChestBigOrbsMax int
	ChestBigDiamondsMin int
	ChestBigDiamondsMax int
	ChestBigShards []int
	ChestBigKeysMin int
	ChestBigKeysMax int
	ChestBigWait int
}

type ServerConfig struct {
	SrvID string
	SrvKey string
	MaxUsers int
	MaxLevels int
	MaxComments int
	MaxPosts int
	HalMusic bool
	Locked bool
}

type GlobalConfig struct {
	ApiEndpoint string
	LogConnector string
	LogEndpoint string
	MaintenanceMode bool
	RedisHost string
	RedisPort string
	RedisPassword string
	RedisDB int
	SavePath string
	ModuleSettings map[string]string
}

type ConfigBlob struct {
	DBConfig     MysqlConfig
	LogConfig    LogConfig
	ChestConfig  ChestConfig
	ServerConfig ServerConfig
}

func (glob *GlobalConfig) LoadById(Srvid string) (ConfigBlob, error){
	rdb:=RedisConn{}
	log:=Logger{}
	if err:=rdb.ConnectBlob(*glob); err!=nil {
		log.LogWarn(rdb,err.Error())
		return ConfigBlob{},err
	}
	conf:=ConfigBlob{}
	t,err:=rdb.DB.Get(rdb.context,Srvid).Result()
	if err!=nil {
		return ConfigBlob{},err
	}
	err=json.Unmarshal([]byte(t),&conf)
	if err!=nil{
		log.LogWarn(rdb,err.Error())
		return ConfigBlob{},err
	}
	rdb.DB.Close()
	return conf, nil
}