package infrastructure

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kyawmyintthein/orange-contrib/optionx"
	"github.com/kyawmyintthein/orange-contrib/tracingx/newrelicx"
	newrelic "github.com/newrelic/go-agent"
)

type MongodbRepositoryCfg struct {
	CustomIDPrefix  string `mapstructure:"custom_id_prefix" json:"custom_id_prefix"`
	EnabledNRTracer bool   `mapstructure:"enabled_newrelic_tracer" json:"enabled_newrelic_tracer"`
}

type BaseMongoRepo struct {
	Config           *MongodbRepositoryCfg
	NRTracer         newrelicx.NewrelicTracer
	MongodbConnector MongodbConnector
}

func NewBaseMongoRepo(cfg *MongodbRepositoryCfg, mongodbConnector MongodbConnector, opts ...optionx.Option) *BaseMongoRepo {
	options := optionx.NewOptions(opts...)
	base := &BaseMongoRepo{
		Config:           cfg,
		MongodbConnector: mongodbConnector,
	}
	// set logger
	newRelicTracer, ok := options.Context.Value(newrelicTracerKey{}).(newrelicx.NewrelicTracer)
	if newRelicTracer == nil && !ok {
		base.Config.EnabledNRTracer = false
	} else {
		base.NRTracer = newRelicTracer
		base.Config.EnabledNRTracer = true
	}
	return base
}

func (base *BaseMongoRepo) GenerateID(ctx context.Context) string {
	return base.Config.CustomIDPrefix + "-" + strings.ToUpper(strconv.FormatInt(time.Now().UnixNano(), 36))
}

func (base *BaseMongoRepo) NrSegment(ctx context.Context, collName string, opName string) newrelic.DatastoreSegment {
	txn := newrelic.FromContext(ctx)
	dbConfig := base.MongodbConnector.Config()
	return base.NRTracer.RecordDatastoreMetric(txn,
		newrelic.DatastoreMongoDB,
		fmt.Sprintf("Collection.%s", collName),
		opName,
		"",
		nil,
		dbConfig.DatabaseHosts,
		"",
		dbConfig.DatabaseName)
}
