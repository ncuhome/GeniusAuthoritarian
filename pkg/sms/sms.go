package sms

import (
	"bytes"
	"errors"
	"github.com/Mmx233/tool"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
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

func (a Ums) transformEncoding(i io.Reader, encoder transform.Transformer) (string, error) {
	reader := transform.NewReader(i, encoder)
	d, err := io.ReadAll(reader)
	return string(d), err
}

func (a Ums) Send(msg string, phone string) error {
	msgGbk, err := a.transformEncoding(bytes.NewBuffer([]byte(msg)), simplifiedchinese.GBK.NewEncoder())

	res, err := a.http.PostRequest(&tool.DoHttpReq{
		Url: "https://smsapi.ums86.com:9600/sms/Api/Send.do",
		Header: map[string]interface{}{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: map[string]interface{}{
			"SpCode":         a.conf.SpCode,
			"LoginName":      a.conf.LoginName,
			"Password":       a.conf.Password,
			"MessageContent": msgGbk,
			"UserNumber":     phone,
			"SerialNumber":   time.Now().UnixNano(),
		},
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	data, err := url.ParseQuery(string(resBytes))
	if err != nil {
		return err
	} else if data.Get("result") != "0" {
		return errors.New(string(resBytes))
	}
	return nil
}
