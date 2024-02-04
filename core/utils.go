package core

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
	"golang.org/x/exp/slices"
	"html"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var loc, _ = time.LoadLocation("Europe/Moscow")

func ClearGDRequest(str string) string {
	return strings.TrimSpace(
		strings.Split(
			strings.Split(
				strings.Split(
					strings.Split(
						strings.Split(
							strings.TrimSpace(
								html.EscapeString(str)), ":")[0], "|")[0], "~")[0], "#")[0], ")")[0])

}

func DoXOR(text string, key string) (output string) {
	for i := 0; i < len(text); i++ {
		output += string(text[i] ^ key[i%len(key)])
	}
	return output
}

func GetDateAgo(date int64) string {
	diff := time.Now().Unix() - date
	if diff < 60 {
		return strconv.FormatInt(diff, 10) + " seconds"
	}
	if diff < 3600 {
		return strconv.FormatInt(diff/60, 10) + " minutes"
	}
	if diff < 86400 {
		return strconv.FormatInt(diff/3600, 10) + " hours"
	}
	if diff < 604800 {
		return strconv.FormatInt(diff/86400, 10) + " days"
	}
	if diff < 604800*4 {
		return strconv.FormatInt(diff/604800, 10) + " weeks"
	}
	if diff < 604800*4*12 {
		return strconv.FormatInt(diff/(604800*4), 10) + " months"
	}
	return strconv.FormatInt(diff/(604800*4*12), 10) + " years"
}

func HashSolo(levelstring string) string {
	hash := make([]byte, 40)
	p := 0
	plen := len(levelstring)
	for i := 0; i < plen; i += (plen / 40) {
		if p > 39 {
			break
		}
		hash[p] = levelstring[i]
		p++
	}
	sha := sha1.New()
	sha.Write([]byte(string(hash) + "xI25fpAapCQg"))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func HashSolo2(lvlstring string) string {
	sha := sha1.New()
	sha.Write([]byte(lvlstring + "xI25fpAapCQg"))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func HashSolo3(lvlstring string) string {
	sha := sha1.New()
	sha.Write([]byte(lvlstring + "oC36fpYaPtdg"))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func HashSolo4(lvlstring string) string {
	sha := sha1.New()
	sha.Write([]byte(lvlstring + "pC26fpYaQCtg"))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func DoGjp(gjp string) string {
	gjp = strings.ReplaceAll(strings.ReplaceAll(gjp, "_", "/"), "-", "+")
	block, err := base64.StdEncoding.DecodeString(gjp)
	if err != nil {
		return ""
	}
	return DoXOR(string(block), "37526")
}

func DoGjp2(password string) string {
	sha := sha1.New()
	sha.Write([]byte(password + "mI29fmAnxgTs"))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func MD5(str string) string {
	md := md5.New()
	md.Write([]byte(str))
	return fmt.Sprintf("%x", md.Sum(nil))
}

func SHA256(str string) string {
	sha := sha256.New()
	sha.Write([]byte(str))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func SHA512(str string) string {
	sha := sha512.New()
	sha.Write([]byte(str))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

type Logger struct {
	Output io.Writer
}

func (lg *Logger) LogErr(module interface{}, message string) {
	sentry.CaptureMessage(message)
	if fmt.Sprintf("%T", module) == "*core.MySQLConn" {
		message = "[MySQL " + module.(*MySQLConn).DBName + "]" + message
	}
	fmt.Println("ERR: ", message)
	ReportFail(message)
}
func (lg *Logger) LogWarn(module interface{}, message string) {
	fmt.Printf("[%T] %s\n", module, message)
}

func (lg *Logger) Must(err error) {
	if err != nil {
		sentry.CaptureException(err)
		fmt.Println("ERR:", err.Error())
		ReportFail(err.Error())
		panic("Must be dereferenced")
	}
}

func (lg *Logger) Should(err error) error {
	if err != nil && err != redis.Nil && err.Error() != ".ignore" {
		sentry.CaptureException(err)
		ReportFail(err.Error())
		lg.LogWarn(err, err.Error())
	}
	return err
}

func sliceRemove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func FilterEmail(email string) bool {
	semail := strings.Split(email, "@")
	if len(semail) != 2 {
		return false
	}
	AllowedEmailProviders := []string{
		"yandex.ru",
		"ya.ru",
		"mail.ru",
		"gmail.com",
		"aol.com",
		"rambler.ru",
		"bk.ru",
		"vk.com",
	}
	if !slices.Contains(AllowedEmailProviders, strings.ToLower(semail[1])) {
		return false
	}
	return true
}

func TryInt(target *int, str string) bool {
	if c, err := strconv.Atoi(str); err == nil {
		*target = c
		return true
	}
	return false
}

func ToInt(b bool) int {
	var i int
	if b {
		i = 1
	}
	return i
}

func GetGDVersion(post url.Values) int {
	version := 21
	if post.Has("gameVersion") {
		if c, err := strconv.Atoi(post.Get("gameVersion")); err == nil {
			version = c
		}
	}
	if version == 20 {
		if post.Has("binaryVersion") {
			if c, err := strconv.Atoi(post.Get("binaryVersion")); err == nil && c > 27 {
				version++
			}
		}
	}
	if version == 21 {
		if post.Has("binaryVersion") {
			if c, err := strconv.Atoi(post.Get("binaryVersion")); err == nil && c > 36 {
				version++
			}
		}
	}
	return version
}

func CheckGDAuth(post url.Values) bool {
	version := GetGDVersion(post)
	if post.Get("accountID") != "" && ((version < 22 && post.Get("gjp") != "") || (version == 22 && post.Get("gjp2") != "")) {
		return true
	}
	return false
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func CleanDoubles(src string, req string) string {
	for strings.Contains(src, req+req) {
		src = strings.ReplaceAll(src, req+req, req)
	}
	return src
}

func Decompose(src string, del string) []int {
	mako := strings.Split(src, del)
	var vs []int
	for _, l := range mako {
		i, err := strconv.Atoi(l)
		if err != nil {
			continue
		}
		vs = append(vs, i)
	}
	return vs
}

func ArrTranslate(arr []int) []string {
	var vs []string
	for _, l := range arr {
		vs = append(vs, strconv.Itoa(l))
	}
	return vs
}

func QuickComma(str string) string {
	return strings.Join(ArrTranslate(Decompose(CleanDoubles(str, ","), ",")), ",")
}

func InArray(arr []string, ele string) bool {
	for _, v := range arr {
		if v == ele {
			return true
		}
	}
	return false
}

func SendMessageDiscord(text string) {
	b, _ := json.Marshal(map[string]string{
		"content": text,
	})

	content := bytes.NewReader(b)

	http.Post("https://discord.com/api/webhooks/1133487072624263238/QSamyTrVKSE9ASxs3VDY_p9Yuawfy-JjtjLna3VCL2qpFsmp2W8N00aAM9PDt12-yPlt",
		"application/json", content)
}

func ReportFail(err string) {
	//http.PostForm("https://api.fruitspace.one/pandora/report", url.Values{"error":{url.QueryEscape(err)}})
	SendMessageDiscord(err)
}

type S3FS struct {
	Endpoint  string
	AccessKey string
	SecretKey string

	Region string
	Bucket string
}

func (s3fs *S3FS) GetFile(path string) ([]byte, error) {
	creds := credentials.NewStaticCredentials(s3fs.AccessKey, s3fs.SecretKey, "")
	cfg := aws.NewConfig().WithEndpoint(s3fs.Endpoint).WithRegion(s3fs.Region).WithCredentials(creds)
	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}
	svc := s3manager.NewDownloader(sess)

	buf := aws.NewWriteAtBuffer([]byte{})
	_, err = svc.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(s3fs.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s3fs *S3FS) PutFile(path string, data []byte) error {
	creds := credentials.NewStaticCredentials(s3fs.AccessKey, s3fs.SecretKey, "")
	cfg := aws.NewConfig().WithEndpoint(s3fs.Endpoint).WithRegion(s3fs.Region).WithCredentials(creds)
	sess, err := session.NewSession(cfg)
	if err != nil {
		return err
	}
	svc := s3manager.NewUploader(sess)

	_, err = svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3fs.Bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s3fs *S3FS) DeleteFile(path string) error {
	creds := credentials.NewStaticCredentials(s3fs.AccessKey, s3fs.SecretKey, "")
	cfg := aws.NewConfig().WithEndpoint(s3fs.Endpoint).WithRegion(s3fs.Region).WithCredentials(creds)
	sess, err := session.NewSession(cfg)
	if err != nil {
		return err
	}
	svc := s3.New(sess)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s3fs.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return err
	}
	return nil
}

func NewS3FS() *S3FS {
	Props := GetKVEnv("S3_CONFIG")
	return &S3FS{
		Endpoint:  Props["endpoint"],
		AccessKey: Props["access_key"],
		SecretKey: Props["secret"],
		Region:    Props["region"],
		Bucket:    Props["bucket"],
	}
}

func GetEnv(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func GetKVEnv(key string) map[string]string {
	val := os.Getenv(key)
	if val == "" {
		return map[string]string{}
	}
	kv := map[string]string{}
	for _, v := range strings.Split(val, ",") {
		kv[strings.SplitN(v, "=", 2)[0]] = strings.SplitN(v, "=", 2)[1]
	}
	return kv
}

func SendAPIWebhook(srvid string, xtype string, data map[string]string) {
	body, _ := json.Marshal(data)
	http.Post("https://api.fruitspace.one/v2/internal/gd/"+srvid+"/webhook?type="+xtype,
		"application/json", bytes.NewReader(body))
}

func DiffToText(stars int, demonDiff int, isFeatured int, isEpic int) string {
	diff := "unrated"
	switch stars {
	case 1:
		diff = "auto"
	case 2:
		diff = "easy"
	case 3:
		diff = "normal"
	case 4:
		fallthrough
	case 5:
		diff = "hard"
	case 6:
		fallthrough
	case 7:
		diff = "harder"
	case 8:
		fallthrough
	case 9:
		diff = "insane"
	case 10:
		diff = "demon"
		switch demonDiff {
		case 3:
			diff += "-easy"
		case 4:
			diff += "-medium"
		case 5:
			diff += "-insane"
		case 6:
			diff += "-extreme"
		default:
			diff += "-hard"
		}
	}
	if isEpic > 0 {
		return diff + "-epic"
	}
	switch isFeatured {
	case 1:
		return diff + "-featured"
	case 2:
		return diff + "-epic"
	case 3:
		return diff + "-mythic"
	case 4:
		return diff + "-legendary"
	}
	return diff
}

type GoMetrics struct {
	startTime time.Time
	lastTime  time.Time
	steps     map[string]int64 // ms took
	lastStep  string
}

func NewGoMetrics() *GoMetrics {
	return &GoMetrics{
		startTime: time.Now(),
		lastTime:  time.Now(),
		steps:     make(map[string]int64),
	}
}

func (gm *GoMetrics) Reset() {
	gm.startTime = time.Now()
	gm.lastTime = time.Now()
}

func (gm *GoMetrics) NewStep(name string) {
	gm.steps[name] = 0
	t := time.Now()
	if len(gm.lastStep) != 0 {
		gm.steps[gm.lastStep] = t.Sub(gm.lastTime).Milliseconds()
	}
	gm.lastStep = name
	gm.lastTime = t // Update total time
}

func (gm *GoMetrics) ExplicitDoneStep(name string) {
	gm.steps[name] = time.Now().Sub(gm.lastTime).Milliseconds()
	gm.lastStep = ""
	gm.lastTime = time.Now()
}

func (gm *GoMetrics) Done() {
	gm.steps[gm.lastStep] = time.Now().Sub(gm.lastTime).Milliseconds()
	gm.lastStep = ""
	gm.lastTime = time.Now()
}

func (gm *GoMetrics) dumpText(del string) string {
	out := ""
	for k, v := range gm.steps {
		out += fmt.Sprintf("%s: %dms"+del, k, v)
	}
	out += fmt.Sprintf("Total: %dms"+del, gm.lastTime.Sub(gm.startTime).Milliseconds())
	return out
}

func (gm *GoMetrics) DumpText() string {
	return gm.dumpText("\n")
}

func (gm *GoMetrics) DumpTextInline() string {
	return gm.dumpText("; ")
}

func (gm *GoMetrics) DumpJSON() string {
	// copy gm.stats
	data, _ := json.Marshal(struct {
		Steps map[string]int64
		Total int64
	}{
		Steps: gm.steps,
		Total: gm.lastTime.Sub(gm.startTime).Milliseconds(),
	})
	return string(data)
}
