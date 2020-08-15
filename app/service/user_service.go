package service

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/dto"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/repository"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/viewmodel"
)

type UserService interface{
	RegisterNewUser(context.Context, *dto.RegisterNewUserDTO) (*viewmodel.UserVM, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func ProvideUserService(userRepository repository.UserRepository) UserService{
	return &userService{userRepository: userRepository}
}

func (userService *userService) RegisterNewUser(ctx context.Context, registerNewUserDTO *dto.RegisterNewUserDTO) (*viewmodel.UserVM, error){
	return &viewmodel.UserVM{}, nil
}