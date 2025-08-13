package user

import (
	"twitter_clone/internal/modules/user/dtos"
	"twitter_clone/internal/pkg/apperror"
	"twitter_clone/internal/pkg/validation"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUp(userData dtos.UserSignUpReq) (dtos.UserSignUpRes, *apperror.AppError)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (r userService) SignUp(userData dtos.UserSignUpReq) (dtos.UserSignUpRes, *apperror.AppError) {

	//  validation user Request
	if validErrors := validation.ValidateSignUpReq(userData); validErrors != nil {
		return dtos.UserSignUpRes{}, validErrors
	}

	// hash password by bcrypt package
	hashedPassword, Herr := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if Herr != nil {
		return dtos.UserSignUpRes{}, apperror.Server("failed to hash password", Herr)
	}
	userData.Password = string(hashedPassword)

	// create pool request to database
	message, err := r.repo.SignUp(userData)
	if err != nil {
		return dtos.UserSignUpRes{}, err
	}

	return dtos.UserSignUpRes{Message: message.Message}, nil
}
