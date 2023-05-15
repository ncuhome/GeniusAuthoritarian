package dto

type AppShow struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	AppCode        string `json:"appCode"`
	Callback       string `json:"callback"`
	PermitAllGroup bool   `json:"permitAllGroup"`
}

type AppShowDetail struct {
	AppShow
	Groups []Group `json:"groups" gorm:"-"`
}

type AppNew struct {
	AppShowDetail
	AppSecret string `json:"appSecret"`
}
