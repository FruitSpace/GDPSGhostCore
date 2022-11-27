//+build linux,amd64 windows,amd64

package main

import (
	"HalogenGhostCore/api"
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	"github.com/getsentry/sentry-go"
	"log"
	"os"
	"time"
)

func main() {
	// Start Sentry so I can sleep well
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://c25abcf97ba74583b8cbeaccfa448669@o1404863.ingest.sentry.io/6737893",
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2*time.Second)

	config:=GenGConfig()
	//core.MaintainTasks(config)
	core.MaintainRoutines(config)

	ghostServer:= api.GhostServer{
		Log: core.Logger{
			Output: connectors.GetWriter("",""),
		},
		Config: config,
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
			[]int{1,2,3,4,5,6},1,6, 3600,
			2000,4000,20,100,
			[]int{1,2,3,4,5,6},1,6, 14400,
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
			100,
			map[string]bool{"discord":true},
		},
		SecurityConfig: core.SecurityConfig{
			DisableProtection: false,
			AutoActivate: false,
			BannedIPs: []string{},
		},
	}
}

func GenGConfig() core.GlobalConfig {
	return core.GlobalConfig{
		EnvOrDefault("MASTER_KEY", "3XTR4OrD1nArY_K3Y_1907"),
		EnvOrDefault("API_ENDPOINT", "https://api.fruitspace.one/sched/gd/api"),
		"stdout",
		"null",
		false,
		EnvOrDefault("REDIS_HOST", "localhost"),
		EnvOrDefault("REDIS_PORT", "6379"),
		EnvOrDefault("REDIS_PASSWORD", ""),
		7,
		EnvOrDefault("SAVE_PATH", "./"),

		map[string]string{
			"rabbitmq_host":"auto",
			"rabbitmq_user":"m41dss",
			"rabbitmq_password":EnvOrDefault("RABBITMQ_PASSWORD", "None"),
		},
	}
}

func EnvOrDefault(key string, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}