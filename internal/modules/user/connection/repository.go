package userconnection

import (
	"context"
	"time"
	userdtos "twitter_clone/internal/modules/user/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetFollowers(userID int64, limit, offset int) ([]userdtos.UsersFollower, *apperror.AppError)
	GetFollowings(userID int64, limit, offset int) ([]userdtos.UsersFollowing, *apperror.AppError)
}

type UserConnectionRepository struct {
	db *pgxpool.Pool
}

func NewUserConnectionRepository(db *pgxpool.Pool) *UserConnectionRepository {
	return &UserConnectionRepository{db}
}

// Followers
func (r *UserConnectionRepository) GetFollowers(userID int64, limit, offset int) ([]userdtos.UsersFollower, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const q = `
		SELECT u.username, u.bio, u.avatar_url
		FROM follows f
		JOIN users u ON u.id = f.follower_id
		WHERE f.following_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.db.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, apperror.DB("failed to fetch followers", err)
	}
	defer rows.Close()

	var res []userdtos.UsersFollower
	for rows.Next() {
		var u userdtos.UsersFollower
		if err := rows.Scan(&u.Username, &u.Bio, &u.AvatarURL); err != nil {
			return nil, apperror.DB("failed to scan follower", err)
		}
		res = append(res, u)
	}
	return res, nil
}

// Followings
func (r *UserConnectionRepository) GetFollowings(userID int64, limit, offset int) ([]userdtos.UsersFollowing, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const q = `
		SELECT u.username, u.bio, u.avatar_url
		FROM follows f
		JOIN users u ON u.id = f.following_id
		WHERE f.follower_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.db.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, apperror.DB("failed to fetch followings", err)
	}
	defer rows.Close()

	var res []userdtos.UsersFollowing
	for rows.Next() {
		var u userdtos.UsersFollowing
		if err := rows.Scan(&u.Username, &u.Bio, &u.AvatarURL); err != nil {
			return nil, apperror.DB("failed to scan following", err)
		}
		res = append(res, u)
	}
	return res, nil
}
