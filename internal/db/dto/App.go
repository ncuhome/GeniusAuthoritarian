package dto

type AppShow struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	AppCode        string `json:"appCode"`
	PermitAllGroup bool   `json:"permitAllGroup"`
}

type AppNew struct {
	AppShow
	AppSecret string `json:"appSecret"`
}

type AppShowDetail struct {
	AppShow
	GroupCount uint    `json:"groupCount"`
	Groups     []Group `json:"groups" gorm:"-"`
}
