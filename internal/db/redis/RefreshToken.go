package redis

import (
	"context"
	"github.com/Mmx233/tool"
	"github.com/go-redis/redis/v8"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
	log "github.com/sirupsen/logrus"
	"go/types"
	"strconv"
	"time"
)

func CancelToken(ctx context.Context, id uint64, validBefore time.Time) error {
	err := NewCanceledToken().Add(ctx, CanceledToken{
		ID:          id,
		ValidBefore: validBefore,
	})
	if err != nil {
		return err
	}
	if err = NewCanceledTokenChannel().Publish(ctx, id); err != nil {
		return err
	}
	return NewRecordedToken().NewStorePoint(id).Destroy(ctx)
}

func NewRecordedToken() tokenStore.TokenStore[types.Nil] {
	return tokenStore.NewTokenStore[types.Nil](Client, keyRecordedToken.String())
}

func NewCanceledTokenChannel() CanceledTokenChannel {
	return CanceledTokenChannel{
		key: keyCanceledToken.String() + "sub",
	}
}

type CanceledTokenChannel struct {
	key string
}

func (a CanceledTokenChannel) Publish(ctx context.Context, id uint64) error {
	return Client.Publish(ctx, a.key, strconv.FormatUint(id, 10)).Err()
}

func (a CanceledTokenChannel) Subscribe(ctx context.Context) *redis.PubSub {
	return Client.Subscribe(ctx, a.key)
}

func NewCanceledToken() CanceledTokenTable {
	return CanceledTokenTable{
		key: keyCanceledToken.String() + "table",
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
