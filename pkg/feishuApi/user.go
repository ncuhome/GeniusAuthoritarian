package feishuApi

import (
	"fmt"
	"github.com/Mmx233/tool"
)

type FsUser struct {
	fs Fs

	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int64  `json:"expires_in"`
	Name             string `json:"name"`
	AvatarUrl        string `json:"avatar_url"`
	AvatarThumb      string `json:"avatar_thumb"`
	AvatarMiddle     string `json:"avatar_middle"`
	AvatarBig        string `json:"avatar_big"`
	OpenId           string `json:"open_id"`
	UnionId          string `json:"union_id"`
	Email            string `json:"email"`
	EnterpriseEmail  string `json:"enterprise_email"`
	UserId           string `json:"user_id"`
	Mobile           string `json:"mobile"`
	TenantKey        string `json:"tenant_key"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	Sid              string `json:"sid"`
}

func (u FsUser) Info() (*UserInfoResp, error) {
	var data UserInfoResp
	return &data, u.fs.doRequest("GET", &data, &tool.DoHttpReq{
		Url: fmt.Sprintf(GetUserInfoUrl, u.OpenId),
		Header: map[string]interface{}{
			"Authorization": "Bearer " + u.AccessToken,
		},
	})
}
