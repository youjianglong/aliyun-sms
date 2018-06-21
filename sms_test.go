package alisms

import "testing"

func TestGetRandString(t *testing.T) {
	t.Log(GetRandString(16))
}

func TestSms(t *testing.T) {
	sms := NewSms("accessKeyID", "accessSecret", "signName", "templateCode", true)
	err := sms.Send("13800138000", `{"code":"1234"}`)
	if err != nil {
		t.Fatal(err)
	}
}
