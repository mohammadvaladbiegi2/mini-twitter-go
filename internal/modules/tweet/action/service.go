package tweetaction

import "twitter_clone/internal/pkg/apperror"

type Service interface {
	Like(userID, tweetID int64) *apperror.AppError
	Dislike(userID, tweetID int64) *apperror.AppError
	RemoveReaction(userID, tweetID int64) *apperror.AppError
}

type TweetActionService struct {
	Repo Repository
}

func NewTweetActionService(repo Repository) *TweetActionService {
	return &TweetActionService{Repo: repo}
}

func (s *TweetActionService) Like(userID, tweetID int64) *apperror.AppError {
	return s.Repo.React(userID, tweetID, true)
}

func (s *TweetActionService) Dislike(userID, tweetID int64) *apperror.AppError {
	return s.Repo.React(userID, tweetID, false)
}

func (s *TweetActionService) RemoveReaction(userID, tweetID int64) *apperror.AppError {
	return s.Repo.RemoveReaction(userID, tweetID)
}
