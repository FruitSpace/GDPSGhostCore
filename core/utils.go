package core

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)


func DoXOR(text string, key string) (output string) {
	for i:=0;i<len(text);i++ {
		output+=string(text[i] ^ key[i%len(key)])
	}
	return output
}

func GetDateAgo(date int64) string {
	diff:=time.Now().Unix()-date
	if diff<60 {return strconv.FormatInt(diff,10)+" seconds"}
	if diff<3600 {return strconv.FormatInt(diff/60,10)+" minutes"}
	if diff<86400 {return strconv.FormatInt(diff/3600,10)+" hours"}
	if diff<604800 {return strconv.FormatInt(diff/86400,10)+" days"}
	if diff<604800*4 {return strconv.FormatInt(diff/604800,10)+" weeks"}
	if diff<604800*4*12 {return strconv.FormatInt(diff/(604800*4),10)+" months"}
	return strconv.FormatInt(diff/(604800*4*12),10)+" years"
}

func HashSolo(levelstring string) string {
	hash:=make([]byte,40)
	p:=0
	plen:=len(levelstring)
	for i:=0; i<plen; i+=(plen/40) {
		if p>39 {break}
		hash[p]=levelstring[i]
		p++
	}
	sha:=sha1.New()
	sha.Write([]byte(string(hash)+"xI25fpAapCQg"))
	return fmt.Sprintf("%x",sha.Sum(nil))
}

func HashSolo2(lvlstring string) string {
	sha:=sha1.New()
	sha.Write([]byte(lvlstring+"xI25fpAapCQg"))
	return fmt.Sprintf("%x",sha.Sum(nil))
}

func Hashsolo3(lvlstring string) string {
	sha:=sha1.New()
	sha.Write([]byte(lvlstring+"oC36fpYaPtdg"))
	return fmt.Sprintf("%x",sha.Sum(nil))
}

func HashSolo4(lvlstring string) string {
	sha:=sha1.New()
	sha.Write([]byte(lvlstring+"pC26fpYaQCtg"))
	return fmt.Sprintf("%x",sha.Sum(nil))
}

func DoGjp(gjp string) string {
	gjp=strings.ReplaceAll(strings.ReplaceAll(gjp,"_","/"),"-","+")
	block,err:=base64.StdEncoding.DecodeString(gjp)
	if err!=nil {return ""}
	return DoXOR(string(block),"37526")
}

func DoGjp2(password string) string {
	sha:=sha1.New()
	sha.Write([]byte(password+"mI29fmAnxgTs"))
	return fmt.Sprintf("%x",sha.Sum(nil))
}

func MD5(str string) string {
	md:=md5.New()
	md.Write([]byte(str))
	return fmt.Sprintf("%x",md.Sum(nil))
}

type Logger struct {
	Output io.Writer
}

func (lg *Logger) LogErr(module interface{}, message string) {
	log.SetOutput(lg.Output)
	log.Panicf("[%T] %s\n",module,message)
}

func (lg *Logger) Must(i interface{}, err error) interface{}{
	if err!=nil {
		lg.LogErr(i,err.Error())
	}
	return i
}


func sliceRemove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}


func FilterEmail(email string) bool {
	semail:=strings.Split(email,"@")
	if len(semail)!=2 {return false}
	AllowedEmailProviders:=[]string{
		"yandex.ru",
		"ya.ru",
		"mail.ru",
		"gmail.com",
		"aol.com",
		"rambler.ru",
		"bk.ru",
		"vk.com",
	}
	if !slices.Contains(AllowedEmailProviders,semail[1]) {return false}
	return true
}