package users

import (
	users "odyssey/m/v2/models/user"
)

func (u *UserUseCase) Create(newUser *users.NewUser) (*users.User, error) {
	result, err := u.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	user := &users.User{
		Id:        result.ID,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
	}

	return user, nil
}
