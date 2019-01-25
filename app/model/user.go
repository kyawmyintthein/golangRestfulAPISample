package model

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"golang.org/x/net/context"
)

type User struct{
	RawID             objectid.ObjectID `json:"raw_id,omitempty" bson:"_id,omitempty"`
	ID                string             `json:"id,omitempty" bson:"id,omitempty"`
	Name string `json:"name" bson:"name" validate:"nonzero"`
}

type CreateUserRequest struct{
	Name string `json:"name" validate:"nonzero"`
}

func (a *CreateUserRequest) Validate(ctx context.Context) error {
	return ValidateFields(a)
}

func (rm *CreateUserRequest) Convert(user *User){
	user.Name = rm.Name
}

type UpdateUserRequest struct{
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (a *UpdateUserRequest) Validate(ctx context.Context) error {
	return ValidateFields(a)
}

func (rm *UpdateUserRequest) Convert(user *User){
	user.ID = rm.ID
	user.Name = rm.Name
}