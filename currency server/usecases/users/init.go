package users

import "odyssey/m/v2/repositories"

type UserUseCase struct {
	userRepo *repositories.UserRepo
}

func NewUserCase(repo *repositories.UserRepo) (*UserUseCase, error) {
	return &UserUseCase{userRepo: repo}, nil
}
