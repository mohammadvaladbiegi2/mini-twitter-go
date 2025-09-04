package timeline

import (
	"twitter_clone/internal/pkg/apperror"
)

type TimeLineService struct {
	repo Repository
}

type Service interface {
	AddDataToRedis(userID int64) (R, *apperror.AppError)
}

func NewTimeLineService(repo Repository) *TimeLineService {
	return &TimeLineService{repo}
}

type R struct {
	ID           int64   `json:"id"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	PasswordHash string  `json:"-"`
	Bio          *string `json:"bio"`
	AvatarURL    *string `json:"avatar_url"`
}

func (t TimeLineService) AddDataToRedis(userID int64) (R, *apperror.AppError) {

	user, err := t.repo.AddDataToRedis(userID)
	if err != nil {
		return R{}, err
	}

	return user, nil
}
