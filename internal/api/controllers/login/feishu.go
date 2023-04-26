package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
)

var FeishuLoginLink = GetLoginLink(feishu.Api.LoginLink)

var FeishuLogin = Login(func(c *gin.Context, code string) (string, []string) {
	user, e := feishu.Api.GetUser(code)
	if e != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed)
		return "", nil
	}

	userInfo, e := user.Info()
	if e != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed)
		return "", nil
	}

	groups, e := service.FeishuGroups.Search(userInfo.User.DepartmentIds)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation)
		return "", nil
	} else if len(groups) == 0 {
		callback.Error(c, callback.ErrFindUnit)
		return "", nil
	}

	var groupSlice = make([]string, len(groups))
	for i, group := range groups {
		groupSlice[i] = group.Name
	}

	return userInfo.User.Name, groupSlice
})
