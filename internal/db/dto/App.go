package dto

type AppShow struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Callback       string `json:"callback"`
	PermitAllGroup bool   `json:"permitAllGroup"`
}

type AppShowOwner struct {
	AppShow
	AppCode string `json:"appCode"`
}

type AppShowDetail struct {
	AppShowOwner
	Groups []Group `json:"groups" gorm:"-"`
}

type AppNew struct {
	AppShowDetail
	AppSecret string `json:"appSecret"`
}

type AppShowWithGroup struct {
	AppShow
	GroupID   uint   `json:"groupID"`
	GroupName string `json:"groupName"`
}
