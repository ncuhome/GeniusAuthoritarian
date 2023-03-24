package feishu

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mmx233/tool"
	"net/http"
	"net/url"
)

func New(ClientID, Secret string, client *http.Client) *FsLogin {
	return &FsLogin{
		appID:  ClientID,
		secret: Secret,
		Http:   tool.NewHttpTool(client),
		tenant: newTenant(),
	}
}

type FsLogin struct {
	appID  string
	secret string
	Http   *tool.Http

	tenant *tenantTokenCache
}

func (f FsLogin) loadTenantAccessToken() (string, error) {
	if token, exist := f.tenant.Load(); exist {
		return token, nil
	}

	res, e := f.Http.PostRequest(&tool.DoHttpReq{
		Url: "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal",
		Body: map[string]string{
			"app_id":     f.appID,
			"app_secret": f.secret,
		},
	})
	if e != nil {
		return "", e
	}
	defer res.Body.Close()

	var data TenantAccessTokenResp
	if e = json.NewDecoder(res.Body).Decode(&data); e != nil {
		return "", e
	} else if data.Code != 0 {
		return "", errors.New(data.Msg)
	}

	if res.StatusCode > 299 {
		return "", fmt.Errorf("server return http status %d", res.StatusCode)
	}

	f.tenant.Set(data.TenantAccessToken, data.Expire)

	return data.TenantAccessToken, nil
}

func (f FsLogin) LoginLink(state string) string {
	return fmt.Sprintf(
		"https://open.feishu.cn/open-apis/authen/v1/user_auth_page_beta?app_id=%s&redirect_uri=https%%3A%%2F%%2Fv.ncuos.com%%2Fapi%%2Fpublic%%2Flogin%%2Ffeishu%%2F&state=%s",
		f.appID,
		url.QueryEscape(state),
	)
}

func (f FsLogin) GetUser(code string) (*FsUser, error) {
	tenantToken, e := f.loadTenantAccessToken()
	if e != nil {
		return nil, e
	}

	res, e := f.Http.PostRequest(&tool.DoHttpReq{
		Url: "https://open.feishu.cn/open-apis/authen/v1/access_token",
		Header: map[string]interface{}{
			"Authorization": "Bearer " + tenantToken,
			"Content-Type":  "application/json",
		},
		Body: map[string]string{
			"grant_type": "authorization_code",
			"code":       code,
		},
	})
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()

	var data OAuth2AccessTokenResp
	if e = json.NewDecoder(res.Body).Decode(&data); e != nil {
		return nil, e
	} else if data.Code != 0 {
		return nil, errors.New(data.Msg)
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("server return http status %d", res.StatusCode)
	}

	data.Data.Http = f.Http
	return &data.Data, nil
}

type FsUser struct {
	Http             *tool.Http `json:"-"`
	AccessToken      string     `json:"access_token"`
	TokenType        string     `json:"token_type"`
	ExpiresIn        int64      `json:"expires_in"`
	Name             string     `json:"name"`
	AvatarUrl        string     `json:"avatar_url"`
	AvatarThumb      string     `json:"avatar_thumb"`
	AvatarMiddle     string     `json:"avatar_middle"`
	AvatarBig        string     `json:"avatar_big"`
	OpenId           string     `json:"open_id"`
	UnionId          string     `json:"union_id"`
	Email            string     `json:"email"`
	EnterpriseEmail  string     `json:"enterprise_email"`
	UserId           string     `json:"user_id"`
	Mobile           string     `json:"mobile"`
	TenantKey        string     `json:"tenant_key"`
	RefreshExpiresIn int        `json:"refresh_expires_in"`
	RefreshToken     string     `json:"refresh_token"`
	Sid              string     `json:"sid"`
}

func (u FsUser) GetFeishuUserInfo() (*UserInfoResp, error) {
	res, e := u.Http.GetRequest(&tool.DoHttpReq{
		Url: "https://open.feishu.cn/open-apis/contact/v3/users/" + u.OpenId,
		Header: map[string]interface{}{
			"Authorization": "Bearer " + u.AccessToken,
		},
	})
	if e != nil {
		return nil, e
	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("server return http status %d", res.StatusCode)
	}

	var data UserInfoResp
	return &data, json.NewDecoder(res.Body).Decode(&data)
}
