package clerrors

import (
	"bytes"
)

func New(id string, code int, messageFormat string, args ...interface{}) error {
	return WithMessage(code, id, messageFormat, args...)
}

func Wrap(err error, id string, code int, messageFormat string, args ...interface{}) error {
	return WithMessage(code, id, messageFormat, args...).Wrap(err)
}

func GetErrorMessages(e error) string {
	return extractFullErrorMessage(e, false)
}

func GetErrorMessagesWithStack(e error) string {
	return extractFullErrorMessage(e, true)
}

func extractFullErrorMessage(e error, includeStack bool) string {
	type causer interface {
		Cause() error
	}

	var ok bool
	var lastClErr error
	errMsg := bytes.NewBuffer(make([]byte, 0, 1024))
	dbxErr := e
	for {
		_, ok := dbxErr.(StackTracer)
		if ok {
			lastClErr = dbxErr
		}

		errorWithFormat, ok := dbxErr.(ErrorFormatter)
		if ok {
			errMsg.WriteString(errorWithFormat.FormattedMessage())
		}

		errorCauser, ok := dbxErr.(causer)
		if ok {
			innerErr := errorCauser.Cause()
			if innerErr == nil {
				break
			}
			dbxErr = innerErr
		} else {
			// We have reached the end and traveresed all inner clerrors.
			// Add last message and exit loop.
			errMsg.WriteString(dbxErr.Error())
			break
		}
		errMsg.WriteString(", ")
	}

	stackError, ok := lastClErr.(StackTracer)
	if includeStack && ok {
		errMsg.WriteString("\nSTACK TRACE:\n")
		errMsg.WriteString(stackError.GetStack())
	}
	return errMsg.String()
}

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
