package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/util"
)

func GoFeishuLogin(c *gin.Context) {
	var f struct {
		Callback string `json:"callback" form:"callback" binding:"required,uri"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm)
		return
	}

	c.Redirect(302, util.Feishu.LoginLink(f.Callback))
}

func FeishuLogin(c *gin.Context) {
	var f struct {
		Code string `json:"code" form:"code" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm)
		return
	}

	_, e := util.Feishu.GetUser(f.Code)
	if e != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed)
		return
	}

	/*token, e := auth.Jwt.GenerateToken(time.Hour * 24 * 7)
	if e != nil {
		callback.Error(c, msg.InnerErr, e)
		return
	}

	callback.Success(c, token)*/
}
