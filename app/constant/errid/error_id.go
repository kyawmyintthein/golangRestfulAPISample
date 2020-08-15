package errid

const (
	UnknownError               string = "UnknownError"
	InvalidRequestPayloadError string = "InvalidRequestPayload"
	DuplicateResourceError     string = "DuplicateResource"
)

var ErrorMapping map[string]string

func init() {
	ErrorMapping = map[string]string{
		UnknownError:           "Server Issues",
		DuplicateResourceError: "{{var_resource}} is already exist",
	}
}
