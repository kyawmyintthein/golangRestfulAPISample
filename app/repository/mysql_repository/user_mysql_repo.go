package mysql_repository

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/model"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/repository"
	base_repository "github.com/kyawmyintthein/golangRestfulAPISample/internal/base-repository"
)

type userMysqlRepo struct{
	*base_repository.BaseSqlRepository
	defaultTable string
}

func ProvideUserRepository(baseSqlRepository *base_repository.BaseSqlRepository) repository.UserRepository{
	return &userMysqlRepo{
		baseSqlRepository,
		"users",
	}
}

func (repo *userMysqlRepo) Create(ctx context.Context, user *model.User) (*model.User, error){
	stmt, err := repo.InsertStatement(ctx, repo.defaultTable, user)
	if err != nil{
		return nil, err
	}
	result, err := repo.DBConnector.DB(ctx).NamedExecContext(ctx, stmt, user)
	if err != nil{
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil{
		return nil, err
	}
	user.ID = id
	return user, nil
}