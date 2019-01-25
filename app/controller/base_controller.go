package controller

import (
	"encoding/json"
	"fmt"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/ecodes"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/model"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/errors"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"io/ioutil"
	"net/http"
	"reflect"
)

type BaseController struct{
	Logging  logging.Logger
	Config   *config.GeneralConfig
}

func (c *BaseController) WriteWithStatus(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func (c *BaseController) WriteJSON(r *http.Request, w http.ResponseWriter, statusCode int, v interface{}) error {
	log := c.Logging.GetLogger(r.Context())
	w.WriteHeader(statusCode)
	resp :=  model.SuccessResponse{
		Success:true,
		Data: v,
	}

	if c.Config.Log.LogLevel == "debug"{
		log.WithField(constant.Identifier, "success_response").Debug(fmt.Sprintf("%+v", resp))
	}

	return c.writeJSON(w, resp)
}

// argument error object need to be CustomError type.
// if not, this function with return with 500 status code as default.
func (c *BaseController) WriteError(r *http.Request, w http.ResponseWriter, err error) error {
	log := c.Logging.GetLogger(r.Context())
	code := getErrorCode(err)
	statusCode := getHttpStatus(code)
	w.WriteHeader(statusCode)
	log.WithError(err).Errorln("failed_response")
	return c.writeJSON(w, &model.ErrorResponse{
		Success: false,
		Error: model.HttpError{
			Code:    code,
			Message: err.Error(),
		},
	})
}

func (c *BaseController) WriteErrorWithMessage(r *http.Request, w http.ResponseWriter, err error, message string) error {
	log := c.Logging.GetLogger(r.Context())
	code := getErrorCode(err)
	statusCode := getHttpStatus(code)

	w.WriteHeader(statusCode)
	log.WithError(err).Errorln("failed_response")
	return c.writeJSON(w, &model.ErrorResponse{
		Success: false,
		Error: model.HttpError{
			Code:    code,
			Message: message,
		},
	})
}

func (c *BaseController) writeJSON(w http.ResponseWriter, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}


// error is not CustomError type return default error code (Internal Server Error)
func getErrorCode(err error) uint32 {
	clError, ok := err.(errors.CustomError)
	if !ok {
		return ecodes.InternalServerError
	}
	return clError.GetCode()
}

func getHttpStatus(code uint32) (status int) {
	firstThreeDigits := code / 100
	switch firstThreeDigits {
	case 400:
		status = http.StatusBadRequest
	case 401:
		status = http.StatusUnauthorized
	case 403:
		status = http.StatusForbidden
	case 404:
		status = http.StatusNotFound
	case 405:
		status = http.StatusMethodNotAllowed
	case 406:
		status = http.StatusNotAcceptable
	case 408:
		status = http.StatusRequestTimeout
	default:
		status = http.StatusInternalServerError
	}
	return
}



func (c *BaseController) decodeAndValidate(r *http.Request, v model.RequestValidator) error {
	err := c.decodeRequestBody(r, v)
	if err != nil {
		return err
	}
	return v.Validate(r.Context())
}

func (c *BaseController) decodeRequestBody(r *http.Request, v interface{}) (err error) {
	log := c.Logging.GetLogger(r.Context())

	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		err := errors.New(ecodes.InternalServerError, constant.ServerIssue)
		log.WithError(err).Errorf("unsupported type to decode %+v ", v)
		return err
	}

	payloadBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if c.Config.Log.LogLevel == "debug"{
		log.WithField(constant.Identifier, "payload").Debug(string(payloadBytes))
	}

	//err = json.NewDecoder(payloadBytes).Decode(v)
	json.Unmarshal(payloadBytes, &v)
	if err != nil {
		err = errors.Wrap(err, ecodes.FailedToDecodeRequestBody, constant.DecodeRequestBodyErr)
		return err
	}

	return nil
}
