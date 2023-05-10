package feishuApi

import (
	"encoding/json"
	"fmt"
	"github.com/Mmx233/tool"
	"net/http"
	"net/url"
)

func New(ClientID, Secret string, client *http.Client) *Fs {
	fs := &Fs{
		appID:  ClientID,
		secret: Secret,
		Http:   tool.NewHttpTool(client),
	}
	fs.tenant = NewTenant(fs)
	return fs
}

// Fs 飞书缩写
type Fs struct {
	appID  string
	secret string
	Http   *tool.Http

	tenant *TenantToken
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

func (f Fs) GetTenantAccessToken() (*TenantAccessTokenResp, error) {
	res, e := f.Http.PostRequest(&tool.DoHttpReq{
		Url: GetTenantAccessTokenUrl,
		Body: map[string]string{
			"app_id":     f.appID,
			"app_secret": f.secret,
		},
	})
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()

	var data TenantAccessTokenResp
	if e = json.NewDecoder(res.Body).Decode(&data); e != nil {
		return nil, e
	} else if data.Code != 0 {
		return nil, fmt.Errorf("feishu api: %s", data.Msg)
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("server return http status %d", res.StatusCode)
	}
	return &data, nil
}

func (f Fs) LoginLink(selfDomain, state string) string {
	redirectUri := "https://" + selfDomain + "/feishu/"

	return fmt.Sprintf(LarkLoginUrl,
		f.appID,
		url.QueryEscape(state),
		url.QueryEscape(redirectUri))
}

func (f Fs) GetUser(code string) (*FsUser, error) {
	tenantToken, e := f.tenant.Load()
	if e != nil {
		return nil, e
	}
	var data FsUser
	data.fs = f
	return &data, f.doRequest("POST", &data, &tool.DoHttpReq{
		Url: GetUserUrl,
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

func (f Fs) LoadDepartmentList() (*ListDepartmentResp, error) {
	tenantToken, e := f.tenant.Load()
	if e != nil {
		return nil, e
	}
	var data ListDepartmentResp
	return &data, f.doRequest("GET", &data, &tool.DoHttpReq{
		Url: ListDepartmentUrl,
		Header: map[string]interface{}{
			"Authorization": "Bearer " + tenantToken,
			"Content-Type":  "application/json",
		},
		Query: map[string]interface{}{
			"page_size": 40,
		},
	})
}

func (f Fs) doLoadUserListRequest(departmentID, pageToken string, pageSize uint) (*ListUserResp, error) {
	tenantToken, e := f.tenant.Load()
	if e != nil {
		return nil, e
	}
	var data ListUserResp
	query := map[string]interface{}{
		"department_id": departmentID,
		"page_size":     pageSize,
	}
	if pageToken != "" {
		query["page_token"] = pageToken
	}
	return &data, f.doRequest("GET", &data, &tool.DoHttpReq{
		Url: ListUsersUrl,
		Header: map[string]interface{}{
			"Authorization": "Bearer " + tenantToken,
			"Content-Type":  "application/json",
		},
		Query: query,
	})
}
func (f Fs) doLoadAllUserListRequest(departmentID string) ([]ListUserContent, error) {
	const pageSize = 99
	var r []ListUserContent
	res, e := f.doLoadUserListRequest(departmentID, "", pageSize)
	if e != nil {
		return nil, e
	}
	r = res.Items
	for res.HasMore {
		res, e = f.doLoadUserListRequest(departmentID, res.PageToken, pageSize)
		if e != nil {
			return nil, e
		}
		r = append(r, res.Items...)
	}
	return r, nil
}

// LoadUserList map key 为飞书部门 OpenID
func (f Fs) LoadUserList() (map[string][]ListUserContent, error) {
	departments, e := f.LoadDepartmentList()
	if e != nil {
		return nil, e
	}

	result := make(map[string][]ListUserContent, len(departments.Items))
	for _, department := range departments.Items {
		var list []ListUserContent
		list, e = f.doLoadAllUserListRequest(department.OpenDepartmentId)
		if e != nil {
			return nil, e
		}
		result[department.OpenDepartmentId] = list
	}
	return result, nil
}
