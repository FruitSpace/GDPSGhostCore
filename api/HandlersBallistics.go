package api

import (
	"HalogenGhostCore/core"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var BallisticsCache map[string]int64
var BadRepIP map[string]int

func PrepareBallistics(req *http.Request) bool {
	IPAddr := req.Header.Get("CF-Connecting-IP")
	if IPAddr == "" {
		IPAddr = req.Header.Get("X-Real-IP")
	}
	if IPAddr == "" {
		IPAddr = strings.Split(req.RemoteAddr, ":")[0]
	}
	tm, ok := BallisticsCache[req.URL.Path]
	if !ok {
		BallisticsCache[req.URL.Path] = time.Now().UnixMilli()
		return false
	}
	if _, ok := BadRepIP[IPAddr]; !ok {
		BadRepIP[IPAddr] = 0
	}
	BallisticsCache[req.URL.Path] = time.Now().UnixMilli()
	if time.Now().UnixMilli()-tm > 2000 && BadRepIP[IPAddr] < 5 {
		BadRepIP[IPAddr]++
		return false
	}
	if time.Now().UnixMilli()-tm > 500 && BadRepIP[IPAddr] < 20 {
		t := fmt.Sprintf("[%s, r=%d] Unusual request speed .4s-2s at `%s`\nResult: `Throttle 5s`", IPAddr, BadRepIP[IPAddr], req.URL.Path)
		BadRepIP[IPAddr] += 2
		log.Println(t)
		go core.SendMessageDiscord(t)
		time.Sleep(time.Second * 5)
		return false
	}
	BadRepIP[IPAddr]++
	t := fmt.Sprintf("[%s, r=%d] Possible DDoS detected (d=%d) at `%s`\nResult: `Throttle 30s + DROP`", IPAddr, BadRepIP[IPAddr], time.Now().UnixMilli()-tm, req.URL.Path)
	log.Println(t)
	go core.SendMessageDiscord(t)
	time.Sleep(time.Second * 30)
	return true
}
