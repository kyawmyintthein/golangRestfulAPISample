package repository

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/model"
)

type UserRepository interface {
	Create(context.Context, *model.User) (*model.User, error)
}
