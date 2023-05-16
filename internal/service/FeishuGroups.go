package service

import (
	"github.com/Mmx233/daoUtil"
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

func (a FeishuGroupsSrv) GetAll(opts ...daoUtil.ServiceOpt) ([]dao.FeishuGroups, error) {
	return (&dao.FeishuGroups{}).GetAll(daoUtil.TxOpts(a.DB, opts...))
}

func (a FeishuGroupsSrv) GetByOpenID(openID []string, opts ...daoUtil.ServiceOpt) ([]dao.FeishuGroups, error) {
	return (&dao.FeishuGroups{}).GetByOpenIDSlice(daoUtil.TxOpts(a.DB, opts...), openID)
}

func (a FeishuGroupsSrv) DeleteSelected(ids []uint) error {
	return (&dao.FeishuGroups{}).DeleteByIDSlice(a.DB, ids)
}

func (a FeishuGroupsSrv) CreateAll(data []dao.FeishuGroups) error {
	return (&dao.FeishuGroups{}).CreateAll(a.DB, data)
}

func (a FeishuGroupsSrv) Search(openID []string) ([]dao.BaseGroup, error) {
	return (&dao.FeishuGroups{}).GetGroupsByOpenIDSlice(a.DB, openID)
}
