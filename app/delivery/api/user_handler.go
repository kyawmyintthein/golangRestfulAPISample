package api

import "github.com/kyawmyintthein/golangRestfulAPISample/app/service"

type UserHandler struct{
	UserService service.UserService
}

func ProvideUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}
