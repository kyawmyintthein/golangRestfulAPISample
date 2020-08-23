package mysql_repository

import (
	"context"

	"github.com/kyawmyintthein/golangRestfulAPISample/app/model"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/repository"
	"github.com/kyawmyintthein/golangRestfulAPISample/infrastructure"
)

type userMysqlRepo struct {
	*infrastructure.BaseSqlRepository
	defaultTable string
}

func ProvideUserRepository(baseSqlRepository *infrastructure.BaseSqlRepository) repository.UserRepository {
	return &userMysqlRepo{
		baseSqlRepository,
		_userTable,
	}
}

func (repo *userMysqlRepo) Create(ctx context.Context, user *model.User) (*model.User, error) {
	stmt, err := repo.InsertStatement(ctx, repo.defaultTable, user)
	if err != nil {
		return nil, err
	}
	result, err := repo.DBConnector.DB(ctx).NamedExecContext(ctx, stmt, user)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}
