package department2BaseGroup

import (
	"container/list"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	log "github.com/sirupsen/logrus"
)

func Init() {
	err := CheckBaseGroups()
	if err != nil {
		log.Fatalln("写入基本组列表失败:", err)
	}
}

func CheckBaseGroups() error {
	dbGroups, err := service.BaseGroups.LoadGroups()
	if err != nil {
		return err
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

		err = service.BaseGroups.CreateGroups(newGroups)
		if err != nil {
			return err
		}
	}

	return nil
}
