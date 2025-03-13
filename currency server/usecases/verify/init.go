package verify

import "crypto/rsa"

type VerifyUseCase struct {
	publicKey *rsa.PublicKey
}

func NewVerifyUseCase(key *rsa.PublicKey) (*VerifyUseCase, error) {
	return &VerifyUseCase{publicKey: key}, nil
}
