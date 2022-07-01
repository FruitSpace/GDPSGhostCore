package core

import "encoding/json"

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
	ChestSmallShardsMin int
	ChestSmallShardsMax int
	ChestSmallKeysMin int
	ChestSmallKeysMax int
	ChestSmallWait int

	ChestBigOrbsMin int
	ChestBigOrbsMax int
	ChestBigDiamondsMin int
	ChestBigDiamondsMax int
	ChestBigShardsMin int
	ChestBigShardsMax int
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
	if err:=rdb.ConnectBlob(*glob); err!=nil {return ConfigBlob{},err}
	conf:=ConfigBlob{}
	t:=rdb.DB.Get(rdb.context,Srvid)
	if err:=t.Err();err!=nil { return ConfigBlob{},err}
	err:=json.Unmarshal([]byte(t.String()),&conf)
	if err!=nil{return ConfigBlob{},err}
	rdb.DB.Close()
	return conf, nil
}