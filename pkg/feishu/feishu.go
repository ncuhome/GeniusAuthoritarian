package feishu

import (
	"encoding/json"
	"fmt"
	"github.com/Mmx233/tool"
	"net/http"
)

func New(ClientID, Secret string, client *http.Client) *FsLogin {
	return &FsLogin{
		ClientID: ClientID,
		Secret:   Secret,
		Http:     tool.NewHttpTool(client),
	}
}

type FsLogin struct {
	ClientID string
	Secret   string
	Http     *tool.Http
}

func (f FsLogin) tenantAccessToken() (string, error) {
	res, e := f.Http.PostRequest(&tool.DoHttpReq{
		Url: "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal",
		Body: map[string]string{
			"app_id":     f.ClientID,
			"app_secret": f.Secret,
		},
	})
	if e != nil {
		return "", e
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return "", fmt.Errorf("server return http status %d", res.StatusCode)
	}

	var data TenantAccessTokenResp
	if e = json.NewDecoder(res.Body).Decode(&data); e != nil {
		return "", e
	}

	if data.Code != 0 {
		return "", fmt.Errorf("server return code %d", res.StatusCode)
	}

	return data.TenantAccessToken, nil
}

func (f FsLogin) GetUserAccessToken(code string) (*OAuth2AccessTokenResp, bool, error) {
	tenantToken, e := f.tenantAccessToken()
	if e != nil {
		return nil, false, e
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
		return nil, false, e
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, false, fmt.Errorf("server return http status %d", res.StatusCode)
	}

	var data OAuth2AccessTokenResp
	if e = json.NewDecoder(res.Body).Decode(&data); e != nil {
		return nil, false, e
	}

	return &data, data.Code == 0, nil
}

func (f FsLogin) GetFeishuUserInfo(userAccessToken string) (*UserInfo, error) {
	res, e := f.Http.GetRequest(&tool.DoHttpReq{
		Url: "https://open.feishu.cn/connect/qrconnect/oauth2/user_info/",
		Header: map[string]interface{}{
			"Authorization": "Bearer " + userAccessToken,
		},
		Query:  nil,
		Body:   nil,
		Cookie: nil,
	})
	if e != nil {
		return nil, e
	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("server return http status %d", res.StatusCode)
	}

	var data UserInfo
	return &data, json.NewDecoder(res.Body).Decode(&data)
}
