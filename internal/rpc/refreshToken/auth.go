package refreshToken

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt/jwtClaims"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/refreshToken/refreshTokenProto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"unsafe"
)

func CheckSignature(req *refreshTokenProto.TokenRequest) error {
	_, appSecret, err := service.App.FirstAppKeyPairByAppCode(req.AppCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Error(codes.Unauthenticated, "appCode not found")
		}
		return status.Error(codes.Internal, "database error")
	}

	signStr := fmt.Sprintf("%s:%s:%s", req.AppCode, appSecret, req.Token)
	h := sha256.New()
	h.Write(unsafe.Slice(unsafe.StringData(signStr), len(signStr)))
	if req.Signature != fmt.Sprintf("%x", h.Sum(nil)) {
		return status.Error(codes.Unauthenticated, "signature invalid")
	}
	return nil
}

func CheckRefreshToken(req *refreshTokenProto.TokenRequest) (*jwtClaims.RefreshToken, error) {
	err := CheckSignature(req)
	if err != nil {
		return nil, err
	}

	claims, valid, err := jwt.ParseRefreshToken(req.Token)
	if err != nil || !valid || claims.AppCode != req.AppCode {
		return nil, status.Error(codes.Unauthenticated, "token invalid")
	}
	return claims, nil
}

func CheckAccessToken(req *refreshTokenProto.TokenRequest) (*jwtClaims.AccessToken, error) {
	err := CheckSignature(req)
	if err != nil {
		return nil, err
	}

	claims, valid, err := jwt.ParseAccessToken(req.Token)
	if err != nil || !valid || claims.AppCode != req.AppCode {
		return nil, status.Error(codes.Unauthenticated, "token invalid")
	}
	return claims, nil
}
