package config

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/infrastructure"
	"github.com/kyawmyintthein/orange-contrib/logx"
	"github.com/kyawmyintthein/orange-contrib/tracingx/newrelicx"
)

type GeneralConfig struct {
	App              AppCfg                              `mapstructure:"app" json:"app"`
	GracefulShutdown GracefulShutdownCfg                 `mapstructure:"graceful_shutdown" json:"graceful_shutdown"`
	NRTracer         newrelicx.NewrelicCfg               `mapstructure:"nr_tracer" json:"nr_tracer"`
	Log              logx.LogCfg                         `mapstructure:"log" json:"log"`
	Swagger          SwaggerCfg                          `mapstructure:"swagger" json:"swagger"`
	MongoDB          infrastructure.MongodbConfig        `mapstructure:"mongodb" json:"mongodb"`
	MysqlDB          infrastructure.SqlDBConfig          `mapstructure:"mysqldb" json:"mysqldb"`
	MongoBaseRepo    infrastructure.MongodbRepositoryCfg `mapstructure:"mongo_base_repo" json:"mongo_base_repo"`
	SqlBaseRepo      infrastructure.SqlRepositoryCfg     `mapstructure:"sqx_base_repo" json:"sqx_base_repo"`
}
