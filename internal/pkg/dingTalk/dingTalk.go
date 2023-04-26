package dingTalk

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/dingTalkApi"
)

var Api = dingTalkApi.New(dingTalkApi.Config{
	ClientID: global.Config.DingTalk.ClientID,
	Secret:   global.Config.DingTalk.Secret,
})
