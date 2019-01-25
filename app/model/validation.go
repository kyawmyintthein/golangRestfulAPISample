package model

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/ecodes"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/errors"
	"golang.org/x/net/context"
	"gopkg.in/validator.v2"
)

type RequestValidator interface {
	Validate(ctx context.Context) error
}

// validateFields checks if the required fields in a model is filled.
func ValidateFields(model interface{}) error {
	err := validator.Validate(model)
	if err != nil {
		errs, ok := err.(validator.ErrorMap)
		if ok {
			for f, _ := range errs {
				return errors.New(ecodes.ValidateField, constant.ValidateFieldErr+"-"+f)
			}
		} else {
			return errors.New(ecodes.ValidationUnknown, constant.ValidationUnknownErr)
		}
	}

	return nil
}
