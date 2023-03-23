package util

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/tools"
)

var Feishu = feishu.New(global.Config.Feishu.ClientID, global.Config.Feishu.Secret, tools.Http.Client)
