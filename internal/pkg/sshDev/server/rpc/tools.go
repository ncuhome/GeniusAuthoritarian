package rpc

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/proto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/sshTool"
)

func TransformAccountArray(a []dto.SshDeploy) []*proto.SshAccount {
	var b = make([]*proto.SshAccount, len(a))
	for i, s := range a {
		b[i] = &proto.SshAccount{
			Username:  sshTool.LinuxAccountName(s.UID),
			PublicKey: s.PublicSsh,
		}
	}
	return b
}

func TransformMsgArray(a []SshAccountMsg) []*proto.SshAccount {
	var b = make([]*proto.SshAccount, len(a))
	for i, s := range a {
		b[i] = s.Rpc()
	}
	return b
}
