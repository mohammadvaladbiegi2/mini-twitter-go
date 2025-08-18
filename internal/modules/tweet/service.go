package tweet

import (
	tweetdtos "twitter_clone/internal/modules/tweet/dto"
	"twitter_clone/internal/pkg/apperror"
	"twitter_clone/internal/pkg/validation"
)

type Service interface {
	CreateTweet(UserID int64, req tweetdtos.CreateTweetReq) (tweetdtos.CreateTweetRes, *apperror.AppError)
}

type TweetService struct {
	Repo Repository
}

func NewTweetService(repo Repository) *TweetService {
	return &TweetService{Repo: repo}
}

func (t TweetService) CreateTweet(UserID int64, req tweetdtos.CreateTweetReq) (tweetdtos.CreateTweetRes, *apperror.AppError) {

	if validationError := validation.ValidateCreateTweetReq(req); validationError != nil {
		return tweetdtos.CreateTweetRes{}, validationError
	}

	tweet, RepoError := t.Repo.CreateTweet(UserID, req)
	if RepoError != nil {
		return tweetdtos.CreateTweetRes{}, RepoError
	}

	return tweet, nil
}
