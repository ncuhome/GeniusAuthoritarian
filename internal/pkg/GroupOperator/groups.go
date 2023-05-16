package GroupOperator

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
	log "github.com/sirupsen/logrus"
)

var Groups = []string{
	departments.UDev,
	departments.UPm,
	departments.UGame,
	departments.UOp,
	departments.UAdm,
	departments.UDes,
	departments.UCe,
	departments.USenior,
}

var GroupRelation map[string]uint

func init() {
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

	var groupRelations = make(map[string]uint, len(Groups))

	var notExistGroups []string
	for _, group := range Groups {
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

	GroupRelation = groupRelations
	return nil
}
