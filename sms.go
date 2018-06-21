package alisms

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// Sms 短信发送
type Sms struct {
	AccessKeyID  string
	AccessSecret string
	SignName     string //签名名称
	TemplateCode string //模板code
	Debug        bool
}

// NewSms 创建新的短信发送
func NewSms(accessKeyID string, accessSecret string, signName string, templateCode string, debug bool) *Sms {
	var a Sms
	a.SignName = signName
	a.TemplateCode = templateCode
	a.AccessKeyID = accessKeyID
	a.AccessSecret = accessSecret
	a.Debug = debug

	return &a
}

// Send 发送
func (t *Sms) Send(numbers string, params string) error {
	var request Request
	request.Format = "JSON"
	request.Version = "2017-05-25"
	request.AccessKeyID = t.AccessKeyID
	request.SignatureMethod = "HMAC-SHA1"
	request.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	request.SignatureVersion = "1.0"
	request.SignatureNonce = GetRandString(16)

	request.Action = "SendSms"
	request.SignName = t.SignName
	request.TemplateCode = t.TemplateCode
	request.PhoneNumbers = numbers
	request.TemplateParam = params
	request.RegionID = "cn-hangzhou"

	url := request.ComposeUrl("GET", t.AccessSecret)
	var resp *http.Response
	var err error
	resp, err = http.Get(url)
	if t.Debug {
		println("GET ->", url)
	}
	if err != nil {
		return err
	}
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if t.Debug {
		println("Response ->", string(b))
	}
	_m := make(map[string](string))
	err = json.Unmarshal(b, &_m)
	if err != nil {
		return err
	}
	message, ok := _m["Message"]
	if ok && strings.ToUpper(message) == "OK" {
		return nil
	}
	if ok {
		return errors.New(message)
	}
	return errors.New("send sms error")
}

var randSourceBytes = append(GetByteRange(97, 122), GetByteRange(65, 90)...)

// GetByteRange get a range bytes
func GetByteRange(start, stop byte) []byte {
	r := make([]byte, stop-start+1)
	i := 0
	for s := start; s <= stop; s++ {
		r[i] = s
		i++
	}
	return r
}

// GetRandString 生成随机字符串
func GetRandString(size int) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	buff := make([]byte, size)
	for i := 0; i < size; i++ {
		buff[i] = randSourceBytes[r.Intn(len(randSourceBytes))]
	}
	return string(buff)
}
