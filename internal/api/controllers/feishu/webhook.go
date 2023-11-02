package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	log "github.com/sirupsen/logrus"
)

func Webhook(c *gin.Context) {
	var event feishu.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	if event.Header.Token != global.Config.Feishu.WebhookVerificationToken {
		callback.Error(c, callback.ErrUnauthorized, err)
		return
	}

	log.Infof("[飞书事件] %s:%s", event.Header.EventType, event.Header.EventID)

}
