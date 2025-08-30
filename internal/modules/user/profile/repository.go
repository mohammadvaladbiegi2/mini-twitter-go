package userprofile

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	userdtos "twitter_clone/internal/modules/user/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	UpdateProfile(userID int64, updateProfileReq userdtos.UpdateProfileReq) (userdtos.UpdateProfileRes, *apperror.AppError)
	GetProfile(userID int64) (userdtos.UserGetProfileRes, *apperror.AppError)
}

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db}
}

func (r ProfileRepository) UpdateProfile(userID int64, req userdtos.UpdateProfileReq) (userdtos.UpdateProfileRes, *apperror.AppError) {
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

func (r *ProfileRepository) GetProfile(userID int64) (userdtos.UserGetProfileRes, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	queryUser := `
		SELECT username, email, bio, avatar_url, follower_count, following_count
		FROM users
		WHERE id = $1
	`

	var res userdtos.UserGetProfileRes
	err := r.db.QueryRow(ctx, queryUser, userID).Scan(
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

	queryTweets := `
    SELECT
        t.id,
        t.content,
        t.like_count,
        t.dislike_count,
		t.reply_count,
		t.bookmark_count,
		t.image_url,
        t.created_at,
        COALESCE(ARRAY_AGG(DISTINCT tg.name) FILTER (WHERE tg.name IS NOT NULL), '{}') AS tags
    FROM tweets t
    LEFT JOIN tweet_tags tt ON tt.tweet_id = t.id
    LEFT JOIN tags tg       ON tg.id = tt.tag_id
    WHERE t.user_id = $1
    GROUP BY t.id
    ORDER BY t.created_at DESC;
	`

	rows, err := r.db.Query(ctx, queryTweets, userID)
	if err != nil {
		return res, apperror.DB("failed to get user tweets", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tw userdtos.TweetForProfile
		var createdAt time.Time
		var tags []string

		if err := rows.Scan(&tw.ID, &tw.Content, &tw.LikeCount, &tw.DislikeCount, &tw.ReplyCount, &tw.BookMarkCount, &tw.ImageURL, &createdAt, &tags); err != nil {
			return res, apperror.DB("failed to scan tweet", err)
		}

		tw.CreatedAt = createdAt.Format(time.RFC3339)
		tw.Tags = tags

		res.Tweets = append(res.Tweets, tw)
	}

	return res, nil
}
