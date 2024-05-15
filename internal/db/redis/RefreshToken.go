package redis

import (
	"context"
	"encoding/json"
	"github.com/Mmx233/tool"
	"github.com/go-redis/redis/v8"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
	log "github.com/sirupsen/logrus"
	"go/types"
	"strconv"
	"time"
	"unsafe"
)

type RefreshTokenStore struct {
	tokenStore.TokenStore[types.Nil]
}

func NewRecordedToken() RefreshTokenStore {
	return RefreshTokenStore{
		tokenStore.NewTokenStore[types.Nil](Client, keyRecordedToken.String()),
	}
}

func (store RefreshTokenStore) NewStorePoint(id uint64) RefreshTokenStorePoint {
	return RefreshTokenStorePoint{
		Point: store.TokenStore.NewStorePoint(id),
		id:    id,
	}
}

type RefreshTokenStorePoint struct {
	tokenStore.Point[types.Nil]
	id uint64
}

func (point RefreshTokenStorePoint) Destroy(ctx context.Context, appCode string, validBefore time.Time) error {
	canceledToken := CanceledToken{
		ID: point.id,
		CanceledTokenPayload: CanceledTokenPayload{
			AppCode:     appCode,
			ValidBefore: validBefore,
		},
	}
	err := NewCanceledToken().Add(ctx, canceledToken)
	if err != nil {
		return err
	}
	if err = NewCanceledTokenChannel().Publish(ctx, canceledToken); err != nil {
		return err
	}
	return point.Point.Destroy(ctx)
}

type CanceledTokenChannel struct {
	key string
}

func NewCanceledTokenChannel() CanceledTokenChannel {
	return CanceledTokenChannel{
		key: keyCanceledToken.String() + "sub",
	}
}

func (channel CanceledTokenChannel) Publish(ctx context.Context, token CanceledToken) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return Client.Publish(ctx, channel.key, data).Err()
}

func (channel CanceledTokenChannel) Subscribe(ctx context.Context) *redis.PubSub {
	return Client.Subscribe(ctx, channel.key)
}

func NewCanceledToken() CanceledTokenTable {
	return CanceledTokenTable{
		key: keyCanceledToken.String() + "table",
	}
}

type CanceledTokenPayload struct {
	AppCode     string    `json:"appCode"`
	ValidBefore time.Time `json:"validBefore"`
}

type CanceledToken struct {
	ID uint64 `json:"id"`
	CanceledTokenPayload
}

func (token CanceledToken) Key() string {
	return strconv.FormatUint(token.ID, 10)
}

type CanceledTokenTable struct {
	key string
}

func (table CanceledTokenTable) Add(ctx context.Context, tokens ...CanceledToken) error {
	fields := make([]interface{}, len(tokens)*2)
	for i, v := range tokens {
		fields[i*2] = v.Key()
		data, err := json.Marshal(v.CanceledTokenPayload)
		if err != nil {
			return err
		}
		fields[i*2+1] = data
	}
	return Client.HSet(ctx, table.key, fields...).Err()
}

func (table CanceledTokenTable) Get(ctx context.Context) ([]CanceledToken, error) {
	result, err := Client.HGetAll(ctx, table.key).Result()
	if err != nil {
		return nil, err
	}
	canceledTokens := make([]CanceledToken, len(result))
	left, right := 0, len(result)-1
	for k, v := range result {
		var canceledToken CanceledToken
		var err error
		canceledToken.ID, err = strconv.ParseUint(k, 10, 64)
		if err != nil {
			log.Errorln("parse id failed", err)
			continue
		}
		if err = json.Unmarshal(unsafe.Slice(unsafe.StringData(v), len(v)), &canceledToken.CanceledTokenPayload); err != nil {
			log.Errorln("parse canceled token failed", err)
			continue
		}
		if canceledToken.ValidBefore.After(time.Now()) {
			canceledTokens[left] = canceledToken
			left++
		} else {
			canceledTokens[right] = canceledToken
			right--
		}
	}
	if left != len(result)-1 {
		go table.clean(canceledTokens[left+1:]...)
	}
	return canceledTokens[:left], nil
}

func (table CanceledTokenTable) clean(tokens ...CanceledToken) {
	defer tool.Recover()
	err := table.remove(context.Background(), tokens...)
	if err != nil {
		log.Errorln("clean canceled token failed", err)
	}
}

func (table CanceledTokenTable) remove(ctx context.Context, tokens ...CanceledToken) error {
	keyGroup := make([]string, len(tokens))
	for i, v := range tokens {
		keyGroup[i] = v.Key()
	}
	return Client.HDel(ctx, table.key, keyGroup...).Err()

}
