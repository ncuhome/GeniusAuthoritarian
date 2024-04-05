package refreshToken

import (
	"context"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/middlewares"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/refreshToken/refreshTokenProto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (s *Server) RefreshToken(_ context.Context, req *refreshTokenProto.TokenRequest) (*refreshTokenProto.AccessToken, error) {
	claims, err := CheckRefreshToken(req)
	if err != nil {
		return nil, err
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
	claims, err := CheckRefreshToken(req)
	if err != nil {
		return nil, err
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
	claims, err := CheckAccessToken(req)
	if err != nil {
		return nil, err
	}

	return &refreshTokenProto.AccessTokenInfo{
		Uid:     uint64(claims.UID),
		Payload: claims.Payload,
	}, nil
}

func (s *Server) GetUserInfo(_ context.Context, req *refreshTokenProto.TokenRequest) (*refreshTokenProto.UserInfo, error) {
	claims, err := CheckAccessToken(req)
	if err != nil {
		return nil, err
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
