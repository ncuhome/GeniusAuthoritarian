package department2BaseGroup

import (
	"container/list"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	log "github.com/sirupsen/logrus"
)

func Init() {
	e := CheckBaseGroups()
	if e != nil {
		log.Fatalln("写入基本组列表失败:", e)
	}
}

func CheckBaseGroups() error {
	dbGroups, e := service.BaseGroups.LoadGroups()
	if e != nil {
		return e
	}

	var notExistGroups = list.New() // string
	for _, group := range global.Departments {
		for _, dbGroup := range dbGroups {
			if group == dbGroup.Name {
				goto next
			}
		}
		notExistGroups.PushBack(group)
	next:
	}

	if notExistGroups.Len() != 0 {
		var newGroups = make([]dao.BaseGroup, notExistGroups.Len())
		el := notExistGroups.Front()
		i := 0
		for el != nil {
			newGroups[i].Name = el.Value.(string)
			el = el.Next()
			i++
		}

		e = service.BaseGroups.CreateGroups(newGroups)
		if e != nil {
			return e
		}
	}

	return nil
}
