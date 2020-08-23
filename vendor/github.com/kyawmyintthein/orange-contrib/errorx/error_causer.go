package errorx

type Causer interface {
	Cause() error
}
