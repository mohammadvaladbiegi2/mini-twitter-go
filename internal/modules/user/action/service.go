package useraction

import "twitter_clone/internal/pkg/apperror"

type Service interface {
	Follow(followerID, followeeID int64) *apperror.AppError
	Unfollow(followerID, followeeID int64) *apperror.AppError
}

type UserActionService struct {
	Repo Repository
}

func NewUserActionService(repo Repository) *UserActionService {
	return &UserActionService{Repo: repo}
}

func (s *UserActionService) Follow(followerID, followeeID int64) *apperror.AppError {
	if followerID == followeeID {
		return apperror.Validation("cannot follow yourself", nil, nil)
	}
	return s.Repo.Follow(followerID, followeeID)
}

func (s *UserActionService) Unfollow(followerID, followeeID int64) *apperror.AppError {
	if followerID == followeeID {
		return apperror.Validation("cannot unfollow yourself", nil, nil)
	}
	return s.Repo.Unfollow(followerID, followeeID)
}
