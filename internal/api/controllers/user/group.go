package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
)

func ListGroups(c *gin.Context) {
	list, e := service.BaseGroups.LoadGroupsForShow()
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	callback.Success(c, list)
}
