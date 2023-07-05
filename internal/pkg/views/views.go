package views

import (
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/agent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	log "github.com/sirupsen/logrus"
	"time"
)

func InitRenewAgent() {
	_, e := agent.AddRegular(&agent.Event{
		T: "0 */3 * * *",
		E: func() {
			defer tool.Recover()

			startAt := time.Now()

			e := Renew()
			if e != nil {
				log.Errorln("更新 app views 失败:", e)
				return
			}

			log.Infof("App views 刷新成功，耗时 %d ms", time.Now().Sub(startAt).Milliseconds())
		},
	})
	if e != nil {
		panic(e)
	}
}

func Renew() error {
	views, e := service.LoginRecord.GetViewCount()
	if e != nil {
		return e
	}

	appSrv, e := service.App.Begin()
	if e != nil {
		return e
	}
	defer appSrv.Rollback()

	for _, v := range views {
		e = appSrv.UpdateViews(v.ID, v.Views)
		if e != nil {
			return e
		}
	}

	return appSrv.Commit().Error
}
