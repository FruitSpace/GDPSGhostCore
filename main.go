//+build linux,amd64 windows,amd64

package main

import (
	"HalogenGhostCore/api"
	"HalogenGhostCore/core"
	"HalogenGhostCore/core/connectors"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

func main() {
	// Start Sentry so I can sleep well
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://5ee98ff065064ac5a4d3e96a55f8cd08@o1368861.ingest.sentry.io/6671765",
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2*time.Second)
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
		"Zero",
		"https://halhost.cc/app/api/gdps_api.php",
		"stdout",
		"null",
		false,
		"localhost",
		"6379",
		"3XTR4OrD1nArY_K3Y_1907",
		7,
		"./",

		map[string]string{
			"rabbitmq_host":"auto",
			"rabbitmq_user":"m41dss",
			"rabbitmq_password":"passw",
		},
	}
}