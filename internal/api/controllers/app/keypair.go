package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/keypair"
)

func ServerPublicKeys(c *gin.Context) {
	jwtPublic, err := keypair.PemMarshalPublic(global.JwtEd25519.PublicKey)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callback.Success(c, gin.H{
		"jwt": jwtPublic,
	})
}
