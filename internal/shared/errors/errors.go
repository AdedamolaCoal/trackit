package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	HttpStatus int
	Code       string
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func Unauthorized(message string) *AppError {
	return &AppError{
		HttpStatus: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    message,
	}
}

func NotFound(resource string) *AppError {
	return &AppError{
		HttpStatus: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf("%s not found", resource),
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		HttpStatus: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    message,
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		HttpStatus: http.StatusConflict,
		Code:       "CONFLICT",
		Message:    message,
	}
}

func Internal(message string, err error) *AppError {
	return &AppError{
		HttpStatus: http.StatusInternalServerError,
		Code:       "INTERNAL",
		Message:    message,
		Err:        err,
	}
}

func BadRequest(message string) *AppError {
	return &AppError{
		HttpStatus: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    message,
	}
}

func InvalidToken(message string) *AppError {
	return &AppError{
		HttpStatus: http.StatusUnauthorized,
		Code:       "INVALID_TOKEN",
		Message:    message,
	}
}

func TokenExpired() *AppError {
	return &AppError{
		HttpStatus: http.StatusUnauthorized,
		Code:       "TOKEN_EXPIRED",
		Message:    "Token has expired",
	}
}

func InvalidCredentials() *AppError {
	return &AppError{
		HttpStatus: http.StatusUnauthorized,
		Code:       "INVALID_CREDENTIALS",
		Message:    "Invalid username or password",
	}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, &target)
}
