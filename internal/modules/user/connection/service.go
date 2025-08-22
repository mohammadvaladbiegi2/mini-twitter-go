package userconnection

import (
	userdtos "twitter_clone/internal/modules/user/dto"
	"twitter_clone/internal/pkg/apperror"
)

type Service interface {
	GetFollowers(userID int64, limit, offset int) ([]userdtos.UsersFollower, *apperror.AppError)
	GetFollowings(userID int64, limit, offset int) ([]userdtos.UsersFollowing, *apperror.AppError)
}

type UserConnectionService struct {
	Repo Repository
}

func NewUserConnectionService(repo Repository) *UserConnectionService {
	return &UserConnectionService{repo}
}

func (u *UserConnectionService) GetFollowers(userID int64, limit, offset int) ([]userdtos.UsersFollower, *apperror.AppError) {

	// TODO handel offset

	users, err := u.Repo.GetFollowers(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserConnectionService) GetFollowings(userID int64, limit, offset int) ([]userdtos.UsersFollowing, *apperror.AppError) {

	// TODO handel offset

	users, err := u.Repo.GetFollowings(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}
