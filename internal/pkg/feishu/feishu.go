package feishu

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/agent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/util"
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
		T: "",
		E: func() {
			e := Departments.LoadRelation()
			if e != nil {
				log.Errorf("更新飞书部门异常: %v", e)
			}
		},
	})
	if e != nil {
		log.Fatalln("添加飞书部门定时更新任务失败:", e)
	}
}

var Api = feishu.New(global.Config.Feishu.ClientID, global.Config.Feishu.Secret, tools.Http.Client)

type fsDepartmentsRelation struct {
	keywords   []string
	department string
}

var fsDepartmentsRelationMap = []fsDepartmentsRelation{
	{
		keywords:   []string{util.UDev},
		department: util.UDev,
	},
	{
		keywords:   []string{util.UAdm, "HR", "管理线"},
		department: util.UAdm,
	},
	{
		keywords:   []string{util.UCe, "发展"},
		department: util.UCe,
	},
	{
		keywords:   []string{util.UDes, "视觉"},
		department: util.UDes,
	},
	{
		keywords:   []string{util.UOp},
		department: util.UOp,
	},
	{
		keywords:   []string{util.UGame},
		department: util.UGame,
	},
	{
		keywords:   []string{util.UPm},
		department: util.UPm,
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
