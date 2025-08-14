package validation

import (
	"regexp"
	"twitter_clone/internal/modules/auth/dtos"
	"twitter_clone/internal/pkg/apperror"
)

func ValidateSignUpReq(req dtos.SignUpReq) *apperror.AppError {
	var validationErrors []map[string]string

	if req.Username == "" {
		validationErrors = append(validationErrors, map[string]string{"username": "Username is required"})
	} else if len(req.Username) < 3 {
		validationErrors = append(validationErrors, map[string]string{"username": "Username must be at least 3 characters"})
	}

	if req.Email == "" {
		validationErrors = append(validationErrors, map[string]string{"email": "Email is required"})
	} else if !isValidEmail(req.Email) {
		validationErrors = append(validationErrors, map[string]string{"email": "Invalid email format"})
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

func isValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}
