package useraction

import (
	"context"
	"time"

	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Follow(followerID, followeeID int64) *apperror.AppError
	Unfollow(followerID, followeeID int64) *apperror.AppError
}

type UserActionRepository struct {
	db *pgxpool.Pool
}

func NewUserActionRepository(db *pgxpool.Pool) *UserActionRepository {
	return &UserActionRepository{db}
}

// Follow user
func (r *UserActionRepository) Follow(followerID, followeeID int64) *apperror.AppError {
	if followerID == followeeID {
		return apperror.Validation("you cannot follow yourself", nil, nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperror.DB("failed to start transaction", err)
	}
	defer tx.Rollback(ctx)

	// if befor this action exist do nothing
	const insert = `
        INSERT INTO follows (follower_id, following_id, created_at)
        VALUES ($1,$2,NOW())
        ON CONFLICT (follower_id, following_id) DO NOTHING;
    `
	tag, err := tx.Exec(ctx, insert, followerID, followeeID)
	if err != nil {
		return apperror.DB("failed to insert relation", err)
	}
	if tag.RowsAffected() == 0 {
		// return already following
		return apperror.Validation("already following", nil, nil)
	}

	// increment follower who followed and incremnt following who follow
	const incFollowing = `UPDATE users SET following_count = following_count + 1 WHERE id=$1;`
	const incFollower = `UPDATE users SET follower_count  = follower_count  + 1 WHERE id=$1;`

	if _, err := tx.Exec(ctx, incFollowing, followerID); err != nil {
		return apperror.DB("failed to update following_count", err)
	}
	if _, err := tx.Exec(ctx, incFollower, followeeID); err != nil {
		return apperror.DB("failed to update follower_count", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return apperror.DB("failed to commit follow", err)
	}
	return nil
}

func (r *UserActionRepository) Unfollow(followerID, followeeID int64) *apperror.AppError {
	if followerID == followeeID {
		return apperror.Validation("you cannot unfollow yourself", nil, nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperror.DB("failed to start transaction", err)
	}
	defer tx.Rollback(ctx)

	const del = `DELETE FROM follows WHERE follower_id=$1 AND following_id=$2;`
	tag, err := tx.Exec(ctx, del, followerID, followeeID)
	if err != nil {
		return apperror.DB("failed to delete relation", err)
	}
	if tag.RowsAffected() == 0 {
		return apperror.NotFound("relation not found", nil)
	}

	const decFollowing = `UPDATE users SET following_count = GREATEST(following_count - 1, 0) WHERE id=$1;`
	const decFollower = `UPDATE users SET follower_count  = GREATEST(follower_count  - 1, 0) WHERE id=$1;`

	if _, err := tx.Exec(ctx, decFollowing, followerID); err != nil {
		return apperror.DB("failed to update following_count", err)
	}
	if _, err := tx.Exec(ctx, decFollower, followeeID); err != nil {
		return apperror.DB("failed to update follower_count", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return apperror.DB("failed to commit unfollow", err)
	}
	return nil
}
