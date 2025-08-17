package auth

import (
	"fmt"
	authdtos "twitter_clone/internal/modules/auth/dtos"
	"twitter_clone/internal/pkg/apperror"
	"twitter_clone/internal/pkg/jwt"
	"twitter_clone/internal/pkg/validation"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	SignUp(userData authdtos.SignUpReq) (authdtos.SignUpRes, *apperror.AppError)
	Login(userData authdtos.LoginReq) (authdtos.LoginRes, *apperror.AppError)
}

type AuthService struct {
	repo Repository
}

func NewAuthService(repo Repository) Service {
	return &AuthService{repo: repo}
}

func (r AuthService) SignUp(userData authdtos.SignUpReq) (authdtos.SignUpRes, *apperror.AppError) {

	//  validation user Request
	if validErrors := validation.ValidateSignUpReq(userData); validErrors != nil {
		return authdtos.SignUpRes{}, validErrors
	}

	// hash password by bcrypt package
	hashedPassword, Herr := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if Herr != nil {
		return authdtos.SignUpRes{}, apperror.Server("failed to hash password", Herr)
	}
	userData.Password = string(hashedPassword)

	// create pool request to database
	user, err := r.repo.SignUp(userData)
	fmt.Println(user)
	if err != nil {
		return authdtos.SignUpRes{}, err
	}

	// generate token
	token, Terr := jwt.BuildToken(user.UserName, user.ID)
	if Terr != nil {
		return authdtos.SignUpRes{}, Terr
	}

	return authdtos.SignUpRes{Token: token}, nil
}

func (r AuthService) Login(userData authdtos.LoginReq) (authdtos.LoginRes, *apperror.AppError) {
	if validErrors := validation.ValidateLoginReq(userData); validErrors != nil {
		return authdtos.LoginRes{}, validErrors
	}

	user, err := r.repo.Login(userData)
	if err != nil {
		return authdtos.LoginRes{}, err
	}

	Cerror := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userData.Password))
	if Cerror != nil {
		return authdtos.LoginRes{}, apperror.UnauthorizedErr("invalid username or password", Cerror)
	}

	token, Terr := jwt.BuildToken(user.UserName, user.ID)
	if Terr != nil {
		return authdtos.LoginRes{}, Terr
	}

	return authdtos.LoginRes{
		Token: token,
	}, nil
}
