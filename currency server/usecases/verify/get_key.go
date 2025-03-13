package verify

import (
	"os"
)

func (u *VerifyUseCase) GetKey() (*[]byte, error) {
	pemDate, err := os.ReadFile("rsa_public_dev.pem")
	if err != nil {
		return nil, err
	}

	return &pemDate, nil
}
