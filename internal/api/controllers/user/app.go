package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"gorm.io/gorm"
	"net/url"
)

func ApplyApp(c *gin.Context) {
	var f struct {
		Name         string `json:"name" form:"name" binding:"required,max=20"`
		Callback     string `json:"callback" form:"callback" binding:"url,required"`
		PermitAll    bool   `json:"permitAll" form:"permitAll"`
		PermitGroups []uint `json:"permitGroups" form:"permitGroups"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	}

	callbackUrl, e := url.Parse(f.Callback)
	if e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	} else if callbackUrl.Scheme != "https" {
		callback.Error(c, nil, callback.ErrForm)
		return
	}

	exist, e := service.SiteWhiteList.Exist(callbackUrl.Host)
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	} else if !exist {
		callback.Error(c, e, callback.ErrSiteNotAllow)
		return
	}

	exist, e = service.App.NameExist(f.Name)
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	} else if exist {
		callback.ErrorWithTip(c, nil, callback.ErrAlreadyExist, "名称已被占用")
		return
	}

	appSrc, e := service.App.Begin()
	if e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	}
	defer appSrc.Rollback()

	uid := tools.GetUserInfo(c).ID
	newApp, e := appSrc.New(uid, f.Name, f.Callback, f.PermitAll)
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	if !f.PermitAll && len(f.PermitGroups) != 0 {
		appGroupSrv := service.AppGroupSrv{DB: appSrc.DB}

		if e = appGroupSrv.BindForApp(newApp.ID, f.PermitGroups); e != nil {
			if e == gorm.ErrRecordNotFound {
				callback.Error(c, nil, callback.ErrGroupNotFound)
				return
			}
			callback.Error(c, e, callback.ErrDBOperation)
			return
		}
	}

	if e = redis.AppCode.Add(newApp.AppCode); e != nil || appSrc.Commit().Error != nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	}

	callback.Success(c, dto.AppNew{
		AppShow: dto.AppShow{
			ID:             newApp.ID,
			Name:           newApp.Name,
			AppCode:        newApp.AppCode,
			PermitAllGroup: newApp.PermitAllGroup,
		},
		AppSecret: newApp.AppSecret,
	})
}

func ListOwnedApp(c *gin.Context) {
	uid := tools.GetUserInfo(c).ID

	apps, e := service.App.GetUserOwnedApp(uid)
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	callback.Success(c, apps)
}
