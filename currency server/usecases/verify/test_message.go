package verify

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

type SignedMessage struct {
	Message   string `json:"message"`
	Signature []byte `json:"signature"`
}

func (u *VerifyUseCase) TestMessage() (*SignedMessage, error) {
	message := "Hello, Golang!"

	// Вычисляем хеш сообщения
	hash := sha256.Sum256([]byte(message))

	// Создаем подпись
	signature, err := rsa.SignPSS(
		rand.Reader,
		u.privateKey,
		crypto.SHA256,
		hash[:],
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &SignedMessage{Message: message, Signature: signature}, nil
}
