package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/cookie"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
)

func GoFeishuLogin(c *gin.Context) {
	var f struct {
		Callback string `json:"callback" form:"callback" binding:"required,uri"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm)
		return
	}

	c.Redirect(302, feishu.Api.LoginLink(f.Callback))
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

	refreshToken, e := jwt.GenerateRefreshToken(user.Name, groups)
	if e != nil {
		callback.Error(c, callback.ErrUnexpected)
		return
	}

	cookie.SetRefreshToken(c, refreshToken)
	callback.Default(c)
}
