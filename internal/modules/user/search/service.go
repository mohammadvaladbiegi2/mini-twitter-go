package usersearch

import (
	userdtos "twitter_clone/internal/modules/user/dto"
	"twitter_clone/internal/pkg/apperror"
)

type Service interface {
	SearchUsersByUsername(prefix string, limit int) ([]userdtos.SearchUsersByUsernameRes, *apperror.AppError)
	GetUserByUsername(userName string) (userdtos.GetUserByUsernameRes, *apperror.AppError)
}

type SearchService struct {
	repo Repository
}

func NewSearchService(repo Repository) *SearchService {
	return &SearchService{repo: repo}
}

func (s *SearchService) SearchUsersByUsername(prefix string, limit int) ([]userdtos.SearchUsersByUsernameRes, *apperror.AppError) {
	return s.repo.SearchUsersByUsername(prefix, limit)
}

func (s *SearchService) GetUserByUsername(userName string) (userdtos.GetUserByUsernameRes, *apperror.AppError) {
	return s.repo.GetUserByUsername(userName)
}
