package service

import (
	"container/list"
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

		LinkOff: true,
	}
	return &t, t.Insert(a.DB)
}

func (a AppSrv) AppCodeExist(appCode string) (bool, error) {
	empty, err := redis.AppCode.IsEmpty()
	if err != nil {
		return false, err
	} else if empty {
		appCodeList, err := (&dao.App{}).GetAppCode(a.DB)
		if err != nil {
			return false, err
		}
		err = redis.AppCode.Add(appCodeList...)
		if err != nil {
			return false, err
		}
	}

	return redis.AppCode.Exist(appCode)
}

func (a AppSrv) UserAccessible(id, uid uint) (bool, error) {
	return (&dao.App{
		ID:  id,
		UID: uid,
	}).UserAccessible(a.DB)
}

func (a AppSrv) CheckAppCode(appCode string) (bool, error) {
	return a.AppCodeExist(appCode)
}

func (a AppSrv) FirstAppByID(id uint) (*dao.App, error) {
	var t = dao.App{
		ID: id,
	}
	return &t, t.FirstByID(a.DB)
}

func (a AppSrv) FirstAppCallbackByID(id uint) (string, error) {
	var t = dao.App{
		ID: id,
	}
	return t.Callback, t.FirstCallbackByID(a.DB)
}

func (a AppSrv) FirstAppCodeByID(id, uid uint, opts ...daoUtil.ServiceOpt) (string, error) {
	var t = dao.App{
		ID:  id,
		UID: uid,
	}
	return t.AppCode, t.FirstAppCodeByID(daoUtil.TxOpts(a.DB, opts...))
}

func (a AppSrv) FirstAppByAppCode(appCode string) (*dao.App, error) {
	var t = dao.App{
		AppCode: appCode,
	}
	return &t, t.FirstByAppCode(a.DB)
}

func (a AppSrv) FirstAppDetailedByIDForUser(id, uid uint, opts ...daoUtil.ServiceOpt) (*dto.AppShowDetail, error) {
	appDetailed, err := (&dao.App{
		ID:  id,
		UID: uid,
	}).FirstDetailedByIdAndUID(daoUtil.TxOpts(a.DB, opts...))
	if err != nil {
		return nil, err
	}

	groups, err := (&dao.BaseGroup{}).GetByAppIdsRelatedForShow(a.DB, appDetailed.ID)
	if err != nil {
		return nil, err
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
	apps, err := (&dao.App{UID: uid}).GetByUIDForShowDetailed(a.DB)
	if err != nil {
		return nil, err
	}

	// 获取各 app 授权组
	if len(apps) > 0 {
		var appIds = make([]uint, len(apps))
		for i, app := range apps {
			appIds[i] = app.ID
		}

		var groupRelatedList []dto.GroupRelateApp
		groupRelatedList, err = (&dao.BaseGroup{}).GetByAppIdsRelatedForShow(a.DB, appIds...)
		if err != nil {
			return nil, err
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

func (a AppSrv) GetUserAccessible(uid uint, isCenterMember bool) ([]dto.AppGroupClassified, error) {
	var appList []dto.AppShowWithGroup
	var err error
	appModel := dao.App{UID: uid}
	if isCenterMember {
		appList, err = appModel.GetAllWithGroup(a.DB)
	} else {
		appList, err = appModel.GetAccessible(a.DB)
	}
	if err != nil {
		return nil, err
	}

	var lastGroupID uint
	var counts = list.New() // *int
	var count *int
	for _, app := range appList {
		if lastGroupID != app.GroupID {
			count = new(int)
			counts.PushBack(count)
			lastGroupID = app.GroupID
		}
		*count++
	}

	i := -1
	countEl := counts.PushFront(nil)
	j := 0
	lastGroupID = 0
	var result = make([]dto.AppGroupClassified, counts.Len()-1)
	for _, app := range appList {
		if lastGroupID != app.GroupID {
			i++
			countEl = countEl.Next()
			j = 0
			result[i].Group = dto.Group{
				ID:   app.GroupID,
				Name: app.GroupName,
			}
			result[i].App = make([]dto.AppShow, *countEl.Value.(*int))
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

func (a AppSrv) GetForUpdateViews(opts ...daoUtil.ServiceOpt) ([]dao.App, error) {
	return (&dao.App{}).GetForUpdateView(daoUtil.TxOpts(a.DB, opts...))
}

func (a AppSrv) DeleteByID(id, uid uint) error {
	return (&dao.App{ID: id, UID: uid}).DeleteByIdForUID(a.DB)
}

func (a AppSrv) UpdateAll(id uint, name, callback string, permitAllGroup bool) error {
	return (&dao.App{
		ID:             id,
		Name:           name,
		Callback:       callback,
		PermitAllGroup: permitAllGroup,
	}).UpdatesByID(a.DB)
}

func (a AppSrv) UpdateViews(id, viewsID uint, views uint64) error {
	return (&dao.App{
		ID:      id,
		Views:   views,
		ViewsID: viewsID,
	}).UpdateViewByID(a.DB)
}

func (a AppSrv) UpdateLinkOff(uid, id uint, linkOff bool) error {
	return (&dao.App{
		ID:      id,
		UID:     uid,
		LinkOff: linkOff,
	}).UpdateLinkOff(a.DB)
}
