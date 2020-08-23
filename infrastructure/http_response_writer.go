package infrastructure

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/errcode"
	"github.com/kyawmyintthein/orange-contrib/errorx"
	"github.com/kyawmyintthein/orange-contrib/logx"
)

type HttpResponseWriter interface {
	RenderJSON(*http.Request, http.ResponseWriter, int, interface{}) error
	Status(w http.ResponseWriter, statusCode int)
	RenderErrorAsJSON(r *http.Request, w http.ResponseWriter, err error, messages ...string) error
	RenderPlainText(r *http.Request, w http.ResponseWriter, statusCode int, v interface{}) error
}

type ResponseFormat struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
}

type ErrorResponse struct {
	Code        int      `json:"code"`
	Message     string   `json:"message"`
	Description string   `json:"description"`
	Errors      []string `json:"errors"`
}

type httpResponseWriter struct {
}

func NewHttpResponseWriter() HttpResponseWriter {
	return &httpResponseWriter{}
}

func (c *httpResponseWriter) RenderPlainText(
	r *http.Request,
	w http.ResponseWriter,
	statusCode int,
	v interface{},
) error {
	//	log := logging.Logger.GetLogger(ctx)
	_, err := fmt.Fprintf(w, v.(string))
	if err != nil {
		_, _ = fmt.Fprintf(w, err.Error())
		w.WriteHeader(500)
		logx.Info(r.Context(), "Response: ", err.Error(), "; StatusCode: ", 500)
		return err
	}
	logx.Info(r.Context(), "Response: ", v.(string), "; StatusCode: ", statusCode)
	w.WriteHeader(statusCode)
	return nil
}

func (c *httpResponseWriter) RenderJSON(r *http.Request, w http.ResponseWriter, statusCode int, v interface{}) error {
	w.WriteHeader(statusCode)
	return c.writeJSON(r, w, v)
}

func (c *httpResponseWriter) writeJSON(r *http.Request, w http.ResponseWriter, v interface{}) error {
	data, err := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json")
	logx.Info(r.Context(), "Response: ", string(data))

	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (c *httpResponseWriter) Status(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

// argument error object need to be CustomError type.
// if not, this function with return with 500 status code as default.
func (c *httpResponseWriter) RenderErrorAsJSON(r *http.Request, w http.ResponseWriter, err error, messages ...string) error {
	code := getErrorCode(err)
	statusCode := getHttpStatus(code)
	var resp interface{}

	statusCode = http.StatusOK
	desc, _ := getTitleAndDescription(messages)
	resp = ResponseFormat{
		Success: false,
		Result: ErrorResponse{
			Code:        code,
			Message:     err.Error(),
			Description: desc,
		},
	}

	w.WriteHeader(statusCode)
	// log.WithError(err).Errorf("HTTP_ERROR_RESPONSE::[%d]", statusCode)
	return c.writeJSON(r, w, resp)

}

// error is not CustomError type return default error code (Internal Server Error)
func getErrorCode(err error) int {
	errWithCode, ok := err.(errorx.ErrorCode)
	if !ok {
		return errcode.InternalServerError
	}
	return errWithCode.Code()
}

func getHttpStatus(code int) (status int) {
	firstThreeDigits := code / 10000
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

func getTitleAndDescription(messages []string) (string, string) {
	var ttl, desc string
	if len(messages) > 0 {
		ttl = messages[0]
	}
	if len(messages) > 1 {
		desc = messages[1]
	}
	return ttl, desc
}
