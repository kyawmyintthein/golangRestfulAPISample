package interfaces

type ErrorFormatter interface {
	GetArgs() []interface{}
	GetMessage() string
	FormattedMessage() string
}
