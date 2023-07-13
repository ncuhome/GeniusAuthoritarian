package views

import (
	"container/list"
	"github.com/Mmx233/daoUtil"
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
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
	loginRecordList, e := loginRecordSrv.GetMultipleViewsIDs(apps)
	if e != nil {
		return e
	}

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

			if loginRecordIndex > len(loginRecordList) {
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
		e = appSrv.UpdateViews(app.ID, app.ViewsID, app.Views)
		if e != nil {
			return e
		}
	}

	return appSrv.Commit().Error
}
