package clerrors

type HttpError interface {
	StatusCode() int
}
