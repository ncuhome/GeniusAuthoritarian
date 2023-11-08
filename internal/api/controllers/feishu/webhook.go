package controllers

import (
	"encoding/json"
	"errors"
	"github.com/Mmx233/daoUtil"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishuApi"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Webhook(c *gin.Context) {
	var event feishuApi.Event
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
		var info feishuApi.UserDeletedEvent
		err := json.Unmarshal(event.Event, &info)
		if err != nil {
			callback.Error(c, callback.ErrForm, err)
			return
		}
		if info.Object.Mobile == "" {
			logger.Errorln("电话号码为空")
		} else {
			logger = logger.WithFields(log.Fields{
				"name":  info.Object.Name,
				"phone": info.Object.Mobile,
			})
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
				logger.Infoln("离职已写入")
			}
		}
	case "contact.user.updated_v3":
		var info feishuApi.UserUpdatedEvent
		err := json.Unmarshal(event.Event, &info)
		if err != nil {
			callback.Error(c, callback.ErrForm, err)
			return
		}
		if info.Object.Mobile == "" {
			return
		}
		user := feishu.NewUser(&info.Object)
		userSrv, err := service.User.Begin()
		if err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}
		defer userSrv.Rollback()
		logger = logger.WithFields(log.Fields{
			"name":  user.Data.Name,
			"phone": user.Data.Mobile,
		})
		if user.IsInvalid() {
			result := userSrv.FrozeByPhone(user.Data.Mobile)
			if result.Error != nil {
				callback.Error(c, callback.ErrDBOperation, result.Error)
				return
			}
			if result.RowsAffected != 0 {
				logger.Infoln("已冻结")
			}
		} else {
			userModel, err := userSrv.FirstByPhone(info.OldObject.Mobile, daoUtil.UnScoped, daoUtil.LockForUpdate)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = userSrv.CreateAll([]dao.User{user.Model()})
				if err != nil {
					callback.Error(c, callback.ErrDBOperation, err)
					return
				}
			} else if err != nil {
				callback.Error(c, callback.ErrDBOperation, err)
				return
			} else if userModel.DeletedAt.Valid {
				err = userSrv.UnFrozeByIDSlice([]uint{userModel.ID})
				if err != nil {
					callback.Error(c, callback.ErrDBOperation, err)
					return
				}
			} else {
				userModelNew := user.Model()
				userModelNew.ID = userModel.ID
				if err = userModelNew.UpdateAllInfoByID(userSrv.DB).Error; err != nil {
					callback.Error(c, callback.ErrDBOperation, err)
					return
				}
				logger.Infoln("信息已更新")

				groupMap, err := (&service.FeishuGroupsSrv{DB: userSrv.DB}).GetGroupMap(daoUtil.LockForShare)
				if err != nil {
					callback.Error(c, callback.ErrDBOperation, err)
					return
				}

				changed, err := user.SyncDepartments(userSrv.DB, userModel.ID, groupMap)
				if err != nil {
					callback.Error(c, callback.ErrDBOperation, err)
					return
				}
				if changed {
					logger.Infoln("部门已同步")
				}
			}
		}
		if err = userSrv.Commit().Error; err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}
	default:
		logger.Warnf("未知的事件类型")
	}

	c.JSON(200, gin.H{})
}
