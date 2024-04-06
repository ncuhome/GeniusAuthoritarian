package redis

import (
	"context"
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
	log "github.com/sirupsen/logrus"
	"go/types"
	"strconv"
	"time"
)

func NewRecordedToken() tokenStore.TokenStore[types.Nil] {
	return tokenStore.NewTokenStore[types.Nil](Client, keyRecordedToken.String())
}

func NewCanceledToken() CanceledTokenTable {
	return CanceledTokenTable{
		key: keyCanceledToken.String(),
	}
}

type CanceledToken struct {
	ID          uint64
	ValidBefore time.Time
}

func (a CanceledToken) Key() string {
	return strconv.FormatUint(a.ID, 10)
}

func (a CanceledToken) Value() string {
	return a.ValidBefore.Format(time.RFC3339)
}

type CanceledTokenTable struct {
	key string
}

func (a CanceledTokenTable) Add(ctx context.Context, id ...CanceledToken) error {
	fields := make([]interface{}, len(id)*2)
	for i, v := range id {
		fields[i*2] = v.Key()
		fields[i*2+1] = v.Value()
	}
	return Client.HSet(ctx, a.key, fields...).Err()
}

func (a CanceledTokenTable) Get(ctx context.Context) ([]uint64, error) {
	result, err := Client.HGetAll(ctx, a.key).Result()
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, len(result))
	left, right := 0, len(result)-1
	for k, v := range result {
		id, err := strconv.ParseUint(k, 10, 64)
		if err != nil {
			log.Errorln("parse id failed", err)
			continue
		}
		validBefore, err := time.Parse(time.RFC3339, v)
		if err != nil {
			log.Errorln("parse time failed", err)
			continue
		}
		if validBefore.After(time.Now()) {
			ids[left] = id
			left++
		} else {
			ids[right] = id
			right--
		}
	}
	if left != len(result)-1 {
		go a.clean(ids[left+1:]...)
	}
	return ids[:left], nil
}

func (a CanceledTokenTable) clean(id ...uint64) {
	defer tool.Recover()
	err := a.remove(context.Background(), id...)
	if err != nil {
		log.Errorln("clean canceled token failed", err)
	}
}

func (a CanceledTokenTable) remove(ctx context.Context, id ...uint64) error {
	keyGroup := make([]string, len(id))
	for i, v := range id {
		keyGroup[i] = strconv.FormatUint(v, 10)
	}
	return Client.HDel(ctx, a.key, keyGroup...).Err()

}
