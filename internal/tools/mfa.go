package tools

import (
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"time"
)

func NewMfa(uid uint) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "GeniusAuth",
		AccountName: fmt.Sprint(uid),
		Algorithm:   otp.AlgorithmSHA512,
	})
}

func VerifyMfa(code, secret string) (bool, error) {
	return totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Algorithm: otp.AlgorithmSHA512,
	})
}
