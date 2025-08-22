package userconnection

import (
	"context"
	"time"
	userdtos "twitter_clone/internal/modules/user/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetFollowers(userID int64, page, pageSize int) ([]userdtos.UsersFollower, int64, *apperror.AppError)
	GetFollowings(userID int64, page, pageSize int) ([]userdtos.UsersFollowing, int64, *apperror.AppError)
}

type UserConnectionRepository struct {
	db *pgxpool.Pool
}

func NewUserConnectionRepository(db *pgxpool.Pool) *UserConnectionRepository {
	return &UserConnectionRepository{db}
}

func normalizePagination(page, pageSize int) (int, int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return page, pageSize, offset
}

// Followers
func (r *UserConnectionRepository) GetFollowers(userID int64, page, pageSize int) ([]userdtos.UsersFollower, int64, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	page, pageSize, offset := normalizePagination(page, pageSize)

	//  total count
	const countQ = `SELECT COUNT(*) FROM follows WHERE following_id = $1`
	var total int64
	if err := r.db.QueryRow(ctx, countQ, userID).Scan(&total); err != nil {
		return nil, 0, apperror.DB("failed to count followers", err)
	}

	const q = `
		SELECT u.username, u.bio, u.avatar_url
		FROM follows f
		JOIN users u ON u.id = f.follower_id
		WHERE f.following_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.db.Query(ctx, q, userID, pageSize, offset)
	if err != nil {
		return nil, 0, apperror.DB("failed to fetch followers", err)
	}
	defer rows.Close()

	var res []userdtos.UsersFollower
	for rows.Next() {
		var u userdtos.UsersFollower
		if err := rows.Scan(&u.Username, &u.Bio, &u.AvatarURL); err != nil {
			return nil, 0, apperror.DB("failed to scan follower", err)
		}
		res = append(res, u)
	}
	return res, total, nil
}

// Followings
func (r *UserConnectionRepository) GetFollowings(userID int64, page, pageSize int) ([]userdtos.UsersFollowing, int64, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	page, pageSize, offset := normalizePagination(page, pageSize)

	//  total count
	const countQ = `SELECT COUNT(*) FROM follows WHERE follower_id = $1`
	var total int64
	if err := r.db.QueryRow(ctx, countQ, userID).Scan(&total); err != nil {
		return nil, 0, apperror.DB("failed to count followings", err)
	}

	const q = `
		SELECT u.username, u.bio, u.avatar_url
		FROM follows f
		JOIN users u ON u.id = f.following_id
		WHERE f.follower_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.db.Query(ctx, q, userID, pageSize, offset)
	if err != nil {
		return nil, 0, apperror.DB("failed to fetch followings", err)
	}
	defer rows.Close()

	var res []userdtos.UsersFollowing
	for rows.Next() {
		var u userdtos.UsersFollowing
		if err := rows.Scan(&u.Username, &u.Bio, &u.AvatarURL); err != nil {
			return nil, 0, apperror.DB("failed to scan following", err)
		}
		res = append(res, u)
	}
	return res, total, nil
}
