package sms

import "net/http"

type UmsConf struct {
	SpCode    string
	LoginName string
	Password  string

	Client *http.Client
}
