package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
)

func FeishuLoginLink(c *gin.Context) {
	var f struct {
		Callback string `json:"callback" form:"callback" binding:"required,uri"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm)
		return
	}

	if ok, e := service.SiteWhiteList.CheckUrl(f.Callback); e != nil {
		callback.Error(c, callback.ErrDBOperation)
		return
	} else if !ok {
		callback.Error(c, callback.ErrSiteNotAllow)
		return
	}

	callback.Success(c, gin.H{
		"url": feishu.Api.LoginLink(c.Request.Host, f.Callback),
	})
}

func FeishuLogin(c *gin.Context) {
	var f struct {
		Code string `json:"code" form:"code" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm)
		return
	}

	user, e := feishu.Api.GetUser(f.Code)
	if e != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed)
		return
	}

	userInfo, e := user.Info()
	if e != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed)
		return
	}

	groups := feishu.Departments.MultiSearch(userInfo.User.DepartmentIds)
	if len(groups) == 0 {
		callback.Error(c, callback.ErrFindUnit)
		return
	}

	token, e := jwt.GenerateAuthToken()
	if e != nil {
		callback.Error(c, callback.ErrUnexpected)
		return
	}

	callback.Success(c, gin.H{
		"token": token,
	})
}
