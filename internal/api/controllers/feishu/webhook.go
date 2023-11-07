package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	log "github.com/sirupsen/logrus"
)

func Webhook(c *gin.Context) {
	var event feishu.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	if event.Type == "url_verification" {
		if event.Token != global.Config.Feishu.WebhookVerificationToken {
			callback.Error(c, callback.ErrUnauthorized, err)
			return
		}
		c.JSON(200, gin.H{
			"challenge": event.Challenge,
		})
		return
	} else {
		if event.Header.Token != global.Config.Feishu.WebhookVerificationToken {
			callback.Error(c, callback.ErrUnauthorized, err)
			return
		}
	}

	logger := log.WithFields(log.Fields{
		"n":    "飞书事件",
		"type": event.Header.EventType,
		"id":   event.Header.EventID,
	})
	logger.Infoln("Received")

	switch event.Header.EventType {
	case "contact.user.deleted_v3":
		var info feishu.UserDeletedEvent
		err := json.Unmarshal(event.Event, &info)
		if err != nil {
			callback.Error(c, callback.ErrForm, err)
			return
		}
		if info.Object.Mobile == "" {
			logger.Errorln("电话号码为空")
		} else {
			userSrv, err := service.User.Begin()
			if err != nil {
				callback.Error(c, callback.ErrDBOperation, err)
				return
			}
			defer userSrv.Rollback()
			result := userSrv.FrozeByPhone(info.Object.Mobile)
			if result.Error != nil || userSrv.Commit().Error != nil {
				callback.Error(c, callback.ErrDBOperation, result.Error)
				return
			}
			if result.RowsAffected != 0 {
				logger.Infof("%s:%s 离职已写入", info.Object.Name, info.Object.Mobile)
			}
		}
	default:
		logger.Warnf("未知的事件类型")
	}

	c.JSON(200, gin.H{})
}
