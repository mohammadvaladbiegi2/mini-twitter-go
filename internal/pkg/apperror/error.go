package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorType فقط یک رشته است، محدود نیست
type ErrorType string

// AppError ساختار اصلی خطا
type AppError struct {
	Type       ErrorType           `json:"type"`             // نوع خطا (قابل سفارشی‌سازی)
	Message    string              `json:"message"`          // پیام کلی
	Errors     []map[string]string `json:"errors,omitempty"` // می‌تونه slice یا map از خطاهای جزئی باشه
	StatusCode int                 `json:"status_code"`      // کد HTTP
	Err        error               `json:"-"`                // خطای اصلی (برای لاگ داخلی)
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
		Errors:     details, // می‌تونه nil باشه
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

func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}
