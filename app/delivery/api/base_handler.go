package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/errors"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/infrastructure"
	"github.com/kyawmyintthein/orange-contrib/logx"
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
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.NewFailedToDecodeRequestBodyError().Wrap(err)
	}
	r.Body.Close()

	if err := json.NewDecoder(bytes.NewBuffer(payload)).Decode(v); err != nil {
		return errors.NewFailedToDecodeRequestBodyError().Wrap(err)
	}
	defer r.Body.Close()
	logx.Debugf(r.Context(), "REQUEST PAYLOAD: %s", string(payload))
	return v.Validate(r.Context())
}

func (h *BaseHandler) URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
