package avatar

import (
	"mime/multipart"
	"twitter_clone/internal/modules/upload"
	uploaddto "twitter_clone/internal/modules/upload/dto"
	uploadavataruser "twitter_clone/internal/modules/upload/user"
	"twitter_clone/internal/pkg/apperror"
)

type Service struct {
	uploadSvc  upload.Service
	userRepo   uploadavataruser.Repository
	uploadRepo upload.Repository
	storage    upload.Storage
}

func NewService(uSvc upload.Service, uRepo uploadavataruser.Repository, upRepo upload.Repository, st upload.Storage) *Service {
	return &Service{
		uploadSvc:  uSvc,
		userRepo:   uRepo,
		uploadRepo: upRepo,
		storage:    st,
	}
}

func (s *Service) UploadAvatar(userID int64, fh *multipart.FileHeader) (*uploaddto.UploadResponse, *apperror.AppError) {
	uploadRes, appErr := s.uploadSvc.UploadFile(userID, fh)
	if appErr != nil {
		return nil, appErr
	}

	oldID, err := s.userRepo.GetAvatarUploadID(userID)
	if err != nil {
		return nil, apperror.DB("failed to get old avatar", err)
	}

	if err := s.userRepo.UpdateAvatar(userID, uploadRes.ID, uploadRes.URL); err != nil {
		return nil, apperror.DB("failed to update avatar", err)
	}

	if oldID > 0 && oldID != uploadRes.ID {
		if old, err := s.uploadRepo.FindByID(oldID); err == nil {
			_ = s.storage.Delete(old.FileName)
			_ = s.uploadRepo.DeleteByID(oldID)
		}
	}

	return uploadRes, nil
}
