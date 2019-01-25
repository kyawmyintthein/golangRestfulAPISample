package service

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/ecodes"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/datastore"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/model"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/errors"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type UserServiceInterface interface{
	Create(ctx context.Context, userRequest *model.CreateUserRequest) (*model.User, error)
	Update(ctx context.Context, userRequest *model.UpdateUserRequest) (*model.User, error)
	DeleteByID(ctx context.Context, ID string) error
	FindByID(ctx context.Context, ID string) (*model.User, error)
	FindAll(ctx context.Context) ([]*model.User, error)
}


type UserService struct{
	Config *config.GeneralConfig
	Logging logging.Logger
	UserDatastore datastore.UserDatastoreInterface
}

func (s *UserService) Create(ctx context.Context, userRequest *model.CreateUserRequest) (*model.User, error){
	user := &model.User{}
	userRequest.Convert(user)
	err := s.UserDatastore.Create(ctx, user)
	if err != nil{
		if datastore.IsDuplicateError(err) {
			err = errors.Wrap(err, ecodes.DuplicateUser, constant.DuplicateUser)
			return user, err
		}
		err = errors.Wrap(err, ecodes.InternalServerError, constant.ServerIssue)
		return user, err
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, userRequest *model.UpdateUserRequest) (*model.User, error){
	user := &model.User{}
	userRequest.Convert(user)
	objID, _ := objectid.FromHex(user.ID)
	err := s.UserDatastore.Update(ctx,  bson.D{{"_id", objID}}, user)
	if err != nil{
		return user, err
	}
	user.RawID = objID
	return user, nil
}

func (s *UserService) FindByID(ctx context.Context, id string) (*model.User, error){
	user := &model.User{}
	objID, _ := objectid.FromHex(id)
	err := s.UserDatastore.FindOne(ctx, bson.D{{"_id", objID}}, user)
	if err != nil{
		if datastore.IsNotFoundError(err){
			err = errors.Wrap(err, ecodes.UserNotFound, constant.UserNotFound)
			return user, err
		}
		err = errors.Wrap(err, ecodes.InternalServerError, constant.ServerIssue)
		return user, err
	}
	return user, nil
}

func (s *UserService) DeleteByID(ctx context.Context, ID string) error{
	err := s.UserDatastore.Delete(ctx, bson.M{"_id": ID})
	if err != nil{
		err = errors.Wrap(err, ecodes.InternalServerError, constant.ServerIssue)
		return err
	}
	return nil
}


func (s *UserService) FindAll(ctx context.Context) ([]*model.User, error){
	users, err := s.UserDatastore.FindAll(ctx, bson.D{})
	if err != nil{
		if datastore.IsNotFoundError(err){
			err = errors.Wrap(err, ecodes.UserNotFound, constant.UserNotFound)
			return users, err
		}
		err = errors.Wrap(err, ecodes.InternalServerError, constant.ServerIssue)
		return users, err
	}
	return users, nil
}