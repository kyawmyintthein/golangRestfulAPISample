package errorx

import "strconv"

const (
	DefaultHttpStatusCode int = 500 // Internal Server Error
)

type HttpError interface {
	StatusCode() int
}

type ErrorWithHttpStatus struct {
	httpStatus int
}

func NewErrorWithHttpStatus(httpStatus int) *ErrorWithHttpStatus {
	return &ErrorWithHttpStatus{httpStatus: httpStatus}
}

func (err *ErrorWithHttpStatus) StatusCode() int {
	return err.httpStatus
}

func GenerateHttpStatusCodeFromErrorCode(code int) int {
	arr := strconv.Itoa(code)
	if len(arr) >= 0 {
		f3d := arr[:3]
		statusCode, err := strconv.Atoi(f3d)
		if err != nil {
			return DefaultHttpStatusCode
		}
		return statusCode
	}
	return DefaultHttpStatusCode

}
