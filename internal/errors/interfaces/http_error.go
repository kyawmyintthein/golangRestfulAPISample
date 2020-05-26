package interfaces

type HttpError interface {
	StatusCode() int
}
