package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/jaeger"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/newrelic"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/option"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"time"
)

const (
	defaultRequestTimeout  time.Duration = 10
	defaultRetryAttempts   uint          = 0
	defaultBackOffDuration time.Duration = 1 * time.Second
)

type HttpClient interface {
	JSONPost(context.Context, interface{}, map[string]string, string, string, ...option.Option) (*http.Response, error)
	JSONPut(context.Context, interface{}, map[string]string, string, string, ...option.Option) (*http.Response, error)
	JSONGet(context.Context, map[string]string, string, string, ...option.Option) (*http.Response, error)
	JSONDelete(context.Context, interface{}, map[string]string, string, string, ...option.Option) (*http.Response, error)
}

type httpClient struct {
	options        option.Options
	logger         logging.KVLogger
	newrelicTracer newrelic.NewrelicTracer
	jaeger         jaeger.Jaeger
}

func NewHttpClient(opts ...option.Option) HttpClient {
	options := option.NewOptions(opts...)

	httpClient := &httpClient{
		options: options,
	}

	//set newrelic
	newrelicTracer, ok := options.Context.Value(newrelicKey{}).(newrelic.NewrelicTracer)
	if newrelicTracer != nil && ok {
		httpClient.newrelicTracer = newrelicTracer
	}

	// set jaeger
	jaeger, ok := options.Context.Value(jaegerKey{}).(jaeger.Jaeger)
	if jaeger != nil || !ok {
		httpClient.jaeger = jaeger
	}

	// set logger
	lgr, ok := options.Context.Value(loggerKey{}).(logging.KVLogger)
	if lgr != nil && ok {
		httpClient.logger = lgr
	} else {
		httpClient.logger = logging.DefaultKVLogger()
	}
	httpClient.logger.Info("[CLHttpClient] HttpClient is successfully initiated")
	return httpClient
}

func (httpClient *httpClient) JSONDelete(ctx context.Context, payload interface{}, headers map[string]string, url string, operationName string, opts ...option.Option) (*http.Response, error) {
	options := option.NewOptions(opts...)
	retryConfig := httpClient.getRetryCfgFromOption(options)

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for hk, hv := range headers {
			req.Header.Set(hk, hv)
		}
	}

	//TODO: improvement
	var span opentracing.Span
	if httpClient.jaeger != nil && httpClient.jaeger.IsEnabled() {
		span = httpClient.jaeger.HttpClientTracer(ctx, req, operationName)
		defer span.Finish()
	}

	resp, err := httpClient.firstAttemptAndRetry(retryConfig, req, operationName)
	if err != nil {
		return nil, err
	}

	//TODO: improvement
	if httpClient.jaeger != nil && httpClient.jaeger.IsEnabled() && span != nil {
		span.SetTag("http.response.status", resp.StatusCode)
		for k, v := range resp.Header {
			span.SetTag(fmt.Sprintf("http.response.header.%s", k), v)
		}
	}

	logger := httpClient.getLogger(ctx)
	logger.InfoKV(logging.KV{"URL": url, "Status": resp.Status, "Headers": resp.Header}, "[JSONDelete] Received response")
	return resp, err
}

func (httpClient *httpClient) JSONPost(ctx context.Context, payload interface{}, headers map[string]string, url string, operationName string, opts ...option.Option) (*http.Response, error) {
	options := option.NewOptions(opts...)
	retryConfig := httpClient.getRetryCfgFromOption(options)

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for hk, hv := range headers {
			req.Header.Set(hk, hv)
		}
	}

	//TODO: improve
	var span opentracing.Span
	if httpClient.jaeger != nil && httpClient.jaeger.IsEnabled() {
		span = httpClient.jaeger.HttpClientTracer(ctx, req, operationName)
		defer span.Finish()
	}

	resp, err := httpClient.firstAttemptAndRetry(retryConfig, req, operationName)
	if err != nil {
		return resp, err
	}

	//TODO: improvement
	if httpClient.jaeger != nil && httpClient.jaeger.IsEnabled() && span != nil {
		span.SetTag("http.response.status", resp.StatusCode)
		for k, v := range resp.Header {
			span.SetTag(fmt.Sprintf("http.response.header.%s", k), v)
		}
	}

	logger := httpClient.getLogger(ctx)
	logger.InfoKV(logging.KV{"URL": url, "Status": resp.Status, "Headers": resp.Header}, "[JSONPost] Received response")
	return resp, err
}

func (httpClient *httpClient) JSONPut(ctx context.Context, payload interface{}, headers map[string]string, url string, operationName string, opts ...option.Option) (*http.Response, error) {
	options := option.NewOptions(opts...)
	retryConfig := httpClient.getRetryCfgFromOption(options)

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for hk, hv := range headers {
			req.Header.Set(hk, hv)
		}
	}

	//TODO: improvement
	var span opentracing.Span
	if httpClient.jaeger != nil && httpClient.jaeger.IsEnabled() {
		span = httpClient.jaeger.HttpClientTracer(ctx, req, operationName)
		defer span.Finish()
	}

	resp, err := httpClient.firstAttemptAndRetry(retryConfig, req, operationName)
	if err != nil {
		return nil, err
	}

	//TODO: improvement
	if httpClient.jaeger != nil && httpClient.jaeger.IsEnabled() && span != nil {
		span.SetTag("http.response.status", resp.StatusCode)
		for k, v := range resp.Header {
			span.SetTag(fmt.Sprintf("http.response.header.%s", k), v)
		}
	}

	logger := httpClient.getLogger(ctx)
	logger.InfoKV(logging.KV{"URL": url, "Status": resp.Status, "Headers": resp.Header}, "[JSONPut] Received response")
	return resp, err
}

func (httpClient *httpClient) JSONGet(ctx context.Context, headers map[string]string, url string, operationName string, opts ...option.Option) (*http.Response, error) {
	options := option.NewOptions(opts...)
	retryConfig := httpClient.getRetryCfgFromOption(options)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for hk, hv := range headers {
			req.Header.Set(hk, hv)
		}
	}

	//TODO: improvement
	var span opentracing.Span
	if httpClient.jaeger != nil && httpClient.jaeger.IsEnabled() {
		span = httpClient.jaeger.HttpClientTracer(ctx, req, operationName)
		defer span.Finish()
	}

	resp, err := httpClient.firstAttemptAndRetry(retryConfig, req, operationName)
	if err != nil {
		return nil, err
	}

	//TODO: improvement
	if httpClient.jaeger != nil && httpClient.jaeger.IsEnabled() && span != nil {
		span.SetTag("http.response.status", resp.StatusCode)
		for k, v := range resp.Header {
			span.SetTag(fmt.Sprintf("http.response.header.%s", k), v)
		}
	}

	logger := httpClient.getLogger(ctx)
	logger.InfoKV(logging.KV{"URL": url, "Status": resp.Status, "Headers": resp.Header}, "[JSONGet] Received response")
	return resp, err
}

func (httpClient *httpClient) firstAttemptAndRetry(retryConfig *RetryCfg, req *http.Request, operationName string) (*http.Response, error) {
	var count uint
	resp, err := httpClient.sendHttpRequest(req, operationName)
	if err != nil {
		if retryConfig.Enabled {
			return resp, err
		}

		for count < retryConfig.MaxRetryAttempts {
			resp, err := httpClient.sendHttpRequest(req, operationName)
			if err != nil {
				if count == retryConfig.MaxRetryAttempts-1 {
					return resp, err
				}

				backOffDuration := defaultBackOffDuration
				if uint(len(retryConfig.BackOffDurations)) >= count {
					backOffDuration = retryConfig.BackOffDurations[count]
				}

				time.Sleep(backOffDuration)
			} else {
				return resp, err
			}
			count++
		}
	}
	return resp, err
}

func (httpClient *httpClient) sendHttpRequest(req *http.Request, name string, opts ...option.Option) (*http.Response, error) {
	client := http.Client{Transport: &nethttp.Transport{}}
	requestTimeout, ok := httpClient.options.Context.Value(httpRequestTimeoutKey{}).(time.Duration)
	if !ok {
		requestTimeout = defaultRequestTimeout
	}
	client.Timeout = requestTimeout * time.Second

	if httpClient.newrelicTracer == nil || !httpClient.newrelicTracer.IsEnabled() {
		return client.Do(req)
	}

	es, err := httpClient.newrelicTracer.RecordExternalMetric(req, name)
	if err == nil {
		defer es.End()
	}

	response, err := client.Do(req)
	return response, err
}

func (httpClient *httpClient) getRetryCfgFromOption(options option.Options) *RetryCfg {
	// set retry config
	retryConfig, ok := options.Context.Value(retryConfig{}).(*RetryCfg)
	if retryConfig == nil || !ok {
		retryConfig = &RetryCfg{
			MaxRetryAttempts: defaultRetryAttempts,
		}
		var i uint = 0
		for i < defaultRetryAttempts {
			retryConfig.BackOffDurations = append(retryConfig.BackOffDurations, defaultBackOffDuration)
			i++
		}
	}

	if uint(len(retryConfig.BackOffDurations)) < retryConfig.MaxRetryAttempts {
		backOffLen := len(retryConfig.BackOffDurations)
		if backOffLen == 0 {
			var i uint = 0
			for i < defaultRetryAttempts {
				retryConfig.BackOffDurations = append(retryConfig.BackOffDurations, defaultBackOffDuration)
				i++
			}
		} else {
			missingAttempts := retryConfig.MaxRetryAttempts - uint(backOffLen)
			lastBackoffDuration := retryConfig.BackOffDurations[backOffLen-1]
			var i uint = 0
			for i < missingAttempts {
				retryConfig.BackOffDurations = append(retryConfig.BackOffDurations, lastBackoffDuration)
				i++
			}
		}
	}
	return retryConfig
}

func (httpClient *httpClient) getLogger(ctx context.Context) logging.KVLogger {
	logger := httpClient.logger
	structureLogger, isExist := logging.GetStructuredLoggerIfExist(ctx)
	if !isExist {
		logger = structureLogger.AsKVLogger()
	}
	return logger
}
