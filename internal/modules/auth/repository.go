package auth

import (
	"context"
	"fmt"
	"time"
	"twitter_clone/internal/modules/auth/dtos"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	SignUp(userData dtos.SignUpReq) (dtos.SignUpRes, *apperror.AppError)
	Login(userData dtos.LoginReq) (dtos.LoginDBRes, *apperror.AppError)
}

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) Repository {
	return &AuthRepository{db: db}
}

func (r AuthRepository) SignUp(userData dtos.SignUpReq) (dtos.SignUpRes, *apperror.AppError) {

	// if request tacke more 5 second the request canceled
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// query for postgrest databae
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	// requst to Data Base
	var userID int64
	err := r.db.QueryRow(ctx, query, userData.Username, userData.Email, userData.Password).Scan(&userID)
	if err != nil {
		fmt.Println(err)
		return dtos.SignUpRes{}, apperror.DB("failed to insert user", err)
	}

	// TODO Generate token
	return dtos.SignUpRes{
		Token: "generate token",
	}, nil
}

func (r AuthRepository) Login(userData dtos.LoginReq) (dtos.LoginDBRes, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var storedHashedPassword string
	var userName string

	query := `
		SELECT username, password_hash
		FROM users
		WHERE username = $1
		LIMIT 1;
	`

	err := r.db.QueryRow(ctx, query, userData.UserName).Scan(&userName, &storedHashedPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return dtos.LoginDBRes{}, apperror.NotFound("user not found", err)
		}
		return dtos.LoginDBRes{}, apperror.DB("failed to fetch user", err)
	}

	return dtos.LoginDBRes{
		UserName:       userName,
		HashedPassword: storedHashedPassword,
	}, nil
}
