package webAuthn

import (
	"fmt"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"unsafe"
)

type User struct {
	ID          []byte
	Name        string
	Credentials []webauthn.Credential
}

func NewUser(uid uint) (*User, error) {

	idStr := fmt.Sprint(uid)
	idBytes := unsafe.Slice(unsafe.StringData(idStr), len(idStr))

	name, err := service.WebAuthn.UserName(uid)
	if err != nil {
		return nil, err
	}

	cred, err := service.WebAuthn.GetCredentials(uid)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:          idBytes,
		Name:        name,
		Credentials: cred,
	}, nil
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

// Deprecated: 此特性已经在新规范中弃用，返回空字符串
func (u User) WebAuthnIcon() string {
	return ""
}
