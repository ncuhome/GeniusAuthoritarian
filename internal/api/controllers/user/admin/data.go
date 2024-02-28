package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	log "github.com/sirupsen/logrus"
	"time"
)

func LoginData(c *gin.Context) {
	var f struct {
		Range string `json:"range" form:"range" binding:"eq=week|eq=month|eq=year"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	startTime := time.Now()
	switch f.Range {
	default:
		log.Warnln("admin login data controller encountered unhandled form parameters")
		fallthrough
	case "week":
		startTime.Add(-time.Hour * 24 * 7)
	case "month":
		startTime.Add(-time.Hour * 24 * 30)
	case "year":
		startTime.Add(-time.Hour * 24 * 365)
	}

	data, err := service.LoginRecord.GetForAdminView(startTime)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, data)
}
