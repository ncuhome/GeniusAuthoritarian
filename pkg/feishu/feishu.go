package feishu

import (
	"encoding/json"
	"fmt"
	"github.com/Mmx233/tool"
	"net/http"
	"net/url"
)

func New(ClientID, Secret string, client *http.Client) *Fs {
	return &Fs{
		appID:  ClientID,
		secret: Secret,
		Http:   tool.NewHttpTool(client),
		tenant: newTenant(),
	}
}

// Fs 飞书缩写
type Fs struct {
	appID  string
	secret string
	Http   *tool.Http

	tenant *tenantTokenCache
}

func (f Fs) doRequest(Method string, data interface{}, opt *tool.DoHttpReq) error {
	res, e := f.Http.Request(Method, opt)
	if e != nil {
		return e
	}
	defer res.Body.Close()

	container := &Resp{
		Data: data,
	}
	if e = json.NewDecoder(res.Body).Decode(container); e != nil {
		return e
	} else if container.Code != 0 {
		return fmt.Errorf("feishu api: %s", container.Msg)
	}

	if res.StatusCode > 299 {
		return fmt.Errorf("server return http status %d", res.StatusCode)
	}

	return nil
}

func (f Fs) loadTenantAccessToken() (string, error) {
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
		return "", fmt.Errorf("feishu api: %s", data.Msg)
	}

	if res.StatusCode > 299 {
		return "", fmt.Errorf("server return http status %d", res.StatusCode)
	}

	f.tenant.Set(data.TenantAccessToken, data.Expire)

	return data.TenantAccessToken, nil
}

func (f Fs) LoginLink(selfDomain, state string) string {
	return fmt.Sprintf(
		"https://open.feishu.cn/open-apis/authen/v1/user_auth_page_beta?app_id=%s&redirect_uri=https%%3A%%2F%%2F"+selfDomain+"%%2Ffeishu%%2F&state=%s",
		f.appID,
		url.QueryEscape(state),
	)
}

func (f Fs) GetUser(code string) (*FsUser, error) {
	tenantToken, e := f.loadTenantAccessToken()
	if e != nil {
		return nil, e
	}
	var data FsUser
	data.fs = f
	return &data, f.doRequest("POST", &data, &tool.DoHttpReq{
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
}

func (f Fs) LoadDepartmentList() (map[string]string, error) {
	tenantToken, e := f.loadTenantAccessToken()
	if e != nil {
		return nil, e
	}
	var data ListDepartmentResp
	if e = f.doRequest("GET", &data, &tool.DoHttpReq{
		Url: "https://open.feishu.cn/open-apis/contact/v3/departments/0/children",
		Header: map[string]interface{}{
			"Authorization": "Bearer " + tenantToken,
			"Content-Type":  "application/json",
		},
		Query: map[string]interface{}{
			"page_size": 40,
		},
	}); e != nil {
		return nil, e
	}

	var m = make(map[string]string, len(data.Items))
	for _, v := range data.Items {
		m[v.OpenDepartmentId] = v.Name
	}
	return m, nil
}
