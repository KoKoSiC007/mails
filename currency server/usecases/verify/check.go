package verify

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func (u *VerifyUseCase) Check(data string, signature []byte) error {
	hash := sha256.Sum256([]byte(data))

	err := rsa.VerifyPSS(
		u.publicKey,
		crypto.SHA256,
		hash[:],
		signature,
		nil,
	)
	if err != nil {
		fmt.Errorf("Invalid signature")
	}

	return nil
}
