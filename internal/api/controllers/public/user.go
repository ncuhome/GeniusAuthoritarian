package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
)

func GetUserPublicInfo(c *gin.Context) {
	var f struct {
		ID []uint `json:"id" form:"id" binding:"required"`
	}
	if err := tools.ShouldBindReused(c, &f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	data, err := service.User.GetUserInfoPublic(f.ID...)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, data)
}
