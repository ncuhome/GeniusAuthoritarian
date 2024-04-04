package sshTool

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/ed25519"
)

func NewSshDevModel(uid uint) (model dao.UserSsh, err error) {
	sshKey, err := ed25519.Generate()
	if err != nil {
		return
	}

	publicPem, privatePem, err := sshKey.MarshalPem()
	if err != nil {
		return
	}
	publicSsh, privateSsh, err := sshKey.MarshalSSH()
	if err != nil {
		return
	}

	return dao.UserSsh{
		UID:        uid,
		PublicPem:  string(publicPem),
		PrivatePem: string(privatePem),
		PublicSsh:  string(publicSsh),
		PrivateSsh: string(privateSsh),
	}, nil
}
