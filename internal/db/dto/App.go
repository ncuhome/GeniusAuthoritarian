package dto

type AppShow struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Callback       string `json:"callback"`
	PermitAllGroup bool   `json:"permitAllGroup"`
	LinkOff        bool   `json:"linkOff"`
	Views          uint64 `json:"views"`
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

type AppGroupClassified struct {
	Group Group     `json:"group"`
	App   []AppShow `json:"app"`
}
