package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sms"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"strings"
)

func SendVerifySms(c *gin.Context) {
	uid := tools.GetUserInfo(c).UID
	phone, err := service.User.FirstPhoneByID(uid)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	phone = strings.TrimPrefix(phone, "+86")

	code, err := redis.NewUserIdentityCode(uid).New()
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	smsRedis := redis.NewSms(phone)
	ok, err := smsRedis.TryLock()
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	} else if !ok {
		callback.Error(c, callback.ErrSmsCoolDown)
		return
	}

	err = sms.Ums.Send("你的验证码为"+code, phone)
	if err != nil {
		callback.Error(c, callback.ErrSendSmsFailed, err)
		_ = smsRedis.UnLock()
		return
	}

	callback.Default(c)
}
