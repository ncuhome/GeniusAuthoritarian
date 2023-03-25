package util

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/tools"
	"strings"
	"sync"
)

var Feishu = feishu.New(global.Config.Feishu.ClientID, global.Config.Feishu.Secret, tools.Http.Client)

type fsDepartmentsRelation struct {
	keywords   []string
	department string
}

var fsDepartmentsRelationMap = []fsDepartmentsRelation{
	{
		keywords:   []string{UDev},
		department: UDev,
	},
	{
		keywords:   []string{UAdm, "HR", "管理线"},
		department: UAdm,
	},
	{
		keywords:   []string{UCe, "发展"},
		department: UCe,
	},
	{
		keywords:   []string{UDes, "视觉"},
		department: UDes,
	},
	{
		keywords:   []string{UOp},
		department: UOp,
	},
	{
		keywords:   []string{UGame},
		department: UGame,
	},
	{
		keywords:   []string{UPm},
		department: UPm,
	},
}

type fsDepartments struct {
	ml sync.RWMutex
	// 部门 id to name
	m map[string]string
}

var FsDepartments fsDepartments

func (d *fsDepartments) LoadRelation() error {
	m, e := Feishu.LoadDepartmentList()
	if e != nil {
		return e
	}

	var result = make(map[string]string, len(fsDepartmentsRelationMap))

	for oid, name := range m {
		for _, r := range fsDepartmentsRelationMap {
			for _, k := range r.keywords {
				if strings.Contains(name, k) {
					result[oid] = r.department
					goto next
				}
			}
		}
	next:
	}

	d.ml.Lock()
	defer d.ml.Unlock()
	d.m = result
	return nil
}

func (d *fsDepartments) Search(id string) (string, bool) {
	d.ml.RLock()
	defer d.ml.RUnlock()

	n, ok := d.m[id]
	return n, ok
}
