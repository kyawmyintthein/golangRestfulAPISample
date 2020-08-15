package injectors

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/mongo"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/sql"
)

func ProvideSqlDBConnector(config *config.GeneralConfig) (sql.SqlDBConnector, error) {
	return sql.NewSQLConnector(&config.SqlDB)
}

func ProvideMongoDBConnector(config *config.GeneralConfig) (mongo.MongodbConnector, error) {
	return mongo.NewMongodbConnector(&config.MongoDB)
}

