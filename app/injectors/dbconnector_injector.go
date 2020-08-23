package injectors

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/infrastructure"
)

func ProvideSqlDBConnector(config *config.GeneralConfig) (infrastructure.SqlDBConnector, error) {
	return infrastructure.NewSQLConnector(&config.MysqlDB)
}

func ProvideMongoDBConnector(config *config.GeneralConfig) (infrastructure.MongodbConnector, error) {
	return infrastructure.NewMongodbConnector(&config.MongoDB)
}
