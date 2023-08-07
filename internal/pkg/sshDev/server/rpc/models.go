package rpc

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/proto"
	"sync/atomic"
)

type SshAccountMsg struct {
	IsDel     bool // 是否为删除用户操作
	IsKill    bool // 是否为结束用户进程操作
	Username  string
	PublicKey string
}

func (a SshAccountMsg) Rpc() *proto.SshAccount {
	return &proto.SshAccount{
		IsDel:     a.IsDel,
		IsKill:    a.IsKill,
		Username:  a.Username,
		PublicKey: a.PublicKey,
	}
}

type SshAccountListElement struct {
	Channel  chan []SshAccountMsg
	IsQuited atomic.Bool
}
