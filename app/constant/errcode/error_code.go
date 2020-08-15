package errcode

const (
	InternalServerError int = iota + 500000
)

const (
	InvalidRequestPayload int = iota + 400000
	DuplicateResource     int = iota
)
