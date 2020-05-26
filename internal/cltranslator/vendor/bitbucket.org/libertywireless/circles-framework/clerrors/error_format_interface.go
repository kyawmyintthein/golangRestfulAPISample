package clerrors

type ErrorFormatter interface {
	GetArgs() []interface{}
	GetMessage() string
	FormattedMessage() string
}
