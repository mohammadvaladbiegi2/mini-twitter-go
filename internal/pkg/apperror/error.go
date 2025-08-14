package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorType string

type AppError struct {
	Type       ErrorType           `json:"type"`
	Message    string              `json:"message"`
	Errors     []map[string]string `json:"errors,omitempty"`
	StatusCode int                 `json:"status_code"`
	Err        error               `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(errType ErrorType, msg string, details []map[string]string, status int, err error) *AppError {
	return &AppError{
		Type:       errType,
		Message:    msg,
		Errors:     details,
		StatusCode: status,
		Err:        err,
	}
}

func Validation(msg string, details []map[string]string, err error) *AppError {
	return New("VALIDATION_ERROR", msg, details, http.StatusBadRequest, err)
}

func DB(msg string, err error) *AppError {
	return New("DATABASE_ERROR", msg, nil, http.StatusInternalServerError, err)
}

func Server(msg string, err error) *AppError {
	return New("SERVER_ERROR", msg, nil, http.StatusInternalServerError, err)
}

func NotFound(msg string, err error) *AppError {
	return New("NOT_FOUND", msg, nil, http.StatusNotFound, err)
}

func UnauthorizedErr(msg string, err error) *AppError {
	return New("UNAUTHORIZED", msg, nil, http.StatusUnauthorized, err)
}

func Forbidden(msg string, err error) *AppError {
	return New("FORBIDDEN", msg, nil, http.StatusForbidden, err)
}

func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}
