package ecodes

const(
	InternalServerError uint32 = 50001

	DatabaseConnnectionFailed uint32 = 50002

	FailedToDecodeRequestBody uint32 = 40002

	ValidateField   uint32 = 40001

	ValidationUnknown  uint32 = 40004

	NotFound   uint32 = 40401

	StaleDataErrorCode = 403001

	// Users
	DuplicateUser = 40004

	UserNotFound   uint32 = 40402

	InvalidRequestParameters uint32 = 40005
)
