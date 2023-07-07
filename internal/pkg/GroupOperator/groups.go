package GroupOperator

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	log "github.com/sirupsen/logrus"
)

func InitGroupRelation() {
	e := LoadGroupRelation()
	if e != nil {
		log.Fatalf("载入部门 id 关系失败: %v", e)
	}
}

func LoadGroupRelation() error {
	dbGroups, e := service.BaseGroups.LoadGroups()
	if e != nil {
		return e
	}

	var groupRelations = make(map[string]uint, len(global.Departments))

	var notExistGroups []string
	for _, group := range global.Departments {
		for _, dbGroup := range dbGroups {
			if group == dbGroup.Name {
				groupRelations[group] = dbGroup.ID
				goto next
			}
		}
		notExistGroups = append(notExistGroups, group)
	next:
	}

	if len(notExistGroups) != 0 {
		var newGroups []dao.BaseGroup
		newGroups, e = service.BaseGroups.CreateGroups(notExistGroups)
		if e != nil {
			return e
		}
		for _, group := range newGroups {
			groupRelations[group.Name] = group.ID
		}
	}

	global.DepartmentRelation = groupRelations
	return nil
}
