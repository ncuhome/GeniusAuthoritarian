package sshDev

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/proto"
	"sync/atomic"
)

type SshAccountMsg struct {
	IsDel     bool
	Username  string
	PublicKey string
}

func (a SshAccountMsg) Rpc() *proto.SshAccount {
	return &proto.SshAccount{
		IsDel:     a.IsDel,
		Username:  a.Username,
		PublicKey: a.PublicKey,
	}
}

type RpcSshAccountListElement struct {
	Channel  chan []SshAccountMsg
	IsQuited atomic.Bool
}
