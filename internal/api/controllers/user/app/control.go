package controllers

import (
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
)

func checkCallback(c *gin.Context, callbackStr string) {
	callbackUrl, e := url.Parse(callbackStr)
	if e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	} else if callbackUrl.Scheme != "https" {
		callback.Error(c, callback.ErrForm)
		return
	}

	exist, e := service.SiteWhiteList.Exist(callbackUrl.Host)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	} else if !exist {
		callback.Error(c, callback.ErrSiteNotAllow, e)
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
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	}

	checkCallback(c, f.Callback)
	if c.IsAborted() {
		return
	}

	appSrc, e := service.App.Begin()
	if e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	}
	defer appSrc.Rollback()

	uid := tools.GetUserInfo(c).ID
	newApp, e := appSrc.New(uid, f.Name, f.Callback, f.PermitAll)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	var groups = make([]dto.Group, 0)
	if !f.PermitAll && len(f.PermitGroups) != 0 {
		appGroupSrv := service.AppGroupSrv{DB: appSrc.DB}

		if groups, e = appGroupSrv.BindForApp(newApp.ID, f.PermitGroups); e != nil {
			if e == gorm.ErrRecordNotFound {
				callback.Error(c, callback.ErrGroupNotFound)
				return
			}
			callback.Error(c, callback.ErrDBOperation, e)
			return
		}
	}

	if e = redis.AppCode.Add(newApp.AppCode); e != nil || appSrc.Commit().Error != nil {
		callback.Error(c, callback.ErrUnexpected, e)
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
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	}

	uid := tools.GetUserInfo(c).ID

	appSrv, e := service.App.Begin()
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}
	defer appSrv.Rollback()

	appModel, e := appSrv.DeleteByID(f.ID, uid)
	if e != nil {
		if e == gorm.ErrRecordNotFound {
			callback.Error(c, callback.ErrAppNotFound)
			return
		}
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	if e = redis.AppCode.Del(appModel.AppCode); e != nil {
		callback.Error(c, callback.ErrUnexpected)
		return
	}

	if e = appSrv.Commit().Error; e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
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
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm, e)
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

	appSrv, e := service.App.Begin()
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}
	defer appSrv.Rollback()

	uid := tools.GetUserInfo(c).ID

	app, e := appSrv.FirstAppDetailedByIDForUser(f.ID, uid, daoUtil.LockForUpdate)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	appInfoChanged := false

	if f.Name != app.Name {
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
	if app.PermitAllGroup != f.PermitAll {
		if f.PermitAll {
			e = appGroupSrv.DeleteAllForApp(app.ID)
			if e != nil {
				callback.Error(c, callback.ErrDBOperation, e)
				return
			}
		} else if len(f.PermitGroups) > 0 {
			app.Groups, e = appGroupSrv.BindForApp(app.ID, f.PermitGroups)
			if e != nil {
				callback.Error(c, callback.ErrDBOperation, e)
				return
			}
		}
		appInfoChanged = true
	} else if !f.PermitAll {
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
			if e = appGroupSrv.UnBindForApp(app.ID, groupToRemove); e != nil {
				callback.Error(c, callback.ErrDBOperation, e)
				return
			}
		}
		if len(groupToCreate) > 0 {
			if _, e = appGroupSrv.BindForApp(app.ID, groupToCreate); e != nil {
				callback.Error(c, callback.ErrDBOperation, e)
				return
			}
		}
	}

	if appInfoChanged {
		if e = appSrv.UpdateAll(app.ID, f.Name, f.Callback, f.PermitAll); e != nil {
			callback.Error(c, callback.ErrDBOperation, e)
			return
		}
	}

	if e = appSrv.Commit().Error; e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	callback.Default(c)
}

func UpdateLinkState(c *gin.Context) {
	var f struct {
		ID      uint `json:"id" form:"id" binding:"required"`
		LinkOff bool `json:"linkOff" form:"linkOff"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	}

	appSrv, e := service.App.Begin()
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}
	defer appSrv.Rollback()

	uid := tools.GetUserInfo(c).ID

	if e = appSrv.UpdateLinkOff(uid, f.ID, f.LinkOff); e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	if e = appSrv.Commit().Error; e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	callback.Default(c)
}
