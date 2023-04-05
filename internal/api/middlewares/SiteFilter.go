package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"net/url"
)

func SiteFilter(c *gin.Context) {
	if c.Request.Method == "GET" {
		return
	}

	ref, ok := c.GetQuery("referer")
	if !ok || ref == "" {
		callback.ErrorWithTip(c, callback.ErrSiteNotAllow, "授权参数缺失")
		return
	}

	u, e := url.Parse(ref)
	if e != nil {
		callback.ErrorWithTip(c, callback.ErrSiteNotAllow, "referer 不合法")
		return
	}

	ok, e = service.SiteWhiteList.Exist(u.Host)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation)
		return
	}

	if !ok {
		callback.Error(c, callback.ErrSiteNotAllow)
	}
}
