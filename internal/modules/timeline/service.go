package timeline

import (
	timelinedto "twitter_clone/internal/modules/timeline/dto"
	"twitter_clone/internal/pkg/apperror"
)

type TimeLineService struct {
	repo Repository
}

type Service interface {
	MyTimeline(userID int64, limit int, cursor *string) (*timelinedto.TimelineResponse, *apperror.AppError)
}

func NewTimeLineService(repo Repository) *TimeLineService {
	return &TimeLineService{repo}
}

func (t TimeLineService) MyTimeline(userID int64, limit int, cursor *string) (*timelinedto.TimelineResponse, *apperror.AppError) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return t.repo.MyTimeline(userID, limit, cursor)
}
