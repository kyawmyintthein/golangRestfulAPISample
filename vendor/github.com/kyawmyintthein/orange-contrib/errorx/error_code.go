package errorx

type ErrorCode interface {
	Code() int
}

type ErrorWithCode struct {
	code int
}

func NewErrorWithCode(code int) *ErrorWithCode {
	return &ErrorWithCode{code: code}
}

func (err *ErrorWithCode) Code() int {
	return err.code
}
