package feishu

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishuApi"
	"github.com/ncuhome/GeniusAuthoritarian/tools"
	"log"
)

func init() {
	if e := DepartmentSync(); e != nil {
		log.Fatalf("同步飞书部门失败: %v", e)
	}
	if e := UserSync(); e != nil {
		log.Fatalf("同步飞书用户失败: %v", e)
	}

	if e := AddDepartmentSyncCron("0 5 * * *"); e != nil {
		log.Fatalf("添加定时同步飞书部门任务失败: %v", e)
	}
	if e := AddUserSyncCron("30 5 * * *"); e != nil {
		log.Fatalf("添加定时同步飞书用户任务失败: %v", e)
	}
}

var Api = feishuApi.New(global.Config.Feishu.ClientID, global.Config.Feishu.Secret, tools.Http.Client)
