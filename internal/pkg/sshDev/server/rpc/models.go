package rpc

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/proto"
)

type SshAccountMsg struct {
	IsDel     bool   `json:"isDel"`  // 是否为删除用户操作
	IsKill    bool   `json:"isKill"` // 是否为结束用户进程操作
	Username  string `json:"username"`
	PublicKey string `json:"publicKey"`
}

func (a SshAccountMsg) Rpc() *proto.SshAccount {
	return &proto.SshAccount{
		IsDel:     a.IsDel,
		IsKill:    a.IsKill,
		Username:  a.Username,
		PublicKey: a.PublicKey,
	}
}
