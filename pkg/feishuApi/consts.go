package feishuApi

const (

	// GetTenantAccessTokenUrl 自建应用获取 tenant_access_token
	GetTenantAccessTokenUrl = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"

	// LarkLoginUrl 飞书登录接口
	LarkLoginUrl = "https://open.feishu.cn/open-apis/authen/v1/user_auth_page_beta?app_id=%s&state=%s&redirect_uri=%s"

	// GetUserUrl 获取用户
	GetUserUrl = "https://open.feishu.cn/open-apis/authen/v1/access_token"

	//ListDepartmentUrl 获取部门列表
	ListDepartmentUrl = "https://open.feishu.cn/open-apis/contact/v3/departments/0/children"

	// ListUsersUrl 列出所有用户
	ListUsersUrl = "https://open.feishu.cn/open-apis/contact/v3/users"

	// GetUserInfoUrl 获取用户信息
	GetUserInfoUrl = "https://open.feishu.cn/open-apis/contact/v3/users/%s"
)
