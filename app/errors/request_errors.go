package errors

import (
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
		errorx.NewErrorX(errid.ErrorMapping[errid.InvalidRequestPayload]),
		errorx.NewErrorWithCode(errcode.InvalidRequestPayloadError),
		errorx.NewErrorWithID(errid.InvalidRequestPayload),
		errorx.NewErrorWithHttpStatus(errorx.GenerateHttpStatusCodeFromErrorCode(errcode.InvalidRequestPayloadError)),
		errorx.NewErrorWithStackTrace(constant.DefaultErrorStackLen, contstant.DefaultErrorCallerLen),
	}
}
