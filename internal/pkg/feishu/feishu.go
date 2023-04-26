package feishu

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/tools"
	"log"
)

func init() {
	if e := RunGroupSync(); e != nil {
		log.Fatalf("添加定时同步飞书部门任务失败: %v", e)
	}
}

var Api = feishu.New(global.Config.Feishu.ClientID, global.Config.Feishu.Secret, tools.Http.Client)
