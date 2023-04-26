package feishu

import (
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/GroupOperator"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/agent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
	log "github.com/sirupsen/logrus"
	"strings"
)

func GroupSync() error {
	departmentList, e := Api.LoadDepartmentList()
	if e != nil {
		return e
	}

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

	var result = make([]dao.FeishuGroups, len(pairedDepartments))
	i := 0
	for _, v := range pairedDepartments {
		result[i] = *v
		i++
	}

	srv, e := service.FeishuGroups.Begin()
	if e != nil {
		return e
	}
	defer srv.Rollback()

	if e = srv.DeleteAll(); e != nil {
		return e
	}
	if e = srv.CreateAll(result); e != nil {
		return e
	}
	return srv.Commit().Error
}

func RunGroupSync() error {
	_, e := agent.AddRegular(&agent.Event{
		T: "0 5 * * *",
		E: func() {
			defer tool.Recover()
			if e := GroupSync(); e != nil {
				log.Errorf("同步飞书部门失败: %v", e)
			} else {
				log.Infoln("飞书部门同步成功")
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
}
