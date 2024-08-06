package main

import (
	"HalogenGhostCore/api"
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"time"
)

func main() {
	// Start Sentry so I can sleep well
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://ef8c6a708a684aa78fdfc0be5a85115b@o1404863.ingest.sentry.io/4504374313222144",
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	config := GenGConfig()

	time.Local, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Println(err)
	}

	DB_USER := EnvOrDefault("DB_USER", "")
	DB_PASS := EnvOrDefault("DB_PASS", "")
	DB_HOST := EnvOrDefault("DB_HOST", "")

	core.DBTunnel, err = sqlx.Connect("mysql", DB_USER+":"+DB_PASS+"@tcp("+DB_HOST+":3306)/")
	if err != nil {
		log.Println("Error while connecting to " + DB_USER + "@localhost: " + err.Error())
		time.Sleep(10 * time.Second)
		main()
	}
	core.DBTunnel.SetMaxOpenConns(100)
	ghostServer := api.GhostServer{
		Log: core.Logger{
			Output: connectors.GetWriter("", ""),
		},
		Config: config,
	}
	api.InitCache(config, 16)
	core.PrepareElection(config)
	defer core.StepDown()

	ghostServer.StartServer("0.0.0.0:1997")
}

func GenGConfig() core.GlobalConfig {
	return core.GlobalConfig{
		MasterKey:      EnvOrDefault("MASTER_KEY", "3XTR4OrD1nArY_K3Y_1907"),
		ApiEndpoint:    EnvOrDefault("API_ENDPOINT", "http://127.0.0.1:6000/sched/gd/api"),
		LogConnector:   "stdout",
		LogEndpoint:    "null",
		RedisHost:      EnvOrDefault("REDIS_HOST", "localhost"),
		RedisPort:      EnvOrDefault("REDIS_PORT", "6379"),
		RedisPassword:  EnvOrDefault("REDIS_PASSWORD", ""),
		RedisDB:        7,
		SavePath:       EnvOrDefault("SAVE_PATH", "./"),
		ModuleSettings: map[string]string{},
	}
}

func EnvOrDefault(key string, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
