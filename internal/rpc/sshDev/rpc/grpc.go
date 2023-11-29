package rpc

import (
	"context"
	"encoding/json"
	"errors"
	redisPkg "github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/middlewares"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/sshDev/client/proto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/sshDev/rpcModel"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"time"
	"unsafe"
)

func Run(token, addr string) error {
	tcpListen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middlewares.UnaryLogger(),
			UnaryTokenAuth(token),
		),
		grpc.ChainStreamInterceptor(
			middlewares.StreamLogger(),
			StreamTokenAuth(token),
		),
	)
	proto.RegisterSshAccountsServer(grpcServer, &SshAccounts{})

	return grpcServer.Serve(tcpListen)
}

type SshAccounts struct {
	proto.UnimplementedSshAccountsServer
}

func (a *SshAccounts) Watch(_ *emptypb.Empty, server proto.SshAccounts_WatchServer) error {
	// 注册监听

	sub := redisPkg.SubScribeSshDev()
	defer sub.Close()

	// 发送现有账号
	sshAccounts, err := service.UserSsh.GetAllExist()
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	err = server.Send(&proto.AccountStream{
		IsInit:   true,
		Accounts: TransformAccountArray(sshAccounts),
	})
	if err != nil {
		return err
	}

	msgChannel := make(chan []rpcModel.SshAccountMsg)
	go func() {
		for {
			msg, err := sub.ReceiveMessage(context.Background())
			if err != nil {
				close(msgChannel)
				return
			}

			if msg.PayloadSlice == nil && msg.Payload != "" {
				msg.PayloadSlice = []string{msg.Payload}
			} else {
				continue
			}

			for _, payload := range msg.PayloadSlice {
				var msgDecoded []rpcModel.SshAccountMsg
				msgBytes := unsafe.Slice(unsafe.StringData(payload), len(payload))
				err = json.Unmarshal(msgBytes, &msgDecoded)
				msgChannel <- msgDecoded
			}
		}
	}()

	for {
		select {
		case messages, ok := <-msgChannel:
			if !ok {
				return errors.New("ssh account status subscription exception")
			}
			err := server.Send(&proto.AccountStream{
				Accounts: TransformMsgArray(messages),
			})
			if err != nil {
				return err
			}
		case <-time.After(time.Minute):
			err := server.Send(&proto.AccountStream{
				IsHeartBeat: true,
			})
			if err != nil {
				return err
			}
		}
	}
}
