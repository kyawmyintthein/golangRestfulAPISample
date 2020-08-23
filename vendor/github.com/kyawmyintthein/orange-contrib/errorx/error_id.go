package errorx

type ErrorID interface {
	ID() string
}

type ErrorWithID struct {
	id string
}

func NewErrorWithID(id string) *ErrorWithID {
	return &ErrorWithID{id: id}
}

func (err *ErrorWithID) ID() string {
	return err.id
}
