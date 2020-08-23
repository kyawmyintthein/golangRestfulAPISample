package errors

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/errcode"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/errid"
	"github.com/kyawmyintthein/orange-contrib/errorx"
)

type FailedToDecodeRequestBodyError struct {
	*errorx.ErrorX
	*errorx.ErrorWithCode
	*errorx.ErrorWithID
	*errorx.ErrorWithHttpStatus
	*errorx.ErrorStacktrace
}

func NewFailedToDecodeRequestBodyError() *FailedToDecodeRequestBodyError {
	return &FailedToDecodeRequestBodyError{
		errorx.NewErrorX(errid.ErrorMapping[errid.InvalidRequestPayloadError]),
		errorx.NewErrorWithCode(errcode.InvalidRequestPayload),
		errorx.NewErrorWithID(errid.InvalidRequestPayloadError),
		errorx.NewErrorWithHttpStatus(errorx.GenerateHttpStatusCodeFromErrorCode(errcode.InvalidRequestPayload)),
		errorx.NewErrorWithStackTrace(constant.DefaultErrorStackLen, constant.DefaultErrorCallerLen),
	}
}
