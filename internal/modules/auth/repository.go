package auth

import (
	"context"
	"time"
	"twitter_clone/internal/modules/auth/dtos"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	SignUp(userData dtos.SignUpReq) (dtos.SignUpResDB, *apperror.AppError)
	Login(userData dtos.LoginReq) (dtos.LoginDBRes, *apperror.AppError)
}

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) Repository {
	return &AuthRepository{db: db}
}

func (r AuthRepository) SignUp(userData dtos.SignUpReq) (dtos.SignUpResDB, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checkQuery := `
		SELECT id FROM users 
		WHERE username = $1 OR email = $2
	`
	var existingID int64
	err := r.db.QueryRow(ctx, checkQuery, userData.Username, userData.Email).Scan(&existingID)
	if err == nil {
		return dtos.SignUpResDB{}, apperror.Validation("user already exists", nil, nil)
	} else if err != pgx.ErrNoRows {
		return dtos.SignUpResDB{}, apperror.DB("failed to check existing user", err)
	}

	insertQuery := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, username;
	`
	var userID int64
	var username string
	err = r.db.QueryRow(ctx, insertQuery, userData.Username, userData.Email, userData.Password).Scan(&userID, &username)
	if err != nil {
		return dtos.SignUpResDB{}, apperror.DB("failed to insert user", err)
	}

	return dtos.SignUpResDB{
		ID:       userID,
		UserName: username,
	}, nil
}

func (r AuthRepository) Login(userData dtos.LoginReq) (dtos.LoginDBRes, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var storedHashedPassword string
	var userName string
	var userID int64

	query := `
		SELECT username, id, password_hash
		FROM users
		WHERE username = $1
		LIMIT 1;
	`

	err := r.db.QueryRow(ctx, query, userData.UserName).Scan(&userName, &userID, &storedHashedPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return dtos.LoginDBRes{}, apperror.NotFound("user not found", err)
		}
		return dtos.LoginDBRes{}, apperror.DB("failed to fetch user", err)
	}

	return dtos.LoginDBRes{
		ID:             userID,
		UserName:       userName,
		HashedPassword: storedHashedPassword,
	}, nil
}
