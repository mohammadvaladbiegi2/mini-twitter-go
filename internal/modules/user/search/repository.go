package usersearch

import (
	"context"
	"errors"
	"time"

	userdtos "twitter_clone/internal/modules/user/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	SearchUsersByUsername(prefix string, limit int) ([]userdtos.SearchUsersByUsernameRes, *apperror.AppError)
	GetUserByUsername(userName string, userID int64) (userdtos.GetUserByUsernameRes, *apperror.AppError)
}

type SearchRepository struct {
	db *pgxpool.Pool
}

func NewSearchRepository(db *pgxpool.Pool) *SearchRepository {
	return &SearchRepository{db}
}

func (r *SearchRepository) SearchUsersByUsername(prefix string, limit int) ([]userdtos.SearchUsersByUsernameRes, *apperror.AppError) {
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

func (r *SearchRepository) GetUserByUsername(userName string, userID int64) (userdtos.GetUserByUsernameRes, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	queryUser := `
		SELECT username, bio, avatar_url, follower_count, following_count
		FROM users
		WHERE username = $1
	`
	var res userdtos.GetUserByUsernameRes
	err := r.db.QueryRow(ctx, queryUser, userName).Scan(
		&res.Username,
		&res.Bio,
		&res.AvatarURL,
		&res.FollowerCount,
		&res.FollowingCount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return userdtos.GetUserByUsernameRes{}, apperror.NotFound("user not found", err)
		}
		return userdtos.GetUserByUsernameRes{}, apperror.DB("failed to get user profile", err)
	}

	queryTweets := `
    SELECT
        t.id,
        t.content,
        t.like_count,
        t.dislike_count,
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
		var tw userdtos.Tweets
		var createdAt time.Time
		var tags []string

		if err := rows.Scan(&tw.ID, &tw.Content, &tw.LikeCount, &tw.DislikeCount, &createdAt, &tags); err != nil {
			return res, apperror.DB("failed to scan tweet", err)
		}

		tw.CreatedAt = createdAt.Format(time.RFC3339)
		tw.Tags = tags

		res.Tweets = append(res.Tweets, tw)
	}

	return res, nil
}
