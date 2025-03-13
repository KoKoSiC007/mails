package verify

import "crypto/rsa"

type VerifyUseCase struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewVerifyUseCase(pubKey *rsa.PublicKey, privKey *rsa.PrivateKey) (*VerifyUseCase, error) {
	return &VerifyUseCase{publicKey: pubKey, privateKey: privKey}, nil
}
