package tweetreply

import (
	"strings"

	tweetdtos "twitter_clone/internal/modules/tweet/dto"
	"twitter_clone/internal/pkg/apperror"
)

type Service interface {
	CreateReply(userID, tweetID int64, req tweetdtos.CreateReplyReq) (tweetdtos.CreateReplyRes, *apperror.AppError)
	GetReplies(tweetID int64, limit int, cursorID *int64) (tweetdtos.GetRepliesRes, *apperror.AppError)
}

type ReplyService struct {
	Repo Repository
}

func NewReplyService(repo Repository) *ReplyService {
	return &ReplyService{Repo: repo}
}

func (s *ReplyService) CreateReply(userID, tweetID int64, req tweetdtos.CreateReplyReq) (tweetdtos.CreateReplyRes, *apperror.AppError) {
	// Validation ساده (تو قبلاً پکیج validation داری؛ اگر خواستی اونجا منتقلش کن)
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return tweetdtos.CreateReplyRes{}, apperror.Validation("Validation failed", []map[string]string{
			{"content": "content is required"},
		}, nil)
	}
	if len([]rune(content)) > 280 {
		return tweetdtos.CreateReplyRes{}, apperror.Validation("Validation failed", []map[string]string{
			{"content": "content must be <= 280 chars"},
		}, nil)
	}

	rp, err := s.Repo.CreateReply(userID, tweetID, content)
	if err != nil {
		return tweetdtos.CreateReplyRes{}, err
	}
	return tweetdtos.CreateReplyRes{Reply: rp}, nil
}

func (s *ReplyService) GetReplies(tweetID int64, limit int, cursorID *int64) (tweetdtos.GetRepliesRes, *apperror.AppError) {
	return s.Repo.GetReplies(tweetID, limit, cursorID)
}
