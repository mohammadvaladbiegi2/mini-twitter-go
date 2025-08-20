package tweetaction

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	React(userID, tweetID int64, isLike bool) *apperror.AppError
	RemoveReaction(userID, tweetID int64) *apperror.AppError
}

type UserActionRepository struct {
	db *pgxpool.Pool
}

func NewUserActionRepository(db *pgxpool.Pool) *UserActionRepository {
	return &UserActionRepository{db}
}

// Like or Dislike a tweet
func (r *UserActionRepository) React(userID, tweetID int64, isLike bool) *apperror.AppError {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperror.DB("failed to start transaction", err)
	}
	defer tx.Rollback(ctx)

	// check for exist tweet
	{
		const checkTweet = `SELECT 1 FROM tweets WHERE id=$1`
		var dummy int
		if err := tx.QueryRow(ctx, checkTweet, tweetID).Scan(&dummy); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				slog.Warn("react.tweet_not_found", "tweetID", tweetID)
				return apperror.NotFound("tweet not found", err)
			}
			return apperror.DB("failed to check tweet", err)
		}
	}

	// check old user reaction
	const check = `SELECT is_like FROM tweet_likes WHERE user_id=$1 AND tweet_id=$2`
	var prev sql.NullBool
	err = tx.QueryRow(ctx, check, userID, tweetID).Scan(&prev)
	slog.Debug("react.check_like", "userID", userID, "tweetID", tweetID, "err", err)

	switch {
	case err == nil:
		if prev.Valid && prev.Bool == isLike {
			return apperror.Validation("already reacted with same type", nil, nil)
		}

		const updateReaction = `UPDATE tweet_likes SET is_like=$1, created_at=NOW() WHERE user_id=$2 AND tweet_id=$3`
		if _, err := tx.Exec(ctx, updateReaction, isLike, userID, tweetID); err != nil {
			return apperror.DB("failed to update reaction", err)
		}

		if isLike {
			// dislike -> like
			_, err := tx.Exec(ctx, `UPDATE tweets SET like_count=like_count+1, dislike_count=GREATEST(dislike_count-1,0) WHERE id=$1`, tweetID)
			if err != nil {
				return apperror.DB("failed to update counts", err)
			}
		} else {
			// like -> dislike
			_, err := tx.Exec(ctx, `UPDATE tweets SET like_count=GREATEST(like_count-1,0), dislike_count=dislike_count+1 WHERE id=$1`, tweetID)
			if err != nil {
				return apperror.DB("failed to update counts", err)
			}
		}

	case errors.Is(err, pgx.ErrNoRows):
		// if dont have old reaction → INSERT -> new reaction
		const insert = `INSERT INTO tweet_likes (user_id, tweet_id, is_like, created_at) VALUES ($1,$2,$3,NOW())`
		if _, err := tx.Exec(ctx, insert, userID, tweetID, isLike); err != nil {
			return apperror.DB("failed to insert reaction", err)
		}

		if isLike {
			_, err := tx.Exec(ctx, `UPDATE tweets SET like_count=like_count+1 WHERE id=$1`, tweetID)
			if err != nil {
				return apperror.DB("failed to update like_count", err)
			}
		} else {
			_, err := tx.Exec(ctx, `UPDATE tweets SET dislike_count=dislike_count+1 WHERE id=$1`, tweetID)
			if err != nil {
				return apperror.DB("failed to update dislike_count", err)
			}
		}

	default:
		return apperror.DB("failed to check reaction", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return apperror.DB("failed to commit transaction", err)
	}
	return nil
}

// (unlike / undislike)
func (r *UserActionRepository) RemoveReaction(userID, tweetID int64) *apperror.AppError {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperror.DB("failed to start transaction", err)
	}
	defer tx.Rollback(ctx)

	const check = `SELECT is_like FROM tweet_likes WHERE user_id=$1 AND tweet_id=$2`
	var isLike bool
	err = tx.QueryRow(ctx, check, userID, tweetID).Scan(&isLike)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Warn("remove.not_found", "userID", userID, "tweetID", tweetID)
		return apperror.NotFound("reaction not found", err)
	} else if err != nil {
		return apperror.DB("failed to check reaction", err)
	}

	const del = `DELETE FROM tweet_likes WHERE user_id=$1 AND tweet_id=$2`
	if _, err := tx.Exec(ctx, del, userID, tweetID); err != nil {
		return apperror.DB("failed to delete reaction", err)
	}

	if isLike {
		_, err := tx.Exec(ctx, `UPDATE tweets SET like_count=GREATEST(like_count-1,0) WHERE id=$1`, tweetID)
		if err != nil {
			return apperror.DB("failed to update like_count", err)
		}
	} else {
		_, err := tx.Exec(ctx, `UPDATE tweets SET dislike_count=GREATEST(dislike_count-1,0) WHERE id=$1`, tweetID)
		if err != nil {
			return apperror.DB("failed to update dislike_count", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return apperror.DB("failed to commit transaction", err)
	}
	return nil
}
