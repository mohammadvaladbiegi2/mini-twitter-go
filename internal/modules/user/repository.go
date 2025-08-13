package user

import (
	"context"
	"fmt"
	"time"
	"twitter_clone/internal/modules/user/dtos"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	SignUp(userData dtos.UserSignUpReq) (dtos.UserSignUpRes, *apperror.AppError)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r userRepository) SignUp(userData dtos.UserSignUpReq) (dtos.UserSignUpRes, *apperror.AppError) {

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
		return dtos.UserSignUpRes{}, apperror.DB("failed to insert user", err)
	}

	return dtos.UserSignUpRes{
		Message: "user created successfully",
	}, nil
}
