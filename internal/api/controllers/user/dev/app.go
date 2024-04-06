package controllers

import (
	"context"
	"errors"
	"github.com/Mmx233/daoUtil"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"gorm.io/gorm"
	"net/url"
	"sort"
	"time"
)

func ListOwnedApp(c *gin.Context) {
	uid := tools.GetUserInfo(c).UID

	apps, err := service.App.GetUserOwnedApp(uid)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, apps)
}

func checkCallback(c *gin.Context, callbackStr string) {
	callbackUrl, err := url.Parse(callbackStr)
	if err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	} else if callbackUrl.Hostname() == "localhost" {
		// 允许调用本地服务
		return
	} else if callbackUrl.Scheme != "https" {
		callback.Error(c, callback.ErrForm)
		return
	}

	exist, err := service.SiteWhiteList.Exist(callbackUrl.Host)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	} else if !exist {
		callback.Error(c, callback.ErrSiteNotAllow, err)
		return
	}
}

func ApplyApp(c *gin.Context) {
	var f struct {
		Name         string `json:"name" form:"name" binding:"required,max=20"`
		Callback     string `json:"callback" form:"callback" binding:"url,required"`
		PermitAll    bool   `json:"permitAll" form:"permitAll"`
		PermitGroups []uint `json:"permitGroups" form:"permitGroups"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	checkCallback(c, f.Callback)
	if c.IsAborted() {
		return
	}

	appSrv, err := service.App.Begin()
	if err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}
	defer appSrv.Rollback()

	exist, err := appSrv.NameExist(f.Name, daoUtil.LockForUpdate)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	} else if exist {
		callback.Error(c, callback.ErrAppNameExist)
		return
	}

	uid := tools.GetUserInfo(c).UID
	newApp, err := appSrv.New(uid, f.Name, f.Callback, f.PermitAll)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	var groups = make([]dto.Group, 0)
	if !f.PermitAll && len(f.PermitGroups) != 0 {
		appGroupSrv := service.AppGroupSrv{DB: appSrv.DB}

		if groups, err = appGroupSrv.BindForApp(newApp.ID, f.PermitGroups); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				callback.Error(c, callback.ErrGroupNotFound)
				return
			}
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}
	}

	if err = appSrv.AddAppCodeToRedis(newApp.AppCode); err != nil || appSrv.Commit().Error != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callback.Success(c, dto.AppNew{
		AppShowDetail: dto.AppShowDetail{
			AppShowOwner: dto.AppShowOwner{
				AppShow: dto.AppShow{
					ID:             newApp.ID,
					Name:           newApp.Name,
					Callback:       newApp.Callback,
					PermitAllGroup: newApp.PermitAllGroup,
					LinkOff:        newApp.LinkOff,
				},
				AppCode: newApp.AppCode,
			},
			Groups: groups,
		},
		AppSecret: newApp.AppSecret,
	})
}

func DeleteApp(c *gin.Context) {
	var f struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	uid := tools.GetUserInfo(c).UID

	appSrv, err := service.App.Begin()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	defer appSrv.Rollback()

	appCode, err := appSrv.FirstAppCodeByID(f.ID, uid, daoUtil.LockForUpdate)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			callback.Error(c, callback.ErrAppNotFound)
			return
		}
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	if err = appSrv.DeleteByID(f.ID, uid); err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	if err = redis.AppCode.Del(appCode); err != nil {
		callback.Error(c, callback.ErrUnexpected)
		return
	}

	if err = appSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Default(c)
}

func ModifyApp(c *gin.Context) {
	var f struct {
		ID           uint   `json:"id" form:"id" binding:"required"`
		Name         string `json:"name" form:"name" binding:"required,max=20"`
		Callback     string `json:"callback" form:"callback" binding:"url,required"`
		PermitAll    bool   `json:"permitAll" form:"permitAll"`
		PermitGroups []uint `json:"permitGroups" form:"permitGroups"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	if len(f.PermitGroups) > 1 {
		sort.Sort(tools.UintSlice(f.PermitGroups))
		for i := 1; i < len(f.PermitGroups); i++ {
			if f.PermitGroups[i-1] == f.PermitGroups[i] {
				callback.Error(c, callback.ErrForm)
				return
			}
		}
	}

	appSrv, err := service.App.Begin()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	defer appSrv.Rollback()

	uid := tools.GetUserInfo(c).UID

	app, err := appSrv.FirstAppDetailedByIDForUser(f.ID, uid, daoUtil.LockForUpdate)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	appInfoChanged := false

	if f.Name != app.Name {
		if exist, err := appSrv.NameExist(f.Name); err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		} else if exist {
			callback.Error(c, callback.ErrAppNameExist)
			return
		}

		appInfoChanged = true
	}

	if f.Callback != app.Callback {
		checkCallback(c, f.Callback)
		if c.IsAborted() {
			return
		}
		appInfoChanged = true
	}

	appGroupSrv := service.AppGroupSrv{DB: appSrv.DB}

	// 更新组关系
	var groupPermisstionChanged bool
	if app.PermitAllGroup != f.PermitAll {
		if f.PermitAll {
			err = appGroupSrv.DeleteAllForApp(app.ID)
			if err != nil {
				callback.Error(c, callback.ErrDBOperation, err)
				return
			}
		} else if len(f.PermitGroups) > 0 {
			app.Groups, err = appGroupSrv.BindForApp(app.ID, f.PermitGroups)
			if err != nil {
				callback.Error(c, callback.ErrDBOperation, err)
				return
			}
		}
		groupPermisstionChanged = !f.PermitAll
		appInfoChanged = true
	}

	if !f.PermitAll {
		exGroups := make([]uint, len(app.Groups))
		for i, group := range app.Groups {
			exGroups[i] = group.ID
		}
		sort.Sort(tools.UintSlice(exGroups))

		var groupToRemove []uint
		var i int
		for k, exGroup := range exGroups {
			if i >= len(f.PermitGroups) {
				groupToRemove = append(groupToRemove, exGroups[k:]...)
				break
			}
			for j := i; j < len(f.PermitGroups); j++ {
				if exGroup == f.PermitGroups[j] {
					i = j + 1
					break
				}
				i = j
				if exGroup < f.PermitGroups[j] {
					groupToRemove = append(groupToRemove, exGroup)
					break
				}
			}
		}

		var groupToCreate = make([]uint, len(f.PermitGroups)-len(exGroups)+len(groupToRemove))
		i = 0
		k := 0
		for _, group := range f.PermitGroups {
			if i < len(exGroups) {
				for j := i; j < len(exGroups); j++ {
					if group == exGroups[j] {
						i = j + 1
						goto nextGroup
					}
					i = j
					if group < exGroups[j] {
						break
					}
				}
			}
			groupToCreate[k] = group
			k++
		nextGroup:
		}

		if len(groupToRemove) > 0 {
			if err = appGroupSrv.UnBindForApp(app.ID, groupToRemove); err != nil {
				callback.Error(c, callback.ErrDBOperation, err)
				return
			}
			groupPermisstionChanged = true
		}
		if len(groupToCreate) > 0 {
			if _, err = appGroupSrv.BindForApp(app.ID, groupToCreate); err != nil {
				callback.Error(c, callback.ErrDBOperation, err)
				return
			}
		}
	}

	if appInfoChanged {
		if err = appSrv.UpdateAll(app.ID, f.Name, f.Callback, f.PermitAll); err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}
	}

	if groupPermisstionChanged {
		// 撤销 token

		loginRecordSrv := service.LoginRecordSrv{DB: appSrv.DB}
		records, err := loginRecordSrv.GetValidForApp(f.ID, daoUtil.LockForUpdate)
		if err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()
		ids := make([]uint, len(records))
		for i, record := range records {
			ids[i] = record.ID
			err = redis.CancelToken(ctx, uint64(record.ID), time.Unix(int64(record.ValidBefore), 0))
			if err != nil {
				if errors.Is(err, redis.Nil) {
					continue
				}
				callback.Error(c, callback.ErrUnexpected, err)
				return
			}
		}

		err = loginRecordSrv.SetDestroyedByIDS(ids)
		if err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}
	}

	if err = appSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Default(c)
}

func UpdateLinkState(c *gin.Context) {
	var f struct {
		ID      uint `json:"id" form:"id" binding:"required"`
		LinkOff bool `json:"linkOff" form:"linkOff"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	appSrv, err := service.App.Begin()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	defer appSrv.Rollback()

	uid := tools.GetUserInfo(c).UID

	if err = appSrv.UpdateLinkOff(uid, f.ID, f.LinkOff); err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	if err = appSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Default(c)
}
