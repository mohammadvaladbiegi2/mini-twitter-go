package auth

import (
	"context"
	"time"
	authdtos "twitter_clone/internal/modules/auth/dtos"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	SignUp(userData authdtos.SignUpReq) (authdtos.SignUpResDB, *apperror.AppError)
	Login(userData authdtos.LoginReq) (authdtos.LoginDBRes, *apperror.AppError)
}

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) Repository {
	return &AuthRepository{db: db}
}

func (r AuthRepository) SignUp(userData authdtos.SignUpReq) (authdtos.SignUpResDB, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	checkUsernameQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	if err := r.db.QueryRow(ctx, checkUsernameQuery, userData.Username).Scan(&exists); err != nil {
		return authdtos.SignUpResDB{}, apperror.DB("failed to check username uniqueness", err)
	}
	if exists {
		return authdtos.SignUpResDB{}, apperror.Validation("Validation failed", []map[string]string{
			{"username": "Username is already taken"},
		}, nil)
	}

	checkEmailQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	if err := r.db.QueryRow(ctx, checkEmailQuery, userData.Email).Scan(&exists); err != nil {
		return authdtos.SignUpResDB{}, apperror.DB("failed to check email uniqueness", err)
	}
	if exists {
		return authdtos.SignUpResDB{}, apperror.Validation("Validation failed", []map[string]string{
			{"email": "Email is already registered"},
		}, nil)
	}

	// شروع تراکنش
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return authdtos.SignUpResDB{}, apperror.DB("failed to start transaction", err)
	}
	defer tx.Rollback(ctx)

	// insert user
	insertUserQuery := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, username
	`
	var userID int64
	var username string
	if err := tx.QueryRow(ctx, insertUserQuery, userData.Username, userData.Email, userData.Password).Scan(&userID, &username); err != nil {
		return authdtos.SignUpResDB{}, apperror.DB("failed to insert user", err)
	}

	// insert wallet
	insertWalletQuery := `
		INSERT INTO wallets (user_id, balance)
		VALUES ($1, 0)
	`
	if _, err := tx.Exec(ctx, insertWalletQuery, userID); err != nil {
		return authdtos.SignUpResDB{}, apperror.DB("failed to create wallet", err)
	}

	// commit تراکنش
	if err := tx.Commit(ctx); err != nil {
		return authdtos.SignUpResDB{}, apperror.DB("failed to commit transaction", err)
	}

	return authdtos.SignUpResDB{
		ID:       userID,
		UserName: username,
	}, nil
}

func (r AuthRepository) Login(userData authdtos.LoginReq) (authdtos.LoginDBRes, *apperror.AppError) {
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
			return authdtos.LoginDBRes{}, apperror.NotFound("user not found", err)
		}
		return authdtos.LoginDBRes{}, apperror.DB("failed to fetch user", err)
	}

	return authdtos.LoginDBRes{
		ID:             userID,
		UserName:       userName,
		HashedPassword: storedHashedPassword,
	}, nil
}
