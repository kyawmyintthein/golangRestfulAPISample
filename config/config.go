package config

import (
	base_repository "github.com/kyawmyintthein/golangRestfulAPISample/internal/base-repository"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/mongo"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/newrelic"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/sql"
)

type GeneralConfig struct {
	SqlDB       sql.SqlDBConfig                  `mapstructure:"sqldb" json:"sqldb"`
	SqlBaseRepo base_repository.SqlRepositoryCfg `mapstructure:"sql_base_repo" json:"sql_base_repo"`

	MongoDB       mongo.MongodbConfig                  `mapstructure:"mongodb" json:"mongodb"`
	MongoBaseRepo base_repository.MongodbRepositoryCfg `mapstructure:"mongo_base_repo" json:"mongo_base_repo"`

	NRTracer         newrelic.NewrelicCfg `mapstructure:"nr_tracer" json:"nr_tracer"`
	App              AppCfg               `mapstructure:"app" json:"app"`
	GracefulShutdown GracefulShutdownCfg  `json:"graceful_shutdown" mapstructure:"graceful_shutdown"`
	Log              logging.LogCfg       `mapstructure:"log" json:"log"`
	Swagger          SwaggerCfg           `mapstructure:"swagger" json:"swagger"`
}
