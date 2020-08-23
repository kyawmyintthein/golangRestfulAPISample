package newrelicx

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"

	"github.com/kyawmyintthein/orange-contrib/logx"
	"github.com/kyawmyintthein/orange-contrib/optionx"
	newrelic "github.com/newrelic/go-agent"
)

type NewrelicTracer interface {
	RecordHandlerMetric(string, func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request))
	RecordDatastoreMetric(newrelic.Transaction, newrelic.DatastoreProduct, string, string, string, map[string]interface{}, string, string, string) newrelic.DatastoreSegment
	RecordExternalMetric(*http.Request, string) (*newrelic.ExternalSegment, error)
	RecordCustomMetric(string, float64)
	RecordFunctionMetric(newrelic.Transaction, string) *newrelic.Segment
	RecordBackgroundProcessMetric(string) (newrelic.Transaction, error)
	IsEnabled() bool
	GetNewRelicApp() newrelic.Application
	ChiMiddleware(next http.Handler) http.Handler
}

type newrelicTracer struct {
	cfg     *NewrelicCfg
	options optionx.Options
	app     newrelic.Application
	logger  logx.Logger
}

func New(cfg *NewrelicCfg, opts ...optionx.Option) (NewrelicTracer, error) {
	options := optionx.NewOptions(opts...)
	newrelicTracer := &newrelicTracer{
		cfg:     cfg,
		options: options,
	}

	// set logger
	lgr, ok := options.Context.Value(loggerKey{}).(logx.Logger)
	if lgr == nil && !ok {
		newrelicTracer.logger = logx.New(&logx.LogCfg{})
	} else {
		newrelicTracer.logger = lgr
	}

	nrConfig := newrelic.NewConfig(cfg.Name, cfg.License)
	app, err := newrelic.NewApplication(nrConfig)
	if err != nil {
		newrelicTracer.cfg.Enabled = false
		newrelicTracer.logger.Errorf(context.Background(), err, "[%s] failed to initiate new-relic tracer", PackageName)
		return newrelicTracer, err
	}
	newrelicTracer.app = app
	newrelicTracer.cfg.Enabled = true
	newrelicTracer.logger.Infof(context.Background(), "[%s] Newrelic Tracer is successfully initiated", PackageName)
	return newrelicTracer, nil
}

func (service *newrelicTracer) ChiMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		urlPattern := getURLPattern(r)
		_, ok := service.cfg.SkipURLs[urlPattern]
		if !ok {
			logx.Infof(r.Context(), "URL: '%s' is skip", urlPattern)
		}
		if service.IsEnabled() && ok {
			txn := ((service.app).StartTransaction(r.URL.Path, w, r)).(newrelic.Transaction)
			defer txn.End()
			r = newrelic.RequestWithTransactionContext(r, txn)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func getURLPattern(req *http.Request) string {
	rctx := chi.RouteContext(req.Context())
	patterns := append(rctx.RoutePatterns[:0:0], rctx.RoutePatterns...)

	for i, pattern := range patterns[:len(patterns)-1] {
		patterns[i] = pattern[:len(pattern)-2]
	}
	fullPattern := strings.Join(patterns, "")
	return fmt.Sprintf("[%s]:%s", strings.ToLower(req.Method), fullPattern)
}

func (service *newrelicTracer) RecordHandlerMetric(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	if service.cfg.Enabled || service.app != nil {
		return newrelic.WrapHandleFunc(service.app, pattern, handler)
	}
	return pattern, handler
}

func (service *newrelicTracer) RecordDatastoreMetric(txn newrelic.Transaction, product newrelic.DatastoreProduct,
	collection string, operation string, parameterizedQuery string, queryParameters map[string]interface{}, host string,
	portPathOrID string, databaseName string) newrelic.DatastoreSegment {
	s := newrelic.DatastoreSegment{
		Product:            product,
		Collection:         collection,
		Operation:          operation,
		ParameterizedQuery: parameterizedQuery,
		QueryParameters:    queryParameters,
		Host:               host,
		PortPathOrID:       portPathOrID,
		DatabaseName:       databaseName,
	}
	s.StartTime = newrelic.StartSegmentNow(txn)
	return s
}

func (service *newrelicTracer) RecordExternalMetric(r *http.Request, opName string) (*newrelic.ExternalSegment, error) {
	if !service.cfg.Enabled || service.app == nil {
		return nil, NotAvailableError()
	}

	txn := newrelic.FromContext(r.Context())
	if txn == nil {
		txn = service.app.StartTransaction(opName, nil, r)
	}
	txn.SetName(opName)
	return newrelic.StartExternalSegment(txn, r), nil
}

func (service *newrelicTracer) RecordCustomMetric(transactionName string, value float64) {
	if service.cfg.Enabled || service.app != nil {
		return
	}
	service.app.RecordCustomMetric(transactionName, value)
}

func (service *newrelicTracer) RecordBackgroundProcessMetric(transactionName string) (newrelic.Transaction, error) {
	if !service.cfg.Enabled || service.app == nil {
		return nil, NotAvailableError()
	}
	return service.app.StartTransaction(transactionName, nil, nil), nil
}

func (service *newrelicTracer) RecordFunctionMetric(txn newrelic.Transaction, fnName string) *newrelic.Segment {
	return newrelic.StartSegment(txn, fnName)
}

func (service *newrelicTracer) IsEnabled() bool {
	return service.cfg.Enabled
}

func (service *newrelicTracer) GetNewRelicApp() newrelic.Application {
	if service.app != nil {
		return service.app
	}
	return nil
}
