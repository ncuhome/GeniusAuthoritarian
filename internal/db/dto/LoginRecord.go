package dto

type LoginRecord struct {
	ID        uint   `json:"id"`
	CreatedAt int64  `json:"createdAt"`
	Useragent string `json:"useragent"`
	IP        string `json:"ip"`
	Target    string `json:"target"`
	Method    string `json:"method"`
}

type LoginRecordOnline struct {
	LoginRecord
	ValidBefore uint64 `json:"validBefore"`

	IsMe bool `json:"isMe" gorm:"-"`
}

type ViewCount struct {
	// App.ID
	ID    uint   `json:"id"`
	Views uint64 `json:"views"`
}

type ViewID struct {
	// LoginRecord.ID
	ID uint
	// App.ID
	AID uint `gorm:"column:aid"`
}

type LoginRecordDataView struct {
	ID        uint  `json:"id"`
	CreatedAt int64 `json:"createdAt"`

	UID uint `json:"uid" gorm:"column:uid"`
	AID uint `json:"aid" gorm:"column:aid"`

	Method      string `json:"method"`
	Destroyed   bool   `json:"destroyed"`
	ValidBefore uint64 `json:"validBefore"`
}

type AdminLoginDataView struct {
	Apps    []AppDataView         `json:"apps"`
	Records []LoginRecordDataView `json:"records"`
}
