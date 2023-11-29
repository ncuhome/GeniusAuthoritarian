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
	"gorm.io/gorm"
	"net"
	"unsafe"
)

func Run(addr string) error {
	tcpListen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middlewares.UnaryLogger(),
		),
		grpc.ChainStreamInterceptor(
			middlewares.StreamLogger(),
		),
	)
	refreshTokenProto.RegisterRefreshTokenServer(grpcServer, &Server{})

	return grpcServer.Serve(tcpListen)
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

	accessToken, err := jwt.GenerateAccessToken(claims.UID, req.AppCode, claims.Payload)
	if err != nil {
		return nil, status.Error(codes.Internal, "generate access token failed")
	}

	return &refreshTokenProto.AccessToken{
		AccessToken: accessToken,
		Payload:     claims.Payload,
	}, nil
}

func (s *Server) DestroyRefreshToken(stream refreshTokenProto.RefreshToken_DestroyRefreshTokenServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return status.Error(codes.Internal, "receive request failed")
		}

		err = s.CheckSignature(req)
		if err != nil {
			return err
		}

		claims, valid, err := jwt.ParseRefreshToken(req.Token)
		if err != nil || !valid {
			return status.Error(codes.Unauthenticated, "token invalid")
		}

		err = redis.NewRefreshToken().NewStorePoint(claims.ID).Destroy(context.Background())
		if err != nil {
			return status.Error(codes.Internal, "destroy token failed")
		}
	}
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
