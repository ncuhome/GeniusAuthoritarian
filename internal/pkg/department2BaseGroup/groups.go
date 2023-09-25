package department2BaseGroup

import (
	"container/list"
	"context"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	log "github.com/sirupsen/logrus"
	"time"
)

func Init(stat redis.SyncStat) {
	err := stat.MustLock(context.Background(), time.Second*30)
	if err != nil {
		log.Fatalln("初始化 base groups 失败:", err)
	}
	defer stat.Unlock(context.Background())

	if err = CheckBaseGroups(); err != nil {
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
