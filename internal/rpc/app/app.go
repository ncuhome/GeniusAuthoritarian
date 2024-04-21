// rpc server for applications
// mutual ssl verify is used

package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/app/appProto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"unsafe"
)

type Server struct {
	appProto.UnimplementedAppServer
}

func (s *Server) GetTokenStatus(ctx context.Context, req *appProto.TokenRequest) (*appProto.TokenStatus, error) {
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

func (s *Server) GetTokenCanceled(_ *emptypb.Empty, srv appProto.App_GetTokenCanceledServer) error {
	msgChan := redis.NewCanceledTokenChannel().Subscribe(context.TODO()).Channel()
	list, err := redis.NewCanceledToken().Get(context.TODO())
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return status.Error(codes.Internal, err.Error())
		}
	}
	var listId = make([]uint64, len(list))
	for i, v := range list {
		listId[i] = v.ID
	}
	if err = srv.Send(&appProto.TokenCanceled{
		Id: listId,
	}); err != nil {
		return err
	}

	for {
		msg, ok := <-msgChan
		if !ok {
			break
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
		if err = srv.Send(&appProto.TokenCanceled{
			Id: listId,
		}); err != nil {
			return err
		}
	}
	return status.Error(codes.DataLoss, "send message failed")
}

func (s *Server) DestroyToken(ctx context.Context, req *appProto.TokenRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "todo")
}

func (s *Server) GetUserInfo(_ context.Context, req *appProto.TokenRequest) (*appProto.UserInfo, error) {
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
