package errors

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/errcode"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/errid"
	"github.com/kyawmyintthein/orange-contrib/errorx"
)

type DuplicateResourceError struct {
	*errorx.ErrorX
	*errorx.ErrorWithCode
	*errorx.ErrorWithID
	*errorx.ErrorWithHttpStatus
	*errorx.ErrorStacktrace
}

func NewDuplicateResourceError(resource string) *DuplicateResourceError {
	return &DuplicateResourceError{
		errorx.NewErrorX(errid.ErrorMapping[errid.DuplicateResourceError], "resource", resource),
		errorx.NewErrorWithCode(errcode.DuplicateResource),
		errorx.NewErrorWithID(errid.DuplicateResourceError),
		errorx.NewErrorWithHttpStatus(errorx.GenerateHttpStatusCodeFromErrorCode(errcode.DuplicateResource)),
		errorx.NewErrorWithStackTrace(constant.DefaultErrorStackLen, constant.DefaultErrorCallerLen),
	}
}
