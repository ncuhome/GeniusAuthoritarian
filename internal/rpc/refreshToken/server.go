package refreshToken

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/middlewares"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	refreshTokenProto "github.com/ncuhome/GeniusAuthoritarianProtos/refreshToken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"unsafe"
)

func NewRpc() *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middlewares.UnaryLogger(),
		),
		grpc.ChainStreamInterceptor(
			middlewares.StreamLogger(),
		),
	)
	refreshTokenProto.RegisterRefreshTokenServer(grpcServer, &Server{})
	return grpcServer
}

type Server struct {
	refreshTokenProto.UnimplementedRefreshTokenServer
}

func (s *Server) CheckSignature(req *refreshTokenProto.TokenRequest) error {
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

func (s *Server) RefreshToken(_ context.Context, req *refreshTokenProto.TokenRequest) (*refreshTokenProto.AccessToken, error) {
	err := s.CheckSignature(req)
	if err != nil {
		return nil, err
	}

	claims, valid, err := jwt.ParseRefreshToken(req.Token)
	if err != nil || !valid {
		return nil, status.Error(codes.Unauthenticated, "token invalid")
	}

	accessToken, err := jwt.GenerateAccessToken(claims.ID, claims.UID, req.AppCode, claims.Payload)
	if err != nil {
		return nil, status.Error(codes.Internal, "generate access token failed")
	}

	return &refreshTokenProto.AccessToken{
		AccessToken: accessToken,
		Payload:     claims.Payload,
	}, nil
}

func (s *Server) DestroyRefreshToken(_ context.Context, req *refreshTokenProto.TokenRequest) (*emptypb.Empty, error) {
	err := s.CheckSignature(req)
	if err != nil {
		return nil, err
	}

	claims, valid, err := jwt.ParseRefreshToken(req.Token)
	if err != nil || !valid {
		return nil, status.Error(codes.Unauthenticated, "token invalid")
	}

	loginRecordSrv, err := service.LoginRecord.Begin()
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}
	defer loginRecordSrv.Rollback()

	err = loginRecordSrv.SetDestroyed(uint(claims.ID))
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}

	err = redis.NewRecordedToken().NewStorePoint(claims.ID).Destroy(context.Background())
	if err != nil {
		return nil, status.Error(codes.Internal, "destroy token failed")
	}

	if err = loginRecordSrv.Commit().Error; err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) VerifyAccessToken(_ context.Context, req *refreshTokenProto.TokenRequest) (*refreshTokenProto.AccessTokenInfo, error) {
	err := s.CheckSignature(req)
	if err != nil {
		return nil, err
	}

	claims, valid, err := jwt.ParseAccessToken(req.Token)
	if err != nil || !valid {
		return nil, status.Error(codes.Unauthenticated, "token invalid")
	}

	return &refreshTokenProto.AccessTokenInfo{
		Uid:     uint64(claims.UID),
		Payload: claims.Payload,
	}, nil
}

func (s *Server) GetUserInfo(_ context.Context, req *refreshTokenProto.TokenRequest) (*refreshTokenProto.UserInfo, error) {
	err := s.CheckSignature(req)
	if err != nil {
		return nil, err
	}

	claims, valid, err := jwt.ParseAccessToken(req.Token)
	if err != nil || !valid {
		return nil, status.Error(codes.Unauthenticated, "token invalid")
	}

	user, err := service.User.UserInfoByID(claims.UID)
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}

	groups, err := service.UserGroups.GetNamesForUser(user.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}

	return &refreshTokenProto.UserInfo{
		Uid:       uint64(user.ID),
		Name:      user.Name,
		AvatarUrl: user.AvatarUrl,
		Groups:    groups,
	}, nil
}
