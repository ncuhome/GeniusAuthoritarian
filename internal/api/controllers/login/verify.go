package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/models/response"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
)

func doVerifyToken(c *gin.Context, token string, groups []string) *jwt.LoginTokenClaims {
	claims, valid, e := jwt.ParseLoginToken(token)
	if e != nil || !valid {
		callback.Error(c, e, callback.ErrUnauthorized)
		return nil
	}

	if len(groups) != 0 {
		var verifiedGroups []string
		for _, targetGroup := range groups {
			for _, existGroup := range claims.Groups {
				if targetGroup == existGroup {
					verifiedGroups = append(verifiedGroups, existGroup)
				}
			}
		}
		if len(verifiedGroups) == 0 {
			callback.Error(c, nil, callback.ErrUnauthorized)
			return nil
		}
		claims.Groups = verifiedGroups
	}

	loginRecordSrv, e := service.LoginRecord.Begin()
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return nil
	}
	defer loginRecordSrv.Rollback()

	if e = loginRecordSrv.Add(claims.UID, claims.IP, claims.Target); e != nil || loginRecordSrv.Commit().Error != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return nil
	}

	return claims
}

func VerifyToken(c *gin.Context) {
	var f struct {
		Token  string   `json:"token" form:"token" binding:"required"`
		Groups []string `json:"groups" form:"groups"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	}

	claims := doVerifyToken(c, f.Token, f.Groups)
	if c.IsAborted() {
		return
	}

	callback.Success(c, response.VerifyTokenSuccess{
		Name:   claims.Name,
		Groups: claims.Groups,
	})
}

func Login(c *gin.Context) {
	var f struct {
		Token string `json:"token" form:"token" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	}

	claims := doVerifyToken(c, f.Token, nil)
	if c.IsAborted() {
		return
	}

	token, e := jwt.GenerateUserToken(claims.UID)
	if e != nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	}

	callback.Success(c, gin.H{
		"token": token,
	})
}
