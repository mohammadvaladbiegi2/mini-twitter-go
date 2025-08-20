package userprofile

import (
	userdtos "twitter_clone/internal/modules/user/dto"
	"twitter_clone/internal/pkg/apperror"
	"twitter_clone/internal/pkg/validation"
)

type Service interface {
	UpdateProfile(userID int64, updateProfileReq userdtos.UpdateProfileReq) (userdtos.UpdateProfileRes, *apperror.AppError)
	GetProfile(userID int64) (userdtos.UserGetProfileRes, *apperror.AppError)
}

type ProfileService struct {
	repo Repository
}

func NewProfileService(repo Repository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) UpdateProfile(userID int64, req userdtos.UpdateProfileReq) (userdtos.UpdateProfileRes, *apperror.AppError) {
	if validErrors := validation.ValidateUpdateProfileReq(req); validErrors != nil {
		return userdtos.UpdateProfileRes{}, validErrors
	}
	return s.repo.UpdateProfile(userID, req)
}

func (s *ProfileService) GetProfile(userID int64) (userdtos.UserGetProfileRes, *apperror.AppError) {
	return s.repo.GetProfile(userID)
}
