package sms

import (
	"errors"
	"fmt"
	"github.com/Mmx233/tool"
	"github.com/ncuhome/urlencode"
	"net/http"
	"net/url"
	"time"
)

func New(conf UmsConf) Ums {
	return Ums{
		conf: conf,
		http: tool.NewHttpTool(conf.Client),
	}
}

type Ums struct {
	http *tool.Http
	conf UmsConf
}

func (a Ums) Send(msg string, phone string) error {
	var e error
	msg, e = urlencode.Encode(msg, "gbk")
	if e != nil {
		return e
	}
	if req, e := http.Post(
		fmt.Sprintf(
			"https://smsapi.ums86.com:8888/sms/Api/Send.do?SpCode=%s&LoginName=%s&Password=%s&MessageContent=%s&UserNumber=%s&SerialNumber=%v",
			a.conf.SpCode, a.conf.LoginName, a.conf.Password, msg, phone, time.Now().UnixNano()),
		"application/x-www-form-urlencoded", nil); e != nil {
		return e
	} else if data, e := a.http.ReadResBodyToString(req.Body); e != nil {
		return e
	} else if res, e := url.ParseQuery(data); e != nil {
		return e
	} else if res.Get("result") != "0" {
		data, _ = urlencode.Decode(data, "gbk")
		return errors.New(data)
	}

	return nil
}
