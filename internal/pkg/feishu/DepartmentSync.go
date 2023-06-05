package feishu

import (
	"github.com/Mmx233/daoUtil"
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/GroupOperator"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/agent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func DepartmentSync() error {
	departmentList, e := Api.LoadDepartmentList()
	if e != nil {
		return e
	}

	// 匹配所有命中关键词的部门，以组名为索引避免出现多个匹配结果
	var pairedDepartments = make(map[string]*dao.FeishuGroups, len(departmentList.Items))
	for _, item := range departmentList.Items {
		for _, relate := range fsDepartmentsRelationMap {
			for _, keyword := range relate.keywords {
				if strings.Contains(item.Name, keyword) {
					pairedDepartments[relate.department] = &dao.FeishuGroups{
						Name:             item.Name,
						OpenDepartmentId: item.OpenDepartmentId,
						GID:              GroupOperator.GroupRelation[relate.department],
					}
					goto next
				}
			}
		}
	next:
	}

	// 转换为组 ID 为索引的映射
	var pairedDepartmentsRelations = make(map[uint]*dao.FeishuGroups, len(pairedDepartments))
	for _, v := range pairedDepartments {
		pairedDepartmentsRelations[v.GID] = v
	}

	var toDelete []uint
	var toCreate []dao.FeishuGroups
	srv, e := service.FeishuGroups.Begin()
	if e != nil {
		return e
	}
	defer srv.Rollback()

	dbFeishuDepartments, e := srv.GetAll(daoUtil.LockForUpdate)
	if e != nil {
		return e
	}

	// 计算数据库 diff
	for _, dbDepartment := range dbFeishuDepartments {
		paired, ok := pairedDepartmentsRelations[dbDepartment.GID]
		if !ok {
			toDelete = append(toDelete, dbDepartment.ID)
			continue
		}
		if paired.Name == dbDepartment.Name && paired.OpenDepartmentId == dbDepartment.OpenDepartmentId {
			delete(pairedDepartmentsRelations, dbDepartment.GID)
		} else {
			toDelete = append(toDelete, dbDepartment.ID)
		}
	}
	toCreate = make([]dao.FeishuGroups, len(pairedDepartmentsRelations))
	i := 0
	for _, department := range pairedDepartmentsRelations {
		toCreate[i] = *department
		i++
	}

	if len(toDelete) != 0 {
		if e = srv.DeleteSelected(toDelete); e != nil {
			return e
		}
	}
	if len(toCreate) != 0 {
		if e = srv.CreateAll(toCreate); e != nil {
			return e
		}
	}
	return srv.Commit().Error
}

func AddDepartmentSyncCron(spec string) error {
	_, e := agent.AddRegular(&agent.Event{
		T: spec,
		E: func() {
			defer tool.Recover()
			startAt := time.Now()
			if e := DepartmentSync(); e != nil {
				log.Errorf("同步飞书部门失败: %v", e)
			} else {
				log.Infof("飞书部门同步成功，耗时 %dms", time.Now().Sub(startAt).Milliseconds())
			}
		},
	})
	return e
}

type fsDepartmentsRelation struct {
	keywords   []string
	department string
}

var fsDepartmentsRelationMap = []fsDepartmentsRelation{
	{
		keywords:   []string{departments.UDev},
		department: departments.UDev,
	},
	{
		keywords:   []string{departments.UAdm, "HR", "管理线"},
		department: departments.UAdm,
	},
	{
		keywords:   []string{departments.UCe, "发展"},
		department: departments.UCe,
	},
	{
		keywords:   []string{departments.UDes, "视觉"},
		department: departments.UDes,
	},
	{
		keywords:   []string{departments.UOp},
		department: departments.UOp,
	},
	{
		keywords:   []string{departments.UGame},
		department: departments.UGame,
	},
	{
		keywords:   []string{departments.UPm},
		department: departments.UPm,
	},
	{
		keywords:   []string{departments.USenior},
		department: departments.USenior,
	},
}
