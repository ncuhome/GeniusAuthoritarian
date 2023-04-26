package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"gorm.io/gorm"
)

var FeishuGroups = FeishuGroupsSrv{dao.DB}

type FeishuGroupsSrv struct {
	*gorm.DB
}

func (a FeishuGroupsSrv) Begin() (FeishuGroupsSrv, error) {
	a.DB = a.DB.Begin()
	return a, a.Error
}

func (a FeishuGroupsSrv) DeleteAll() error {
	return (&dao.FeishuGroups{}).DeleteAll(a.DB)
}

func (a FeishuGroupsSrv) CreateAll(data []dao.FeishuGroups) error {
	return (&dao.FeishuGroups{}).CreateAll(a.DB, data)
}
