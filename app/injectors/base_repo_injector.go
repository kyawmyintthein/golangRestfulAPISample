package injectors

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/infrastructure"
	"github.com/kyawmyintthein/orange-contrib/tracingx/newrelicx"
)

func ProvideBaseSqlRepo(config *config.GeneralConfig,
	sqlDBConnector infrastructure.SqlDBConnector,
	nrTracer newrelicx.NewrelicTracer) *infrastructure.BaseSqlRepository {
	return infrastructure.NewBaseSqlRepository(&config.SqlBaseRepo, sqlDBConnector, infrastructure.WithNewrelicTracer(nrTracer))
}

func ProvideBaseMongoRepo(config *config.GeneralConfig,
	mongodbConnector infrastructure.MongodbConnector,
	nrTracer newrelicx.NewrelicTracer) *infrastructure.BaseMongoRepo {
	return infrastructure.NewBaseMongoRepo(&config.MongoBaseRepo, mongodbConnector, infrastructure.WithNewrelicTracer(nrTracer))
}
