package tweetmedia

import (
	"mime/multipart"
	"twitter_clone/internal/modules/upload"
	uploaddto "twitter_clone/internal/modules/upload/dto"
	"twitter_clone/internal/pkg/apperror"
)

type Repository interface {
	UpdateTweetImage(tweetID, uploadID int64, imageURL string) error
}

type Service struct {
	uploadSvc  upload.Service
	tweetRepo  Repository
	uploadRepo upload.Repository
	storage    upload.Storage
}

func NewService(uSvc upload.Service, tRepo Repository, upRepo upload.Repository, st upload.Storage) *Service {
	return &Service{
		uploadSvc:  uSvc,
		tweetRepo:  tRepo,
		uploadRepo: upRepo,
		storage:    st,
	}
}

func (s *Service) UploadTweetImage(userID, tweetID int64, fh *multipart.FileHeader) (*uploaddto.UploadResponse, *apperror.AppError) {
	uploadRes, appErr := s.uploadSvc.UploadFile(userID, fh)
	if appErr != nil {
		return nil, appErr
	}

	if err := s.tweetRepo.UpdateTweetImage(tweetID, uploadRes.ID, uploadRes.URL); err != nil {
		return nil, apperror.DB("failed to update tweet image", err)
	}

	return uploadRes, nil
}
