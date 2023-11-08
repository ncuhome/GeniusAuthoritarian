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
	res, err := f.Http.Request(Method, opt)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	container := &Resp{
		Data: data,
	}
	if err = json.NewDecoder(res.Body).Decode(container); err != nil {
		return err
	} else if container.Code != 0 {
		return fmt.Errorf("feishu api: %s", container.Msg)
	}

	if res.StatusCode > 299 {
		return fmt.Errorf("server return http status %d", res.StatusCode)
	}

	return nil
}

func (f Fs) GetTenantAccessToken() (*TenantAccessTokenResp, error) {
	res, err := f.Http.PostRequest(&tool.DoHttpReq{
		Url: "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal",
		Body: map[string]string{
			"app_id":     f.appID,
			"app_secret": f.secret,
		},
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data TenantAccessTokenResp
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	} else if data.Code != 0 {
		return nil, fmt.Errorf("feishu api: %s", data.Msg)
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("server return http status %d", res.StatusCode)
	}
	return &data, nil
}

func (f Fs) LoginLink(selfDomain, state string) string {
	return fmt.Sprintf(
		"https://open.feishu.cn/open-apis/authen/v1/user_auth_page_beta?app_id=%s&redirect_uri=https%%3A%%2F%%2F"+selfDomain+"%%2Ffeishu%%2F&state=%s",
		f.appID,
		url.QueryEscape(state),
	)
}

func (f Fs) GetUser(code string) (*FsUser, error) {
	tenantToken, err := f.tenant.Load()
	if err != nil {
		return nil, err
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

func (f Fs) LoadDepartmentList() (*ListDepartmentResp, error) {
	tenantToken, err := f.tenant.Load()
	if err != nil {
		return nil, err
	}
	var data ListDepartmentResp
	return &data, f.doRequest("GET", &data, &tool.DoHttpReq{
		Url: "https://open.feishu.cn/open-apis/contact/v3/departments/0/children",
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
	tenantToken, err := f.tenant.Load()
	if err != nil {
		return nil, err
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
		Url: "https://open.feishu.cn/open-apis/contact/v3/users/find_by_department",
		Header: map[string]interface{}{
			"Authorization": "Bearer " + tenantToken,
			"Content-Type":  "application/json",
		},
		Query: query,
	})
}

func (f Fs) doLoadAllUserListRequest(departmentID string) ([]User, error) {
	const pageSize = 50
	var r []User
	res, err := f.doLoadUserListRequest(departmentID, "", pageSize)
	if err != nil {
		return nil, err
	}
	r = res.Items
	for res.HasMore {
		res, err = f.doLoadUserListRequest(departmentID, res.PageToken, pageSize)
		if err != nil {
			return nil, err
		}
		r = append(r, res.Items...)
	}
	return r, nil
}

// LoadUserList 键名为部门 OpenID，只能挨个部门获取用户列表
func (f Fs) LoadUserList() (map[string][]User, error) {
	departments, err := f.LoadDepartmentList()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]User, len(departments.Items))
	for _, department := range departments.Items {
		var list []User
		list, err = f.doLoadAllUserListRequest(department.OpenDepartmentId)
		if err != nil {
			return nil, err
		}
		result[department.OpenDepartmentId] = list
	}
	return result, nil
}
