package userconnection

import (
	userdtos "twitter_clone/internal/modules/user/dto"
	"twitter_clone/internal/pkg/apperror"
)

type PaginatedResponse[T any] struct {
	Data        []T   `json:"data"`
	TotalCount  int64 `json:"total_count"`
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	HasNext     bool  `json:"has_next"`
}

type Service interface {
	GetFollowers(userID int64, page, pageSize int) (*PaginatedResponse[userdtos.UsersFollower], *apperror.AppError)
	GetFollowings(userID int64, page, pageSize int) (*PaginatedResponse[userdtos.UsersFollowing], *apperror.AppError)
}

type UserConnectionService struct {
	Repo Repository
}

func NewUserConnectionService(repo Repository) *UserConnectionService {
	return &UserConnectionService{repo}
}

func (u *UserConnectionService) GetFollowers(userID int64, page, pageSize int) (*PaginatedResponse[userdtos.UsersFollower], *apperror.AppError) {
	users, total, err := u.Repo.GetFollowers(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	hasNext := int64(page*pageSize) < total

	return &PaginatedResponse[userdtos.UsersFollower]{
		Data:        users,
		TotalCount:  total,
		CurrentPage: page,
		PageSize:    pageSize,
		HasNext:     hasNext,
	}, nil
}

func (u *UserConnectionService) GetFollowings(userID int64, page, pageSize int) (*PaginatedResponse[userdtos.UsersFollowing], *apperror.AppError) {
	users, total, err := u.Repo.GetFollowings(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	hasNext := int64(page*pageSize) < total

	return &PaginatedResponse[userdtos.UsersFollowing]{
		Data:        users,
		TotalCount:  total,
		CurrentPage: page,
		PageSize:    pageSize,
		HasNext:     hasNext,
	}, nil
}
