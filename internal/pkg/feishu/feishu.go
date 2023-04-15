package feishu

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/agent"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/tools"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
)

func init() {
	e := Departments.LoadRelation()
	if e != nil {
		log.Fatalln("载入飞书部门失败:", e)
	}
	_, e = agent.AddRegular(&agent.Event{
		T: "0 5 * * *",
		E: func() {
			e := Departments.LoadRelation()
			if e != nil {
				log.Errorf("更新飞书部门异常: %v", e)
				return
			}
			log.Infoln("飞书部门已刷新")
		},
	})
	if e != nil {
		log.Panicln("添加飞书部门定时更新任务失败:", e)
	}
}

var Api = feishu.New(global.Config.Feishu.ClientID, global.Config.Feishu.Secret, tools.Http.Client)

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

type fsDepartments struct {
	ml sync.RWMutex
	// 部门 id to name
	m map[string]string
}

var Departments fsDepartments

func (d *fsDepartments) LoadRelation() error {
	m, e := Api.LoadDepartmentList()
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

func (d *fsDepartments) MultiSearch(id []string) []string {
	var r []string
	for _, i := range id {
		v, ok := d.Search(i)
		if ok {
			r = append(r, v)
		}
	}
	return r
}
