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

func NewDuplicateResourceError() *DuplicateResourceError {
	return &DuplicateResourceError{
		errorx.NewErrorX(errid.ErrorMapping[errid.DuplicateResource]),
		errorx.NewErrorWithCode(errcode.DuplicateResourceError),
		errorx.NewErrorWithID(errid.DuplicateResource),
		errorx.NewErrorWithHttpStatus(errorx.GenerateHttpStatusCodeFromErrorCode(errcode.DuplicateResourceError)),
		errorx.NewErrorWithStackTrace(constant.DefaultErrorStackLen, contstant.DefaultErrorCallerLen),
	}
}
