package signature

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Gen 传入的必须是结构体
func Gen(form any) string {
	v := reflect.ValueOf(form).Elem()
	t := v.Type()
	num := t.NumField()
	var signMap = make(map[string]string, num)
	var keySlice = make([]string, num)
	var signStrLen = num*2 - 1
	for i := 0; i < num; i++ {
		key := t.Field(i).Tag.Get("json")
		keySlice[i] = key
		signMap[key] = fmt.Sprint(v.Field(i).Interface())
		signStrLen += len(key) + len(signMap[key])
	}
	sort.Strings(keySlice)

	var signStr strings.Builder
	signStr.Grow(signStrLen)
	for i, key := range keySlice {
		if i != 0 {
			signStr.Write([]byte("&"))
		}
		signStr.Write([]byte(key + "=" + signMap[key]))
	}

	h := sha256.New()
	h.Write([]byte(signStr.String()))
	return fmt.Sprintf("%x", h.Sum(nil))
}

type VerifyClaims struct {
	Token     string `json:"token"`
	AppCode   string `json:"appCode"`
	TimeStamp int64  `json:"timeStamp"`
	AppSecret string `json:"appSecret"`
}

func Check(signature string, claims *VerifyClaims) bool {
	return signature == Gen(claims)
}
