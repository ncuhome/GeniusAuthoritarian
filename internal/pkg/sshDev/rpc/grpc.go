package rpc

import (
	"container/list"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/proto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"sync"
)

var MsgChannel chan []SshAccountMsg

func Run(token string) error {
	tcpListen, err := net.Listen("tcp", ":80")
	if err != nil {
		return err
	}

	MsgChannel = make(chan []SshAccountMsg)
	rpcSshAccounts := SshAccounts{}
	go rpcSshAccounts.Broadcaster()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(TokenAuthUnary(token)),
		grpc.StreamInterceptor(TokenAuthStream(token)),
	)
	proto.RegisterSshAccountsServer(grpcServer, &rpcSshAccounts)

	return grpcServer.Serve(tcpListen)
}

type SshAccounts struct {
	proto.UnimplementedSshAccountsServer

	list     list.List // *SshAccountListElement
	listLock sync.Mutex
}

func (a *SshAccounts) Watch(_ *emptypb.Empty, server proto.SshAccounts_WatchServer) error {
	// 发送现有账号
	sshAccounts, err := service.UserSsh.GetAllExist()
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	for _, account := range sshAccounts {
		err = server.Send(&proto.SshAccount{
			Username:  LinuxAccountName(account.UID),
			PublicKey: account.PublicSsh,
		})
		if err != nil {
			return err
		}
	}

	msgChan, unregister := a.RegisterWatcher()
	defer unregister()
	for {
		messages := <-msgChan
		for _, msg := range messages {
			err := server.Send(msg.Rpc())
			if err != nil {
				return err
			}
		}
	}
}

func (a *SshAccounts) RegisterWatcher() (chan []SshAccountMsg, func()) {
	channel := make(chan []SshAccountMsg, 2) // 添加 buffer 防死锁
	elContent := SshAccountListElement{
		Channel: channel,
	}

	a.listLock.Lock()
	el := a.list.PushBack(&elContent)
	a.listLock.Unlock()

	return channel, func() {
		/*
			先添加退出标记后清空管道
			由于广播是单线程操作，就算正好未及时读取到退出标记也能正常工作
		*/

		elContent.IsQuited.Store(true)

		// 清空管道
		for {
			select {
			case <-channel:
			default:
				goto removeElement
			}
		}

	removeElement:
		a.listLock.Lock()
		defer a.listLock.Unlock()
		a.list.Remove(el)
	}
}

func (a *SshAccounts) Broadcaster() {
	for {
		messages := <-MsgChannel

		a.listLock.Lock()
		el := a.list.Front()
		for el != nil {
			elContent := el.Value.(*SshAccountListElement)
			if !elContent.IsQuited.Load() {
				elContent.Channel <- messages
			}
			el = el.Next()
		}
		a.listLock.Unlock()
	}
}
