package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Mmx233/daoUtil"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/sshDev/sshDevClient/sshTool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/sshDev/sshDevModel"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishuApi"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/rand"
	"time"
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

// 逻辑太长了，抽象不了，受不了
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
	var phone string
	// 处理电话号码更换的情况
	if oldUser.Data.Mobile != "" {
		phone = oldUser.Data.Mobile
	} else {
		phone = user.Data.Mobile
	}
	userModel, err := userSrv.FirstByPhone(phone, daoUtil.UnScoped, daoUtil.LockForUpdate)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if !user.IsInvalid() {
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
				err = service.UserGroupsSrv{DB: userSrv.DB}.CreateAll(user.Departments(groupMap).Models(models[0].ID))
				if err != nil {
					callback.Error(c, callback.ErrDBOperation, err)
					return
				}
				logger.Infoln("用户已追加创建")
				return
			} else {
				// 数据库中不存在且用户无效
				return
			}
		} else {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}
	}

	redisUserOperator := redis.NewUserJwt().NewOperator(userModel.ID)

	if user.IsInvalid() != userModel.DeletedAt.Valid {
		if userModel.DeletedAt.Valid { // 用户有效但被冻结，执行解冻
			err = userSrv.UnFrozeByIDSlice([]uint{userModel.ID})
			if err != nil {
				callback.Error(c, callback.ErrDBOperation, err)
				return
			}
			logger.Infoln("用户已解冻")
		} else { // 用户有效但未冻结，执行冻结
			result := userSrv.FrozeByPhone(user.Data.Mobile)
			if result.Error != nil {
				callback.Error(c, callback.ErrDBOperation, result.Error)
				return
			}
			err = redisUserOperator.Del(context.Background())
			if err != nil {
				callback.Error(c, callback.ErrUnexpected, err)
				return
			}
			if result.RowsAffected != 0 {
				logger.Infoln("用户已冻结")
			}
			return
		}
	}

	if userModel.DeletedAt.Valid || !oldUser.IsModelEmpty() { // 刚刚执行了用户解冻或有效字段变更时，更新数据库用户信息字段
		userModelNew := user.Model()
		userModelNew.ID = userModel.ID
		if err = userModelNew.UpdateAllInfoByID(userSrv.DB).Error; err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}
		logger.Infoln("用户信息已更新")
	}

	if userModel.DeletedAt.Valid || len(oldUser.Data.DepartmentIds) != 0 { // 刚刚执行了用户解冻或部门变动时同步用户部门
		groupMap, err := (&service.FeishuGroupsSrv{DB: userSrv.DB}).GetGroupMap(daoUtil.LockForShare)
		if err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}

		userGroupSrv := service.UserGroupsSrv{DB: userSrv.DB}

		// 同步前是否是研发组成员
		prevIsDeveloper, err := userGroupSrv.IsUnitMember(userModel.ID, departments.UDev)
		if err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}

		// 执行部门同步,写入数据库
		err = user.Departments(groupMap).Sync(userSrv.DB, userModel.ID)
		if err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}

		_, err = redisUserOperator.ChangeOperateID(context.Background())
		if err != nil {
			callback.Error(c, callback.ErrUnexpected, err)
			return
		}

		logger.Infoln("用户部门已同步")

		// 同步后是否是研发组成员
		nowIsDeveloper, err := userGroupSrv.IsUnitMember(userModel.ID, departments.UDev)
		if err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}

		// 处理研发身份变更,同步 SSH 权限
		if (prevIsDeveloper || nowIsDeveloper) && prevIsDeveloper != nowIsDeveloper {
			var msg sshDevModel.SshAccountMsg
			userName := sshTool.LinuxAccountName(userModel.ID)
			if prevIsDeveloper { // 以前是现在不是
				msg = sshDevModel.SshAccountMsg{
					IsDel:    true,
					Username: userName,
				}
			} else { // 以前不是现在是
				model, err := sshTool.NewSshDevModel(rand.New(rand.NewSource(time.Now().UnixNano())), userModel.ID)
				if err != nil {
					callback.Error(c, callback.ErrUnexpected, err)
					return
				}

				err = service.UserSshSrv{DB: userSrv.DB}.CreateAll([]dao.UserSsh{model})
				if err != nil {
					callback.Error(c, callback.ErrDBOperation, err)
					return
				}

				msg = sshDevModel.SshAccountMsg{
					Username:  userName,
					PublicKey: model.PublicSsh,
				}
			}
			err = redis.PublishSshDev([]sshDevModel.SshAccountMsg{msg})
			if err != nil {
				callback.Error(c, callback.ErrDBOperation, err)
				return
			}
			logger.Infoln("研发 SSH 已同步")
		}
	}

	if err = userSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
}
