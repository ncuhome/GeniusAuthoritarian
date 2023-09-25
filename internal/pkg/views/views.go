package views

// 应用访问量增量刷新

import (
	"container/list"
	"github.com/Mmx233/daoUtil"
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/cronAgent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	log "github.com/sirupsen/logrus"
	"time"
)

func InitRenewAgent() {
	_, err := cronAgent.AddRegular(&cronAgent.Event{
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
	if err != nil {
		panic(err)
	}
}

func Renew() error {
	appSrv, err := service.App.Begin()
	if err != nil {
		return err
	}
	defer appSrv.Rollback()

	apps, err := appSrv.GetForUpdateViews(daoUtil.LockForUpdate)
	if err != nil {
		return err
	}

	loginRecordSrv := service.LoginRecordSrv{DB: appSrv.DB}
	loginRecordList, err := loginRecordSrv.GetMultipleViewsIDs(apps)
	if err != nil {
		return err
	}

	if len(apps) != 0 && len(loginRecordList) != 0 {
		var appIndex, loginRecordIndex int
		var loginRecord = loginRecordList[loginRecordIndex]
		var appUpdated = list.New() // *dao.App
		for appIndex < len(apps) {
			app := apps[appIndex]

			if app.ID != loginRecord.AID {
				appIndex++
				continue
			}

			appUpdated.PushBack(&app)
			app.ViewsID = loginRecord.ID
			for {
				app.Views++
				loginRecordIndex++

				if loginRecordIndex >= len(loginRecordList) {
					goto doUpdate
				}

				loginRecord = loginRecordList[loginRecordIndex]
				if loginRecord.AID != app.ID {
					appIndex++
					break
				}
			}
		}

	doUpdate:

		for el := appUpdated.Front(); el != nil; el = el.Next() {
			app := el.Value.(*dao.App)
			err = appSrv.UpdateViews(app.ID, app.ViewsID, app.Views)
			if err != nil {
				return err
			}
		}
	}

	return appSrv.Commit().Error
}
