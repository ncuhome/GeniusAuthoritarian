package feishu

type UserInfo struct {
	AvatarUrl  string `json:"AvatarUrl"`
	Name       string `json:"Name"`
	Email      string `json:"Email"`
	Status     int    `json:"Status"`
	EmployeeId string `json:"EmployeeId"`
	Mobile     string `json:"Mobile"`
}

type TenantAccessTokenResp struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int64  `json:"expire"`
}

type OAuth2AccessTokenResp struct {
	Data OAuth2AccessToken `json:"data"`
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
}

type OAuth2AccessToken struct {
	AccessToken  string `json:"access_token"`
	AvatarUrl    string `json:"avatar_url"`
	ExpiresIn    int64  `json:"expires_in"`
	Name         string `json:"name"`
	EnName       string `json:"en_name"`
	OpenId       string `json:"open_id"`
	TenantKey    string `json:"tenant_key"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
}
