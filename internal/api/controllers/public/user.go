package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"strconv"
	"strings"
)

func GetUserPublicInfo(c *gin.Context) {
	var f struct {
		ID string `json:"id" form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	idStrArr := strings.Split(f.ID, ",")
	idNums := make([]uint, len(idStrArr))
	for i, str := range idStrArr {
		id, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			callback.Error(c, callback.ErrForm, err)
			return
		}
		idNums[i] = uint(id)
	}

	data, err := service.User.GetUserInfoPublic(idNums...)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, data)
}
