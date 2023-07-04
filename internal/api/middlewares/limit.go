package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
)

func LimitGroup(groups ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userGroups := tools.GetUserInfo(c).Groups
		for _, userGroup := range userGroups {
			for _, allowGroup := range groups {
				if userGroup == allowGroup {
					return
				}
			}
		}

		callback.Error(c, callback.ErrInsufficientPermissions)
	}
}
