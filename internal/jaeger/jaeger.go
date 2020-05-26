package jaeger

import (
	"context"
	"fmt"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/option"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	"io"
	"net/http"
	"strings"
)

type Jaeger interface {
	MiddlewareTracer(http.Handler) http.Handler
	HttpClientTracer(context.Context, *http.Request, string) opentracing.Span
	Close() error
	IsEnabled() bool
}

type jaegerService struct {
	cfg    *JaegerCfg
	tracer opentracing.Tracer
	closer io.Closer
	logger logging.KVLogger
}

func New(cfg *JaegerCfg, opts ...option.Option) (Jaeger, error) {
	options := option.NewOptions(opts...)

	jaegerService := jaegerService{
		cfg: cfg,
	}

	jaegerConfiguration := &jconfig.Configuration{
		Sampler: &jconfig.SamplerConfig{
			Type:              cfg.SamplerType,
			Param:             cfg.SamplerParam,
			SamplingServerURL: cfg.SamplingServerURL,
		},
		Reporter: &jconfig.ReporterConfig{
			LogSpans:           cfg.LogSpans,
			CollectorEndpoint:  cfg.ReporterCollectorEndpoint,
			LocalAgentHostPort: cfg.LocalAgentPort,
		},
	}
	jaegerConfiguration.ServiceName = cfg.LocalServiceName
	jaegerConfiguration.Sampler.Type = cfg.SamplerType
	jaegerConfiguration.Sampler.Param = cfg.SamplerParam

	// set custom logger
	kvLogger, ok := options.Context.Value(loggerKey{}).(logging.KVLogger)
	if kvLogger != nil || ok {
		jaegerService.logger = kvLogger
	} else {
		jaegerService.logger = logging.DefaultKVLogger()
	}

	// set jaeger logger
	jaegerLogger, ok := options.Context.Value(jaegerLoggerKey{}).(jaeger.Logger)
	if jaegerLogger != nil || !ok {
		jaegerLogger = jaeger.StdLogger
		jaegerService.logger.Info("[CLJaeger] Jaeger is using jaeger.StdLogger")
	}

	tracer, closer, err := jaegerConfiguration.NewTracer(
		jconfig.Logger(jaegerLogger),
	)

	if err != nil {
		jaegerService.cfg.Enabled = false
		jaegerService.logger.Errorf(err, "[CLJaeger] failed to initiate Jaeger")
		return &jaegerService, err
	}

	jaegerService.tracer = tracer
	jaegerService.closer = closer

	jaegerService.logger.Info("[CLJaeger] Jaeger is successfully initiated")
	return &jaegerService, nil
}

func (jaegerService *jaegerService) HttpClientTracer(ctx context.Context, req *http.Request, operationName string) opentracing.Span {
	if !jaegerService.cfg.Enabled {
		return nil
	}
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, jaegerService.tracer, operationName)
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, req.URL.String())
	ext.HTTPMethod.Set(span, req.Method)
	for k, v := range req.Header {
		span.SetTag(fmt.Sprintf("request.header.%s", strings.ToLower(k)), v)
	}

	req = req.WithContext(ctx)
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	jaegerService.logger.Info("Span Injected")
	return span
}

func (jaegerService *jaegerService) MiddlewareTracer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !jaegerService.cfg.Enabled {
			return
		}
		requestSpan := jaegerService.newRequestSpan(r)
		r = r.WithContext(opentracing.ContextWithSpan(r.Context(), requestSpan))
		w = &statusCodeTracker{w, 200}

		//traceID := r.Header.Get("Circle-Trace-Id")
		//if traceID != "" {
		//	traceID = strings.Split(traceID, ":")[0]
		//	w.Header().Set("X-Trace-Id", traceID)
		//	// adds traceID to a context and get from it latter
		//	r = r.WithContext(jaegerService.tracer.WithContextValue(r.Context(), traceID))
		//}
		next.ServeHTTP(w, r)
		for k, v := range w.Header() {
			requestSpan.SetTag(fmt.Sprintf("response.header.%s", strings.ToLower(k)), v)
		}
		ext.HTTPStatusCode.Set(requestSpan, uint16(w.(*statusCodeTracker).status))
		defer requestSpan.Finish()
	})
}

func (jaegerService *jaegerService) newRequestSpan(r *http.Request) opentracing.Span {
	var span opentracing.Span
	operation := fmt.Sprintf("HTTP %s %s", r.Method, r.URL.String())
	carrier := opentracing.HTTPHeadersCarrier(r.Header)
	wireContext, err := jaegerService.tracer.Extract(opentracing.HTTPHeaders, carrier)
	if err != nil {
		span = jaegerService.tracer.StartSpan(operation)
	} else {
		span = jaegerService.tracer.StartSpan(operation, opentracing.ChildOf(wireContext))
	}

	// it adds the trace ID to the http headers
	if err := span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier); err != nil {
		ext.Error.Set(span, true)
	} else {
		r = r.WithContext(opentracing.ContextWithSpan(r.Context(), span))
	}

	ext.HTTPMethod.Set(span, r.Method)
	ext.HTTPUrl.Set(span, r.URL.String())
	for k, v := range r.Header {
		span.SetTag(fmt.Sprintf("request.header.%s", strings.ToLower(k)), v)
	}

	return span
}

func (jaegerService *jaegerService) Close() error {
	if !jaegerService.cfg.Enabled {
		return nil
	}
	jaegerService.cfg.Enabled = true
	return jaegerService.closer.Close()
}

func (jaegerService *jaegerService) IsEnabled() bool {
	return jaegerService.cfg.Enabled
}
