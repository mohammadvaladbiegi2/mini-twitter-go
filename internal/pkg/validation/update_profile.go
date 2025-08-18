package validation

import (
	userdtos "twitter_clone/internal/modules/user/dto"
	"twitter_clone/internal/pkg/apperror"
)

func ValidateUpdateProfileReq(req userdtos.UpdateProfileReq) *apperror.AppError {
	var validationErrors []map[string]string

	if req.Username == "" && req.Bio == nil {
		validationErrors = append(validationErrors, map[string]string{
			"error": "At least one field (username or bio) must be provided",
		})
		return apperror.Validation("Validation failed", validationErrors, nil)
	}

	if req.Username != "" {
		if len(req.Username) < 3 {
			validationErrors = append(validationErrors, map[string]string{"username": "Username must be at least 3 characters"})
		} else if len(req.Username) > 30 {
			validationErrors = append(validationErrors, map[string]string{"username": "Username must not exceed 30 characters"})
		}
	}

	if req.Bio != nil && len(*req.Bio) > 160 {
		validationErrors = append(validationErrors, map[string]string{"bio": "Bio must not exceed 160 characters"})
	}

	if len(validationErrors) > 0 {
		return apperror.Validation("Validation failed", validationErrors, nil)
	}

	return nil
}
