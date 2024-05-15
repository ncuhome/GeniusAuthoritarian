// rpc server for applications
// mutual ssl verify is used

package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/app/appProto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/middlewares"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
	"time"
	"unsafe"
)

func NewRpc() *grpc.Server {
	caPool := x509.NewCertPool()
	caPool.AddCert(global.CaIssuer.CaCert)

	var rpcCert *tls.Certificate
	var rpcCertValidBefore time.Time
	var rpcCertLock sync.RWMutex

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middlewares.UnaryLogger(),
			UnaryAuth(),
		),
		grpc.ChainStreamInterceptor(
			middlewares.StreamLogger(),
			StreamAuth(),
		),
		grpc.Creds(credentials.NewTLS(&tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				rpcCertLock.RLock()
				if !rpcCertValidBefore.IsZero() && rpcCertValidBefore.After(time.Now().Add(time.Minute*5)) {
					defer rpcCertLock.RUnlock()
					return rpcCert, nil
				}
				rpcCertLock.RUnlock()
				rpcCertLock.Lock()
				defer rpcCertLock.Unlock()
				if rpcCertValidBefore.IsZero() || rpcCertValidBefore.Before(time.Now().Add(time.Minute*5)) {
					fullChain, key, err := global.CaIssuer.IssueServer(info.ServerName, time.Now().AddDate(0, 1, 0))
					if err != nil {
						return nil, err
					}
					newCert, err := tls.X509KeyPair(fullChain, key)
					if err != nil {
						return nil, err
					}
					rpcCert = &newCert
				}
				return rpcCert, nil
			},
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  caPool,
		})),
	)
	appProto.RegisterAppServer(grpcServer, &Server{})
	return grpcServer
}

type Server struct {
	appProto.UnimplementedAppServer
}

func (s *Server) GetTokenStatus(ctx context.Context, req *appProto.TokenIDRequest) (*appProto.TokenStatus, error) {
	err := redis.NewRecordedToken().NewStorePoint(req.Id).Get(ctx, nil)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return &appProto.TokenStatus{
				Valid: false,
			}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &appProto.TokenStatus{
		Valid: true,
	}, nil
}

func (s *Server) WatchTokenOperation(_ *emptypb.Empty, srv appProto.App_WatchTokenOperationServer) error {
	// register listen channels
	canceledTokenChan := redis.NewCanceledTokenChannel().Subscribe(context.TODO()).Channel()
	// todo watch operation id change

	// load current data
	operationIDMap, err := redis.NewUserJwt().GetOperationTable(context.TODO())
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return status.Error(codes.Internal, err.Error())
		}
	}
	list, err := redis.NewCanceledToken().Get(context.TODO())
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return status.Error(codes.Internal, err.Error())
		}
	}
	var operationIDList = make([]*appProto.UserOperationID, 0, len(operationIDMap))
	for uid, operationID := range operationIDMap {
		operationIDList = append(operationIDList, &appProto.UserOperationID{
			Uid:         uint64(uid),
			OperationId: operationID,
		})
	}
	var canceledTokenList = make([]uint64, len(list))
	for i, v := range list {
		canceledTokenList[i] = v.ID
	}

	// send current data to client
	if err = srv.Send(&appProto.TokenOperation{
		UserOperation:   operationIDList,
		CanceledTokenId: canceledTokenList,
	}); err != nil {
		return err
	}

	for {
		select {
		case msg, ok := <-canceledTokenChan:
			if !ok {
				goto endStream
			}
			if msg.Payload != "" {
				msg.PayloadSlice = append(msg.PayloadSlice, msg.Payload)
			}

			canceledTokenList := make([]uint64, len(msg.PayloadSlice))
			for i, payload := range msg.PayloadSlice {
				var tokenCanceled redis.CanceledToken
				err := json.Unmarshal(unsafe.Slice(unsafe.StringData(payload), len(payload)), &tokenCanceled)
				if err != nil {
					return status.Error(codes.Internal, err.Error())
				}
				canceledTokenList[i] = tokenCanceled.ID
			}
			if err = srv.Send(&appProto.TokenOperation{
				CanceledTokenId: canceledTokenList,
			}); err != nil {
				return err
			}
		}
	}
endStream:
	return status.Error(codes.DataLoss, "send message failed")
}

func (s *Server) DestroyToken(ctx context.Context, req *appProto.RefreshTokenRequest) (*emptypb.Empty, error) {
	claims, valid, err := jwt.ParseRefreshToken(req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	} else if !valid {
		return nil, nil
	}

	if GetAuthInfo(ctx).AppCode != claims.AppCode {
		return nil, status.Error(codes.Unauthenticated, "token ownership not correct")
	}
	err = redis.NewRecordedToken().NewStorePoint(claims.ID).Destroy(ctx, claims.AppCode, claims.ExpiresAt.Time)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetUserInfo(_ context.Context, req *appProto.UserIDRequest) (*appProto.UserInfo, error) {
	user, err := service.User.UserInfoByID(uint(req.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}
	groups, err := service.UserGroups.GetNamesForUser(user.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}
	return &appProto.UserInfo{
		Uid:       req.Id,
		Name:      user.Name,
		AvatarUrl: user.AvatarUrl,
		Groups:    groups,
	}, nil
}
