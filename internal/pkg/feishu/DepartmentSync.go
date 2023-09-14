package feishu

import (
	"container/list"
	"github.com/Mmx233/daoUtil"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/agent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/backoff"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func DepartmentSync() error {
	departmentRelation, err := service.BaseGroups.LoadGroupsRelation()
	if err != nil {
		return err
	}

	departmentList, err := Api.LoadDepartmentList()
	if err != nil {
		return err
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
						GID:              departmentRelation[relate.department].ID,
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

	var toDelete = list.New() // uint
	var toCreate []dao.FeishuGroups
	srv, err := service.FeishuGroups.Begin()
	if err != nil {
		return err
	}
	defer srv.Rollback()

	dbFeishuDepartments, err := srv.GetAll(daoUtil.LockForUpdate)
	if err != nil {
		return err
	}

	// 计算数据库 diff
	for _, dbDepartment := range dbFeishuDepartments {
		paired, ok := pairedDepartmentsRelations[dbDepartment.GID]
		if !ok {
			toDelete.PushBack(dbDepartment.ID)
			continue
		}
		if paired.Name == dbDepartment.Name && paired.OpenDepartmentId == dbDepartment.OpenDepartmentId {
			delete(pairedDepartmentsRelations, dbDepartment.GID)
		} else {
			toDelete.PushBack(dbDepartment.ID)
		}
	}
	toCreate = make([]dao.FeishuGroups, len(pairedDepartmentsRelations))
	i := 0
	for _, department := range pairedDepartmentsRelations {
		toCreate[i] = *department
		i++
	}

	if toDelete.Len() != 0 {
		var toDeleteSlice = make([]uint, toDelete.Len())
		el := toDelete.Front()
		i = 0
		for el != nil {
			toDeleteSlice[i] = el.Value.(uint)
			i++
			el = el.Next()
		}
		if err = srv.DeleteSelected(toDeleteSlice); err != nil {
			return err
		}
	}
	if len(toCreate) != 0 {
		if err = srv.CreateAll(toCreate); err != nil {
			return err
		}
	}
	return srv.Commit().Error
}

func AddDepartmentSyncCron(spec string) error {
	worker := backoff.New(backoff.Conf{
		Content: func() error {
			startAt := time.Now()
			err := DepartmentSync()
			if err != nil {
				log.Errorf("同步飞书部门失败: %v", err)
			} else {
				log.Infof("飞书部门同步成功，总耗时 %dms", time.Now().Sub(startAt).Milliseconds())
			}
			return err
		},
		MaxRetryDelay: time.Minute * 30,
	})

	_, err := agent.AddRegular(&agent.Event{
		T: spec,
		E: worker.Start,
	})
	return err
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
