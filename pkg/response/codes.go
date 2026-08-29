package response

const (
	//	Auth
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
	ErrCodeInvalidToken       = "INVALID_TOKEN"
	ErrCodeTokenExpired       = "TOKEN_EXPIRED"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"

	//	Validation
	ErrCodeBadRequest = "BAD_REQUEST"
	ErrCodeValidation = "VALIDATION_ERROR"

	// Resource
	ErrCodeNotFound = "NOT_FOUND"
	ErrCodeConflict = "CONFLICT"

	//	Server
	ErrCodeInternal = "INTERNAL_ERROR"
)
