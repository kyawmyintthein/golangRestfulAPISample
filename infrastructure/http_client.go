package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/newrelic/go-agent"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"net/http"
	"time"
)

type HttpClient interface{
	JSONPost(ctx context.Context, payload interface{}, headers interface{}, url string, newrelicMetricName string) (*http.Response, error)
}


type httpClient struct{
	config *config.GeneralConfig
}

func NewHttpClient(config *config.GeneralConfig, newRelicApp newrelic.Application, enableNewrelic bool) HttpClient{
	return &httpClient{
		config: config,
	}
}

func (httpClient *httpClient) JSONPost(ctx context.Context, payload interface{}, headers interface{}, url string, newrelicMetricName string) (*http.Response, error) {

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)


	resp, err := httpClient.sendHttpRequest(req, newrelicMetricName)
	if err != nil {
		return nil, err
	}

	return resp, err
}


func (httpClient *httpClient) sendHttpRequest(req *http.Request, name string) (*http.Response, error) {
	client := http.Client{Transport: &nethttp.Transport{}}
	client.Timeout = time.Duration(10 * time.Second)
	return client.Do(req)
}

