package api_errors

import (
	error_const "github.com/kyawmyintthein/golangRestfulAPISample/app/constant/error-const"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/errors"
)

type DuplicateResourceError struct{
	*errors.BaseError
}

func NewDuplicateResourceError(resource string) *DuplicateResourceError{
	baseErr := errors.NewError(error_const.DuplicateResource, error_const.DuplicateResourceError, error_const.ErrorMapping[error_const.DuplicateResourceError], "resource", resource)
	return &DuplicateResourceError{
		baseErr,
	}
}

func (e *DuplicateResourceError) Wrap(err error) error {
	e.BaseError.Wrap(err)
	return e
}

