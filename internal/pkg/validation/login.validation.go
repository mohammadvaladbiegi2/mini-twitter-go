package validation

import (
	"twitter_clone/internal/modules/auth/dtos"
	"twitter_clone/internal/pkg/apperror"
)

func ValidateLoginReq(req dtos.LoginReq) *apperror.AppError {
	var validationErrors []map[string]string

	if req.UserName == "" {
		validationErrors = append(validationErrors, map[string]string{"username": "Username is required"})
	} else if len(req.UserName) < 3 {
		validationErrors = append(validationErrors, map[string]string{"username": "Username must be at least 3 characters"})
	}

	if req.Password == "" {
		validationErrors = append(validationErrors, map[string]string{"password": "Password is required"})
	} else if len(req.Password) < 6 {
		validationErrors = append(validationErrors, map[string]string{"password": "Password must be at least 6 characters"})
	}

	if len(validationErrors) > 0 {
		return apperror.Validation("Validation failed", validationErrors, nil)
	}

	return nil
}
