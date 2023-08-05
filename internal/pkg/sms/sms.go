package sms

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/sms"
)

var Ums sms.Ums

func init() {
	Ums = sms.New(sms.UmsConf{
		SpCode:    global.Config.Sms.SpCode,
		LoginName: global.Config.Sms.LoginName,
		Password:  global.Config.Sms.Password,
		Client:    tools.Http.Client,
	})
}
