package validation

import (
	tweetdtos "twitter_clone/internal/modules/tweet/dto"
	"twitter_clone/internal/pkg/apperror"
)

func ValidateCreateTweetReq(req tweetdtos.CreateTweetReq) *apperror.AppError {
	var validationErrors []map[string]string

	if req.Content == "" {
		validationErrors = append(validationErrors, map[string]string{"content": "Content is required"})
	} else if len(req.Content) > 280 {
		validationErrors = append(validationErrors, map[string]string{"content": "Content must not exceed 280 characters"})
	}
	if len(req.Tags) > 0 {

		if len(req.Tags) > 10 {
			validationErrors = append(validationErrors, map[string]string{"tags": "You can add up to 10 tags only"})
		}

		for _, t := range req.Tags {
			if len(t) > 30 {
				validationErrors = append(validationErrors, map[string]string{"tags": "Each tag must not exceed 30 characters"})
			}
		}
	}

	if len(validationErrors) > 0 {
		return apperror.Validation("Validation failed", validationErrors, nil)
	}

	return nil
}
