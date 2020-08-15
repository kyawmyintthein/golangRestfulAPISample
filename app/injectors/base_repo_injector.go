package injectors

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	base_repository "github.com/kyawmyintthein/golangRestfulAPISample/internal/base-repository"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/mongo"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/newrelic"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/sql"
)

func ProvideBaseSqlRepo(config *config.GeneralConfig,
	sqlDBConnector sql.SqlDBConnector,
	nrTracer newrelic.NewrelicTracer) (*base_repository.BaseSqlRepository) {
	baseSqlRepo := &base_repository.BaseSqlRepository{
		Config: &config.SqlBaseRepo,
		NRTracer: nrTracer,
		DBConnector: sqlDBConnector,
	}
	return baseSqlRepo
}


func ProvideBaseMongoRepo(config *config.GeneralConfig,
	mongodbConnector mongo.MongodbConnector,
	nrTracer newrelic.NewrelicTracer) (*base_repository.BaseMongoRepo) {
	baseSqlRepo := &base_repository.BaseMongoRepo{
		Config: &config.MongoBaseRepo,
		NRTracer: nrTracer,
		MongodbConnector: mongodbConnector,
	}
	return baseSqlRepo
}
