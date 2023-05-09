package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
)

func CheckSignature(signature, secret string, data any) (bool, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, err
	}
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(jsonData)
	signatureGen := h.Sum(nil)

	return string(signatureGen) == signature, nil
}
