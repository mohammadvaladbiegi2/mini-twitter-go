package user

import (
	userdtos "twitter_clone/internal/modules/user/dto"
	"twitter_clone/internal/pkg/apperror"
	"twitter_clone/internal/pkg/validation"
)

type Service interface {
	UpdateProfile(userID int64, updateProfileReq userdtos.UpdateProfileReq) (userdtos.UpdateProfileRes, *apperror.AppError)
	GetProfile(userID int64) (userdtos.UserGetProfileRes, *apperror.AppError)
}

type UserService struct {
	Repo Repository
}

func NewUserService(repo Repository) *UserService {
	return &UserService{Repo: repo}
}

func (u UserService) UpdateProfile(userID int64, updateProfileReq userdtos.UpdateProfileReq) (userdtos.UpdateProfileRes, *apperror.AppError) {

	if validErrors := validation.ValidateUpdateProfileReq(updateProfileReq); validErrors != nil {
		return userdtos.UpdateProfileRes{}, validErrors
	}

	userupdate, err := u.Repo.UpdateProfile(userID, updateProfileReq)
	if err != nil {
		return userdtos.UpdateProfileRes{}, err
	}

	return userupdate, nil
}
func (u UserService) GetProfile(userID int64) (userdtos.UserGetProfileRes, *apperror.AppError) {

	userProfile, err := u.Repo.GetProfile(userID)
	if err != nil {
		return userdtos.UserGetProfileRes{}, err
	}

	return userProfile, nil
}
