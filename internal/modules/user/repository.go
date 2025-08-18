package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	userdtos "twitter_clone/internal/modules/user/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	UpdateProfile(userID int64, updateProfileReq userdtos.UpdateProfileReq) (userdtos.UpdateProfileRes, *apperror.AppError)
	GetProfile(userID int64) (userdtos.UserGetProfileRes, *apperror.AppError)
	SearchUsersByUsername(prefix string, limit int) ([]userdtos.SearchUsersByUsernameRes, *apperror.AppError)
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db}
}

func (r UserRepository) UpdateProfile(userID int64, req userdtos.UpdateProfileReq) (userdtos.UpdateProfileRes, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if req.Username != "" {
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND id != $2)`
		if err := r.db.QueryRow(ctx, checkQuery, req.Username, userID).Scan(&exists); err != nil {
			return userdtos.UpdateProfileRes{}, apperror.DB("failed to check username uniqueness", err)
		}
		if exists {
			return userdtos.UpdateProfileRes{}, apperror.Validation("Validation failed", []map[string]string{
				{"username": "Username is already taken"},
			}, nil)
		}
	}

	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if req.Username != "" {
		setClauses = append(setClauses, fmt.Sprintf("username = $%d", argPos))
		args = append(args, req.Username)
		argPos++
	}

	if req.Bio != nil {
		setClauses = append(setClauses, fmt.Sprintf("bio = $%d", argPos))
		args = append(args, req.Bio)
		argPos++
	}

	setClauses = append(setClauses, "updated_at = NOW()")

	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = $%d RETURNING username, bio",
		strings.Join(setClauses, ", "),
		argPos,
	)
	args = append(args, userID)

	var res userdtos.UpdateProfileRes
	if err := r.db.QueryRow(ctx, query, args...).Scan(&res.Username, &res.Bio); err != nil {
		return userdtos.UpdateProfileRes{}, apperror.DB("failed to update profile", err)
	}

	return res, nil
}

func (r UserRepository) GetProfile(userID int64) (userdtos.UserGetProfileRes, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT username, email, bio, avatar_url, follower_count, following_count
		FROM users
		WHERE id = $1
	`

	var res userdtos.UserGetProfileRes
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&res.Username,
		&res.Email,
		&res.Bio,
		&res.AvatarURL,
		&res.FollowerCount,
		&res.FollowingCount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return userdtos.UserGetProfileRes{}, apperror.NotFound("user not found", err)
		}
		return userdtos.UserGetProfileRes{}, apperror.DB("failed to get user profile", err)
	}

	return res, nil
}

func (r *UserRepository) SearchUsersByUsername(prefix string, limit int) ([]userdtos.SearchUsersByUsernameRes, *apperror.AppError) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	const q = `
		SELECT username, bio, avatar_url
		FROM users
		WHERE LOWER(username) LIKE LOWER($1) || '%'
		ORDER BY username
		LIMIT $2;
	`

	rows, err := r.db.Query(ctx, q, prefix, limit)
	if err != nil {
		return nil, apperror.DB("failed to search users", err)
	}
	defer rows.Close()

	users := make([]userdtos.SearchUsersByUsernameRes, 0, limit)
	for rows.Next() {
		var u userdtos.SearchUsersByUsernameRes
		if err := rows.Scan(&u.Username, &u.Bio, &u.AvatarURL); err != nil {
			return nil, apperror.DB("failed to scan user row", err)
		}
		users = append(users, u)
	}
	if rows.Err() != nil {
		return nil, apperror.DB("rows iteration error", rows.Err())
	}

	return users, nil
}
