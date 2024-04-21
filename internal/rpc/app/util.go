package app

import "context"

const keyAuthInfo = "auth-info"

type AuthInfo struct {
	AppCode string
}

func SetAuthInfo(ctx context.Context, info *AuthInfo) context.Context {
	return context.WithValue(ctx, keyAuthInfo, info)
}
func GetAuthInfo(ctx context.Context) *AuthInfo {
	return ctx.Value(keyAuthInfo).(*AuthInfo)
}
