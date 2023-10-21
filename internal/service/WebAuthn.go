package service

import (
	"encoding/json"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
	"unsafe"
)

var WebAuthn = WebAuthnSrv{dao.DB}

type WebAuthnSrv struct {
	*gorm.DB
}

func (a WebAuthnSrv) Begin() (WebAuthnSrv, error) {
	a.DB = a.DB.Begin()
	return a, a.Error
}

func (a WebAuthnSrv) UserName(id uint) (string, error) {
	model := dao.User{ID: id}
	return model.Name, model.FirstForPasskey(a.DB)
}

func (a WebAuthnSrv) Add(uid uint, cred *webauthn.Credential) (*dto.UserCredential, error) {
	credBytes, err := json.Marshal(cred)
	if err != nil {
		return nil, err
	}

	model := dao.UserWebauthn{
		UID:        uid,
		Credential: unsafe.String(unsafe.SliceData(credBytes), len(credBytes)),
	}
	return &dto.UserCredential{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
	}, model.Insert(a.DB)
}

func (a WebAuthnSrv) GetCredentials(uid uint) ([]webauthn.Credential, error) {
	data, err := (&dao.UserWebauthn{UID: uid}).GetByUID(a.DB)
	if err != nil {
		return nil, err
	}

	var cred = make([]webauthn.Credential, len(data))
	for i := 0; i < len(data); i++ {
		credStr := data[i]
		credBytes := unsafe.Slice(unsafe.StringData(credStr), len(credStr))
		err = json.Unmarshal(credBytes, &cred[i])
		if err != nil {
			return nil, err
		}
	}
	return cred, nil
}

func (a WebAuthnSrv) ListUserCredForShow(uid uint) ([]dto.UserCredential, error) {
	return (&dao.UserWebauthn{UID: uid}).GetByUidForShow(a.DB)
}
