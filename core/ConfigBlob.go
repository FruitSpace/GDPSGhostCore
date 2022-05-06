package core

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
}

type ConfigBlob struct {
	DBConfig     MysqlConfig
	LogConfig    LogConfig
	ChestConfig  ChestConfig
	ServerConfig ServerConfig
}

func (blob ConfigBlob) LoadById(Srvid string, glob *GlobalConfig) {

}