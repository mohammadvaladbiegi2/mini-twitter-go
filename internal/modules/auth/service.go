package auth

import (
	"fmt"
	"twitter_clone/internal/modules/auth/dtos"
	"twitter_clone/internal/pkg/apperror"
	"twitter_clone/internal/pkg/jwt"
	"twitter_clone/internal/pkg/validation"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	SignUp(userData dtos.SignUpReq) (dtos.SignUpRes, *apperror.AppError)
	Login(userData dtos.LoginReq) (dtos.LoginRes, *apperror.AppError)
}

type AuthService struct {
	repo Repository
}

func NewAuthService(repo Repository) Service {
	return &AuthService{repo: repo}
}

func (r AuthService) SignUp(userData dtos.SignUpReq) (dtos.SignUpRes, *apperror.AppError) {

	//  validation user Request
	if validErrors := validation.ValidateSignUpReq(userData); validErrors != nil {
		return dtos.SignUpRes{}, validErrors
	}

	// hash password by bcrypt package
	hashedPassword, Herr := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if Herr != nil {
		return dtos.SignUpRes{}, apperror.Server("failed to hash password", Herr)
	}
	userData.Password = string(hashedPassword)

	// create pool request to database
	user, err := r.repo.SignUp(userData)
	fmt.Println(user)
	if err != nil {
		return dtos.SignUpRes{}, err
	}

	// generate token
	token, Terr := jwt.BuildToken(user.UserName, user.ID)
	if Terr != nil {
		return dtos.SignUpRes{}, Terr
	}

	return dtos.SignUpRes{Token: token}, nil
}

func (r AuthService) Login(userData dtos.LoginReq) (dtos.LoginRes, *apperror.AppError) {
	if validErrors := validation.ValidateLoginReq(userData); validErrors != nil {
		return dtos.LoginRes{}, validErrors
	}

	user, err := r.repo.Login(userData)
	if err != nil {
		return dtos.LoginRes{}, err
	}

	Cerror := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userData.Password))
	if Cerror != nil {
		return dtos.LoginRes{}, apperror.UnauthorizedErr("invalid username or password", Cerror)
	}

	token, Terr := jwt.BuildToken(user.UserName, user.ID)
	if Terr != nil {
		return dtos.LoginRes{}, Terr
	}

	return dtos.LoginRes{
		Token: token,
	}, nil
}
