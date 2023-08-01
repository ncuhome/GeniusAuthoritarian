package sshDev

import (
	"container/list"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"sync"
)

var msgChannel chan []SshAccountMsg

func Run(token string) error {
	tcpListen, err := net.Listen("tcp", ":80")
	if err != nil {
		return err
	}

	msgChannel = make(chan []SshAccountMsg)
	rpcSshAccounts := RpcSshAccounts{}
	go rpcSshAccounts.Broadcaster()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(TokenAuth(token)))
	proto.RegisterSshAccountsServer(grpcServer, &rpcSshAccounts)

	return grpcServer.Serve(tcpListen)
}

type RpcSshAccounts struct {
	proto.UnimplementedSshAccountsServer

	list     list.List // *RpcSshAccountListElement
	listLock sync.Mutex
}

func (a *RpcSshAccounts) Watch(_ *emptypb.Empty, server proto.SshAccounts_WatchServer) error {
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

func (a *RpcSshAccounts) RegisterWatcher() (chan []SshAccountMsg, func()) {
	channel := make(chan []SshAccountMsg, 2) // 添加 buffer 防死锁
	elContent := RpcSshAccountListElement{
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

func (a *RpcSshAccounts) Broadcaster() {
	for {
		messages := <-msgChannel

		a.listLock.Lock()
		el := a.list.Front()
		for el != nil {
			elContent := el.Value.(*RpcSshAccountListElement)
			if !elContent.IsQuited.Load() {
				elContent.Channel <- messages
			}
			el = el.Next()
		}
		a.listLock.Unlock()
	}
}
