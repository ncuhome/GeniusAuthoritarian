package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// User2Groups bind user with permission groups
type User2Groups struct {
	ID uint `gorm:"primarykey"`
	// User.ID
	UID  uint  `gorm:"index;index:user_group_idx,unique;not null;column:uid;"`
	User *User `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE"`
	// BaseGroup.ID
	GID       uint       `gorm:"index;index:user_group_idx,unique;not null;column:gid"`
	BaseGroup *BaseGroup `gorm:"foreignKey:GID;constraint:OnDelete:CASCADE"`
}

func (a *User2Groups) TableName() string {
	return "user2groups"
}

func (a *User2Groups) sqlJoinUsers(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN users ON users.id=user2groups.uid AND users.deleted_at IS NULL")
}

func (a *User2Groups) sqlJoinBaseGroups(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN base_groups ON base_groups.id=user2groups.gid")
}

func (a *User2Groups) sqlGetUserGroupsByUID(tx *gorm.DB) *gorm.DB {
	groupModel := &BaseGroup{}
	tx = tx.Model(groupModel)
	tx = groupModel.sqlJoinUserGroups(tx)
	return tx.Where("user2groups.uid=?", a.UID)
}

func (a *User2Groups) ExistByName(tx *gorm.DB, groupName string) (bool, error) {
	var t bool
	tx = tx.Model(a).Select("1")
	tx = a.sqlJoinBaseGroups(tx)
	return t, tx.Where("base_groups.name=? AND user2groups.uid=?", groupName, a.UID).Limit(1).Find(&t).Error
}

func (a *User2Groups) InsertAll(tx *gorm.DB, data []User2Groups) error {
	return tx.Model(a).Create(data).Error
}

func (a *User2Groups) GetUserGroupNamesByUID(tx *gorm.DB) ([]string, error) {
	var t []string
	return t, a.sqlGetUserGroupsByUID(tx).Select("base_groups.name").Find(&t).Error
}

func (a *User2Groups) GetGetUserGroupIdsByUID(tx *gorm.DB) ([]uint, error) {
	var t []uint
	return t, a.sqlGetUserGroupsByUID(tx).Select("base_groups.id").Find(&t).Error
}

func (a *User2Groups) GetUserGroupsForShowByUID(tx *gorm.DB) ([]dto.Group, error) {
	var t = make([]dto.Group, 0)
	return t, a.sqlGetUserGroupsByUID(tx).Find(&t).Error
}

func (a *User2Groups) GetUserGroupsForPubByUIDWithPreOrder(tx *gorm.DB, uid ...uint) ([]dto.GroupWithOrder, error) {
	var t []dto.GroupWithOrder
	groupModel := &BaseGroup{}
	tx = tx.Model(groupModel)
	tx = groupModel.sqlJoinUserGroups(tx)
	return t, tx.Select("base_groups.id", "base_groups.name", "user2groups.uid").Where("user2groups.uid IN ?", uid).
		Clauses(clause.OrderBy{Expression: clause.Expr{SQL: "FIELD(user2groups.uid,?),user2groups.id DESC", Vars: []interface{}{uid}, WithoutParentheses: true}}).
		Find(&t).Error
}

func (a *User2Groups) GetUserGroupsForAppCodeByUID(tx *gorm.DB, appCode string) *gorm.DB {
	appGroupsTx := (&AppGroup{}).sqlGetGroupsByAppCode(tx, appCode).Select("base_groups.name")
	return a.sqlGetUserGroupsByUID(appGroupsTx)
}

// GetUserGroupsLimited 根据指定组范围获取用户所在组
func (a *User2Groups) GetUserGroupsLimited(tx *gorm.DB, groups []string) ([]BaseGroup, error) {
	var t []BaseGroup
	groupModel := &BaseGroup{}
	tx = tx.Model(groupModel)
	tx = groupModel.sqlJoinUserGroups(tx)
	return t, tx.Where("user2groups.uid=? AND base_groups.name IN ?", a.UID, groups).
		Find(&t).Error
}

func (a *User2Groups) GetAllNotFrozen(tx *gorm.DB) ([]User2Groups, error) {
	var t []User2Groups
	tx = tx.Model(a)
	tx = a.sqlJoinUsers(tx)
	return t, tx.Order("uid,gid").Find(&t).Error
}

func (a *User2Groups) DeleteByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Model(a).Delete(a, "id IN ?", ids).Error
}

func (a *User2Groups) DeleteNotInGidSliceByUID(tx *gorm.DB, ids []uint) *gorm.DB {
	return tx.Model(a).Where("uid=? AND gid NOT IN ?", a.UID, ids).Delete(nil)
}
