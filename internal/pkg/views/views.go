package views

import (
	"github.com/Mmx233/daoUtil"
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/agent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	log "github.com/sirupsen/logrus"
	"time"
)

func InitRenewAgent() {
	_, e := agent.AddRegular(&agent.Event{
		T: "0 6,12,16,20,23 * * *",
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
	appSrv, e := service.App.Begin()
	if e != nil {
		return e
	}
	defer appSrv.Rollback()

	apps, e := appSrv.GetForUpdateViews(daoUtil.LockForUpdate)
	if e != nil {
		return e
	}

	loginRecordSrv := service.LoginRecordSrv{DB: appSrv.DB}
	for _, app := range apps {
		var loginRecordIdList []uint
		loginRecordIdList, e = loginRecordSrv.GetViewIDs(app.ID, app.ViewsID)
		if e != nil {
			return e
		}

		if len(loginRecordIdList) != 0 {
			e = appSrv.UpdateViews(app.ID, loginRecordIdList[0], app.Views+uint64(len(loginRecordIdList)))
			if e != nil {
				return e
			}
		}
	}

	return appSrv.Commit().Error
}
