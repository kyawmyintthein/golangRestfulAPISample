package usecase

import (
	"context"

	"github.com/kyawmyintthein/golangRestfulAPISample/app/dto"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/repository"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/viewmodel"
)

type UserUsecase interface {
	RegisterNewUser(context.Context, *dto.RegisterNewUserDTO) (*viewmodel.UserVM, error)
}

type userUserase struct {
	userRepository repository.UserRepository
}

func ProvideUserUsecase(userRepository repository.UserRepository) UserUsecase {
	return &userUserase{userRepository: userRepository}
}

func (userService *userUserase) RegisterNewUser(ctx context.Context, registerNewUserDTO *dto.RegisterNewUserDTO) (*viewmodel.UserVM, error) {
	return &viewmodel.UserVM{}, nil
}
