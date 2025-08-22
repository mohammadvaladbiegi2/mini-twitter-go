package tweetbookmark

import (
	tweetdtos "twitter_clone/internal/modules/tweet/dto"
	"twitter_clone/internal/pkg/apperror"
)

type Service interface {
	Add(userID, tweetID int64) *apperror.AppError
	Remove(userID, tweetID int64) *apperror.AppError
	List(userID int64, limit int, cursor string) (tweetdtos.BookmarksListRes, *apperror.AppError)
}

type BookmarkService struct {
	repo Repository
}

func NewBookMarkService(repo Repository) *BookmarkService {
	return &BookmarkService{repo: repo}
}

func (s *BookmarkService) Add(userID, tweetID int64) *apperror.AppError {
	return s.repo.AddBookmark(userID, tweetID)
}

func (s *BookmarkService) Remove(userID, tweetID int64) *apperror.AppError {
	return s.repo.RemoveBookmark(userID, tweetID)
}

func (s *BookmarkService) List(userID int64, limit int, cursor string) (tweetdtos.BookmarksListRes, *apperror.AppError) {
	afterT, afterID, err := DecodeCursor(cursor)
	if err != nil {
		return tweetdtos.BookmarksListRes{}, apperror.Validation("invalid cursor", nil, err)
	}

	items, lastT, lastID, appErr := s.repo.ListBookmarks(userID, limit, afterT, afterID)
	if appErr != nil {
		return tweetdtos.BookmarksListRes{}, appErr
	}

	var next *string
	if lastT != nil && lastID != nil {
		c := EncodeCursor(*lastT, *lastID)
		next = &c
	}

	return tweetdtos.BookmarksListRes{
		Bookmarks:  items,
		Limit:      limit,
		Count:      len(items),
		NextCursor: next,
	}, nil
}
