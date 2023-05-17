package service

import (
	"fmt"
	"github.com/Mmx233/daoUtil"
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

var App = AppSrv{dao.DB}

type AppSrv struct {
	*gorm.DB
}

func (a AppSrv) Begin() (AppSrv, error) {
	a.DB = a.DB.Begin()
	return a, a.Error
}

func (a AppSrv) This(host string) *dao.App {
	return &dao.App{
		Name:           global.ThisAppName,
		Callback:       fmt.Sprintf("https://%s/login", host),
		PermitAllGroup: true,
	}
}

func (a AppSrv) New(uid uint, name, callback string, permitAll bool) (*dao.App, error) {
	randSrc := rand.NewSource(time.Now().UnixNano())
	var t = dao.App{
		Name:           name,
		UID:            uid,
		AppCode:        tool.RandString(randSrc, 8),
		AppSecret:      tool.RandString(randSrc, 100),
		Callback:       callback,
		PermitAllGroup: permitAll,
	}
	return &t, t.Insert(a.DB)
}

func (a AppSrv) NameExist(name string) (bool, error) {
	return (&dao.App{Name: name}).NameExist(a.DB)
}

func (a AppSrv) AppCodeExist(appCode string) (bool, error) {
	list, e := redis.AppCode.Load()
	if e != nil {
		if e == redis.Nil {
			list, e = (&dao.App{}).GetAppCode(a.DB)
			if e != nil {
				return false, e
			}
			_ = redis.AppCode.Add(list...)
		} else {
			return false, e
		}
	}

	for _, v := range list {
		if v == appCode {
			return true, nil
		}
	}
	return false, nil
}

func (a AppSrv) CheckAppCode(appCode string) (bool, error) {
	return a.AppCodeExist(appCode)
}

func (a AppSrv) FirstAppByAppCode(appCode string) (*dao.App, error) {
	var t = dao.App{
		AppCode: appCode,
	}
	return &t, t.FirstForLogin(a.DB)
}

func (a AppSrv) FirstAppDetailedByIDForUser(id, uid uint, opts ...daoUtil.ServiceOpt) (*dto.AppShowDetail, error) {
	appDetailed, e := (&dao.App{
		ID:  id,
		UID: uid,
	}).FirstDetailedByIdAndUID(a.DB)
	if e != nil {
		return nil, e
	}

	groups, e := (&dao.BaseGroup{}).GetByAppIdsRelatedForShow(a.DB, appDetailed.ID)
	if e != nil {
		return nil, e
	}
	appDetailed.Groups = make([]dto.Group, len(groups))
	for i, group := range groups {
		appDetailed.Groups[i] = group.Group
	}
	return appDetailed, nil
}

func (a AppSrv) FirstAppKeyPair(id uint) (string, string, error) {
	var t = dao.App{
		ID: id,
	}
	return t.AppCode, t.AppSecret, t.FirstAppKeyPairByID(a.DB)
}

func (a AppSrv) GetUserOwnedApp(uid uint) ([]dto.AppShowDetail, error) {
	apps, e := (&dao.App{UID: uid}).GetByUIDForShowDetailed(a.DB)
	if e != nil {
		return nil, e
	}

	// 获取各 app 授权组
	if len(apps) > 0 {
		var appIds = make([]uint, len(apps))
		for i, app := range apps {
			appIds[i] = app.ID
		}

		var groupRelatedList []dto.GroupRelateApp
		groupRelatedList, e = (&dao.BaseGroup{}).GetByAppIdsRelatedForShow(a.DB, appIds...)
		if e != nil {
			return nil, e
		}

		var groupCount = make(map[uint]int, len(apps))
		for _, groupRelated := range groupRelatedList {
			count, _ := groupCount[groupRelated.AppID]
			groupCount[groupRelated.AppID] = count + 1
		}

		var AppIdToAppIndexMap = make(map[uint]int, len(apps))
		for i, app := range apps {
			AppIdToAppIndexMap[app.ID] = i
			length, _ := groupCount[app.ID]
			apps[i].Groups = make([]dto.Group, length)
		}

		for _, groupRelated := range groupRelatedList {
			appIndex := AppIdToAppIndexMap[groupRelated.AppID]
			apps[appIndex].Groups[len(apps[appIndex].Groups)-groupCount[groupRelated.AppID]] = groupRelated.Group
			groupCount[groupRelated.AppID]--
		}
	}

	return apps, nil
}

func (a AppSrv) GetUserAccessible(uid uint) ([]dto.AppGroupClassified, error) {
	list, e := (&dao.App{UID: uid}).GetAccessible(a.DB)
	if e != nil {
		return nil, e
	}

	var i = -1
	var lastGroupID uint
	var count []int
	for _, app := range list {
		if lastGroupID != app.GroupID {
			i++
			count = append(count, 0)
			lastGroupID = app.GroupID
		}
		count[i]++
	}

	i = -1
	j := 0
	lastGroupID = 0
	var result = make([]dto.AppGroupClassified, len(count))
	for _, app := range list {
		if lastGroupID != app.GroupID {
			i++
			j = 0
			result[i].Group = dto.Group{
				ID:   app.GroupID,
				Name: app.GroupName,
			}
			result[i].App = make([]dto.AppShow, count[i])
			lastGroupID = app.GroupID
		}
		result[i].App[j] = app.AppShow
		j++
	}
	return result, nil
}

func (a AppSrv) GetPermitAll() ([]dto.AppShow, error) {
	return (&dao.App{}).GetPermitAll(a.DB)
}

func (a AppSrv) DeleteByID(id, uid uint) error {
	result := (&dao.App{ID: id, UID: uid}).DeleteByIdForUID(a.DB)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (a AppSrv) UpdateAll(id uint, name, callback string, permitAllGroup bool) error {
	return (&dao.App{
		ID:             id,
		Name:           name,
		Callback:       callback,
		PermitAllGroup: permitAllGroup,
	}).UpdatesByID(a.DB)
}
