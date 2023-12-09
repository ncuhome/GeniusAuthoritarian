package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mmx233/tool"
	"math/rand"
	"time"
)

func NewUserIdentityCode(uid uint) UserIdentityCode {
	return UserIdentityCode{
		key: keyUserIdentityCode.String() + fmt.Sprint(uid),
	}
}

type UserIdentityCode struct {
	key string
}

// New 新建身份校验码，五分钟有效，每用户仅存在一个
func (a UserIdentityCode) New() (string, error) {
	randSource := rand.NewSource(time.Now().UnixNano())
	code := fmt.Sprint(tool.RandNum(rand.New(randSource), 12345, 99999))
	return code, Client.Set(context.Background(), a.key, code, time.Minute*5).Err()
}

// VerifyAndDestroy 校验并销毁 code
func (a UserIdentityCode) VerifyAndDestroy(code string) (bool, error) {
	rCode, err := Client.Get(context.Background(), a.key).Result()
	if err != nil {
		if errors.Is(err, Nil) {
			err = nil
		}
		return false, err
	}
	return rCode == code, a.Destroy()
}

func (a UserIdentityCode) Destroy() error {
	return Client.Del(context.Background(), a.key).Err()
}
