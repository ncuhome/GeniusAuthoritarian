package store

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/aliyun/oss"
	log "github.com/sirupsen/logrus"
)

var Root *oss.Storage

var Avatar *oss.Storage

func init() {
	var e error
	Root, e = oss.NewRoot(global.Config.Aliyun.Endpoint,
		global.Config.Aliyun.AccessKey, global.Config.Aliyun.SecretKey,
		global.Config.Aliyun.Bucket)
	if e != nil {
		log.Fatalf("初始化 oss 失败: %v", e)
	}

	Avatar = Root.NewDir("avatar")
}
