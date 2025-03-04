package users

import (
	"odyssey/m/v2/models/user"
)

func (u *UserUseCase) Auth(email, password string) (*user.User, error) {
	result, err := u.userRepo.Get(email, password)
	if err != nil {
		return nil, err
	}

	user := &user.User{
		Id:        result.ID,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
	}

	return user, nil
}
