package api_errors

import (
	error_const "github.com/kyawmyintthein/golangRestfulAPISample/app/constant/error-const"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/errors"
)

type UnknownError struct{
	*errors.BaseError
}

func NewUnknownError() *UnknownError{
	baseErr := errors.NewError(error_const.InternalServerError, error_const.UnknownError, error_const.ErrorMapping[error_const.UnknownError])
	return &UnknownError{
		baseErr,
	}
}

func (e *UnknownError) Wrap(err error) error {
	e.BaseError.Wrap(err)
	return e
}