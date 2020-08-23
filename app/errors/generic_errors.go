package errors

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/errcode"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/errid"
	"github.com/kyawmyintthein/orange-contrib/errorx"
)

type UnknownError struct {
	*errorx.ErrorX
	*errorx.ErrorWithCode
	*errorx.ErrorWithID
	*errorx.ErrorWithHttpStatus
	*errorx.ErrorStacktrace
}

func NewUnknownError() *UnknownError {
	return &UnknownError{
		errorx.NewErrorX(errid.ErrorMapping[errid.UnknownError]),
		errorx.NewErrorWithCode(errcode.InternalServerError),
		errorx.NewErrorWithID(errid.UnknownError),
		errorx.NewErrorWithHttpStatus(errorx.GenerateHttpStatusCodeFromErrorCode(errcode.InternalServerError)),
		errorx.NewErrorWithStackTrace(constant.DefaultErrorStackLen, constant.DefaultErrorCallerLen),
	}
}
