package api

import (
	"HalogenGhostCore/core"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var BallisticsCache map[string]int64
var BadRepIP map[string]int
var clusterFuck map[string][]string // map[timespan][]ip

const BADREP_THRESHOLD = 20
const CLUSTERFUCK_THRESHOLD = 10

func PrepareBallistics(req *http.Request) bool {
	IPAddr := ipOf(req)
	if thr, ok := BadRepIP[IPAddr]; ok && thr > BADREP_THRESHOLD {
		// Here we block DDoSers from fucking everywhere
		return false
	}
	//rl := ratelimit.New(10)
	// Last timestamp for this path
	tm, ok := BallisticsCache[req.URL.Path]
	if !ok {
		BallisticsCache[req.URL.Path] = time.Now().UnixMilli()
		return false
	}
	BallisticsCache[req.URL.Path] = time.Now().UnixMilli()

	if time.Now().UnixMilli()-tm > 30000 {
		// I mean how the fuck are we going to stop that?
		return false
	}

	if _, ok := BadRepIP[IPAddr]; !ok {
		BadRepIP[IPAddr] = 0
	}

	// If registered from single IP and multiple requests less in than 30 seconds
	if time.Now().UnixMilli()-tm > 2000 && BadRepIP[IPAddr] < 5 {
		BadRepIP[IPAddr]++
		return false
	}
	if time.Now().UnixMilli()-tm > 400 && BadRepIP[IPAddr] < BADREP_THRESHOLD {
		t := fmt.Sprintf("[%s, rep=%d] Unusual request speed .4s-2s at `%s`\nResult: `Throttle 5s`", IPAddr, BadRepIP[IPAddr], req.URL.Path)
		BadRepIP[IPAddr] += 2
		go core.SendMessageDiscord(t)
		time.Sleep(time.Second * 5)
		return false
	}
	clusterFuckTime := time.Now().Format("2006-01-02 15:04")
	if _, ok := clusterFuck[clusterFuckTime]; !ok {
		// Current minute wasn't registered, check previous
		if len(clusterFuck) > 0 {
			// We had previous keys - most likely from previous minute
			for k, v := range clusterFuck {
				// Most likely there will be only one key
				t := fmt.Sprintf("[Ballistics Clusterfuck] Found unfinished block at %s\n```\n%s\n```", k, strings.Join(v, "\n"))
				go core.SendMessageDiscord(t)
				break
			}
			clusterFuck = make(map[string][]string)
		}
		clusterFuck[clusterFuckTime] = []string{}
	}

	clusterFuck[clusterFuckTime] = append(clusterFuck[clusterFuckTime], IPAddr)
	if len(clusterFuck[clusterFuckTime]) > CLUSTERFUCK_THRESHOLD {
		// 99% it's a fucking DDoS
		for _, ip := range clusterFuck[clusterFuckTime] {
			BadRepIP[ip] = 100
		}
		t := fmt.Sprintf("[Ballistics Clusterfuck] Banned IPs at %s\n```\n%s\n```", clusterFuckTime, strings.Join(clusterFuck[clusterFuckTime], "\n"))
		go core.SendMessageDiscord(t)
		return true
	}

	BadRepIP[IPAddr] += 5
	time.Sleep(time.Second * 30)
	return true
}
