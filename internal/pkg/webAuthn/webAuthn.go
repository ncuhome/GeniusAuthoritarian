package webAuthn

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	log "github.com/sirupsen/logrus"
)

var Client *webauthn.WebAuthn

func init() {
	var err error
	Client, err = webauthn.New(&webauthn.Config{
		RPDisplayName: global.ThisAppName,
		RPID:          global.Config.WebAuthn.ID,
		RPOrigins:     []string{global.Config.WebAuthn.Origin},
	})
	if err != nil {
		log.Fatalln("webauthn init failed:", err)
	}
}

func Options() {
	Client.BeginRegistration()
}
