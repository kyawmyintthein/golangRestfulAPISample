package api

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/api_errors"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/infrastructure"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"io/ioutil"
	"net/http"
)

type BaseHandler struct {
	infrastructure.HttpResponseWriter
	Config *config.GeneralConfig
}

func ProvideBaseHandler(config *config.GeneralConfig, httpResponseWriter infrastructure.HttpResponseWriter) *BaseHandler {
	return &BaseHandler{
		httpResponseWriter,
		config,
	}
}

func (h *BaseHandler) DecodeAndValidate(r *http.Request, v infrastructure.RequestValidator) error {
	log := logging.GetStructuredLogger(r.Context())

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return api_errors.NewFailedToDecodeRequestBodyError().Wrap(err)
	}
	r.Body.Close()

	if err := json.NewDecoder(bytes.NewBuffer(payload)).Decode(v); err != nil {
		return api_errors.NewFailedToDecodeRequestBodyError().Wrap(err)
	}
	defer r.Body.Close()
	log.Debugf("REQUEST PAYLOAD: %s", string(payload))
	return v.Validate(r.Context())
}


func (h *BaseHandler) URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
