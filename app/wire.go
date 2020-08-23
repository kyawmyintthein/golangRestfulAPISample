//+build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/delivery/api"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/injectors"
	mongoRepo "github.com/kyawmyintthein/golangRestfulAPISample/app/repository/mongo_repository"
	mysqlRepo "github.com/kyawmyintthein/golangRestfulAPISample/app/repository/mysql_repository"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/usecase"
	"github.com/kyawmyintthein/golangRestfulAPISample/infrastructure"
)

func NewApp(configFilePaths ...string) (*restApiApplication, error) {
	panic(wire.Build(
		injectors.ProvideConfig,
		injectors.ProvideLogger,
		injectors.ProvideRouter,
		injectors.ProvideHttpServer,
		injectors.ProvideNewRelic,
		infrastructure.NewHttpResponseWriter,
		infrastructure.ProvideStringHelper,

		injectors.ProvideSqlDBConnector,
		injectors.ProvideMongoDBConnector,

		injectors.ProvideBaseSqlRepo,
		injectors.ProvideBaseMongoRepo,

		mysqlRepo.ProvideUserRepository,
		mongoRepo.ProvideArticleRepository,

		usecase.ProvideUserUsecase,
		usecase.ProvideArticleUsecase,

		api.ProvideBaseHandler,
		api.ProvideHealthHandler,
		api.ProvideShutdownHandler,
		api.ProvideUserHandler,
		api.ProvideArticleHandler,
		wire.Struct(new(restApiApplication), "*"),
	))
}
