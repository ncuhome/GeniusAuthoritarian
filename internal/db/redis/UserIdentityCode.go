package redis

import (
	"context"
	"fmt"
	"github.com/Mmx233/tool"
	"math/rand"
	"time"
)

var UserIdentityCode = UserIdentityCodeHelper{
	key: keyUserIdentityCode.String(),
}

type UserIdentityCodeHelper struct {
	key string
}

func (a UserIdentityCodeHelper) genKey(uid uint) string {
	return a.key + "-" + fmt.Sprint(uid)
}

// New 新建身份校验码，五分钟有效，每用户仅存在一个
func (a UserIdentityCodeHelper) New(uid uint) (string, error) {
	randSource := rand.NewSource(time.Now().UnixNano())
	code := fmt.Sprint(tool.RandNum(rand.New(randSource), 12345, 99999))
	return code, Client.Set(context.Background(), a.genKey(uid), code, time.Minute*5).Err()
}

// VerifyAndDestroy 校验并销毁 code
func (a UserIdentityCodeHelper) VerifyAndDestroy(uid uint, code string) (bool, error) {
	rCode, err := Client.Get(context.Background(), a.genKey(uid)).Result()
	if err != nil {
		if err == Nil {
			err = nil
		}
		return false, err
	}
	return rCode == code, a.Destroy(uid)
}

func (a UserIdentityCodeHelper) Destroy(uid uint) error {
	return Client.Del(context.Background(), a.genKey(uid)).Err()
}
