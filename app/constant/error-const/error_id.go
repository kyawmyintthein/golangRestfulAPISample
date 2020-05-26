package error_const

const(
	UnknownError string = "UnknownError"
)

var ErrorMapping map[string]string

func init(){
	ErrorMapping = map[string]string{
		UnknownError: "Server Issues",
	}
}
