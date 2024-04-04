package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/keypair"
	"time"
)

func ServerPublicKeys(c *gin.Context) {
	jwtPublic, err := keypair.PemMarshalPublic(keypair.FormatECDSA, global.JwtEd25519.PublicKey)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callback.Success(c, gin.H{
		"jwt": jwtPublic,
		"ca":  global.CaIssuer.CaCertPem,
	})
}

func RpcClientCredential(c *gin.Context) {
	appCode := tools.GetAppCode(c)

	validBefore := time.Now().AddDate(0, 0, 7)
	certPem, privatePem, err := global.CaIssuer.Issue([]string{appCode}, validBefore)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callback.Success(c, gin.H{
		"cert":        certPem,
		"key":         privatePem,
		"validBefore": validBefore.Unix(),
	})
}
