package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/models/response"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	log "github.com/sirupsen/logrus"
)

func VerifyToken(c *gin.Context) {
	var f struct {
		Token  string   `json:"token" form:"token" binding:"required"`
		Groups []string `json:"groups" form:"groups"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm)
		return
	}

	claims, valid, e := jwt.ParseLoginToken(f.Token)
	if e != nil || !valid {
		log.Debugln("decode token failed:", e)
		callback.Error(c, callback.ErrUnauthorized)
		return
	}

	if len(f.Groups) != 0 {
		var verifiedGroups []string
		for _, targetGroup := range f.Groups {
			for _, existGroup := range claims.Groups {
				if targetGroup == existGroup {
					verifiedGroups = append(verifiedGroups, existGroup)
				}
			}
		}
		if len(verifiedGroups) == 0 {
			callback.Error(c, callback.ErrUnauthorized)
			return
		}
		claims.Groups = verifiedGroups
	}

	loginRecordSrv, e := service.LoginRecord.Begin()
	if e != nil {
		callback.Error(c, callback.ErrDBOperation)
		return
	}
	defer loginRecordSrv.Rollback()

	if e = loginRecordSrv.Add(claims.UID, claims.Target); e != nil || loginRecordSrv.Commit().Error != nil {
		callback.Error(c, callback.ErrDBOperation)
		return
	}

	callback.Success(c, response.VerifyTokenSuccess{
		Name:   claims.Name,
		Groups: claims.Groups,
	})
}
