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
		userDeleted(c, logger, event.Event)
	case "contact.user.updated_v3":
		userUpdated(c, logger, event.Event)
	default:
		logger.Warnf("未知的事件类型")
	}
	if c.IsAborted() {
		return
	}

	c.JSON(200, gin.H{})
}

func userDeleted(c *gin.Context, logger *log.Entry, event json.RawMessage) {
	var info feishuApi.UserDeletedEvent
	err := json.Unmarshal(event, &info)
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
}

func userUpdated(c *gin.Context, logger *log.Entry, event json.RawMessage) {
	var info feishuApi.UserUpdatedEvent
	err := json.Unmarshal(event, &info)
	if err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}
	if info.Object.Mobile == "" {
		return
	}
	user := feishu.NewUser(&info.Object)
	oldUser := feishu.NewUser(&info.OldObject)
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
			logger.Infoln("用户已冻结")
		}
	} else {
		var phone string
		if oldUser.Data.Mobile != "" {
			phone = oldUser.Data.Mobile
		} else {
			phone = user.Data.Mobile
		}
		userModel, err := userSrv.FirstByPhone(phone, daoUtil.UnScoped, daoUtil.LockForUpdate)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			models := []dao.User{user.Model()}
			err = userSrv.CreateAll(models)
			if err != nil {
				callback.Error(c, callback.ErrDBOperation, err)
				return
			}
			groupMap, err := (&service.FeishuGroupsSrv{DB: userSrv.DB}).GetGroupMap(daoUtil.LockForShare)
			if err != nil {
				callback.Error(c, callback.ErrDBOperation, err)
				return
			}
			err = service.UserGroupsSrv{DB: userSrv.DB}.CreateAll(user.DepartmentModels(models[0].ID, groupMap))
			if err != nil {
				callback.Error(c, callback.ErrDBOperation, err)
				return
			}
			logger.Infoln("用户已追加创建")
		} else if err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		} else {
			if userModel.DeletedAt.Valid {
				err = userSrv.UnFrozeByIDSlice([]uint{userModel.ID})
				if err != nil {
					callback.Error(c, callback.ErrDBOperation, err)
					return
				}
				logger.Infoln("用户已解冻")
			}

			if userModel.DeletedAt.Valid || !oldUser.IsModelEmpty() {
				userModelNew := user.Model()
				userModelNew.ID = userModel.ID
				if err = userModelNew.UpdateAllInfoByID(userSrv.DB).Error; err != nil {
					callback.Error(c, callback.ErrDBOperation, err)
					return
				}
				logger.Infoln("用户信息已更新")
			}

			if userModel.DeletedAt.Valid || len(oldUser.Data.DepartmentIds) != 0 {
				groupMap, err := (&service.FeishuGroupsSrv{DB: userSrv.DB}).GetGroupMap(daoUtil.LockForShare)
				if err != nil {
					callback.Error(c, callback.ErrDBOperation, err)
					return
				}

				err = user.SyncDepartments(userSrv.DB, userModel.ID, groupMap)
				if err != nil {
					callback.Error(c, callback.ErrDBOperation, err)
					return
				}
				logger.Infoln("用户部门已同步")
			}
		}
	}
	if err = userSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
}
