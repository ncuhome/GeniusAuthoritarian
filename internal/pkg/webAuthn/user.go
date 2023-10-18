package webAuthn

import (
	"fmt"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"unsafe"
)

type User struct {
	ID          []byte
	Name        string
	Credentials []webauthn.Credential
}

func NewUser(model *dao.User) (User, error) {
	idStr := fmt.Sprint(model.ID)
	idBytes := unsafe.Slice(unsafe.StringData(idStr), len(idStr))

	var cred []webauthn.Credential
	return User{
		ID:          idBytes,
		Name:        model.Name,
		Credentials: cred,
	}, nil //TODO load creds
}

func (u User) WebAuthnID() []byte {
	return u.ID
}

func (u User) WebAuthnName() string {
	return u.Name
}

func (u User) WebAuthnDisplayName() string {
	return u.Name
}

func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func (u User) WebAuthnIcon() string {
	// 此特性已经在新规范中弃用，返回空字符串即可
	return ""
}
