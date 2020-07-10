package api_errors

import (
	error_const "github.com/kyawmyintthein/golangRestfulAPISample/app/constant/error-const"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/errors"
)

type FailedToDecodeRequestBodyError struct{
	*errors.BaseError
}

func NewFailedToDecodeRequestBodyError() *FailedToDecodeRequestBodyError{
	baseErr := errors.NewError(error_const.InvalidRequestPayload, error_const.InvalidRequestPayloadError, error_const.ErrorMapping[error_const.InvalidRequestPayloadError])
	return &FailedToDecodeRequestBodyError{
		baseErr,
	}
}

func (e *FailedToDecodeRequestBodyError) Wrap(err error) error {
	e.BaseError.Wrap(err)
	return e
}
