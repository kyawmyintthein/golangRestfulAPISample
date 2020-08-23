package api

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/app/usecase"
)

type UserHandler struct {
	UserUsecase usecase.UserUsecase
}

func ProvideUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{UserUsecase: userUsecase}
}
