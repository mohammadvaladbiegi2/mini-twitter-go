package tweetreply

import (
	"context"
	"errors"
	"time"

	tweetdtos "twitter_clone/internal/modules/tweet/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateReply(userID, tweetID int64, content string) (tweetdtos.Reply, *apperror.AppError)
	GetReplies(tweetID int64, limit int, cursorID *int64) (tweetdtos.GetRepliesRes, *apperror.AppError)
}

type ReplyRepository struct {
	db *pgxpool.Pool
}

func NewReplyRepository(db *pgxpool.Pool) *ReplyRepository {
	return &ReplyRepository{db: db}
}

func (r *ReplyRepository) CreateReply(userID, tweetID int64, content string) (tweetdtos.Reply, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return tweetdtos.Reply{}, apperror.DB("failed to start tx (create reply)", err)
	}
	defer tx.Rollback(ctx)

	// check tweet exist
	{
		const q = `SELECT 1 FROM tweets WHERE id=$1`
		var d int
		if err := tx.QueryRow(ctx, q, tweetID).Scan(&d); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return tweetdtos.Reply{}, apperror.NotFound("tweet not found", err)
			}
			return tweetdtos.Reply{}, apperror.DB("failed to check tweet existence", err)
		}
	}

	// create reply
	const insert = `
		INSERT INTO tweet_replies (tweet_id, user_id, content)
		VALUES ($1,$2,$3)
		RETURNING id, created_at;
	`
	var (
		replyID   int64
		createdAt time.Time
	)
	if err := tx.QueryRow(ctx, insert, tweetID, userID, content).Scan(&replyID, &createdAt); err != nil {
		return tweetdtos.Reply{}, apperror.DB("failed to insert reply", err)
	}

	// increment reply_count in table tweets
	const bump = `UPDATE tweets SET reply_count = reply_count + 1 WHERE id=$1`
	if _, err := tx.Exec(ctx, bump, tweetID); err != nil {
		return tweetdtos.Reply{}, apperror.DB("failed to bump reply_count", err)
	}

	// get user who reply
	const qUser = `SELECT username, avatar_url FROM users WHERE id=$1`
	var username string
	var avatarURL *string
	if err := tx.QueryRow(ctx, qUser, userID).Scan(&username, &avatarURL); err != nil {
		return tweetdtos.Reply{}, apperror.DB("failed to fetch reply user info", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return tweetdtos.Reply{}, apperror.DB("failed to commit create reply", err)
	}

	return tweetdtos.Reply{
		ID:        replyID,
		TweetID:   tweetID,
		UserID:    userID,
		Username:  username,
		AvatarURL: avatarURL,
		Content:   content,
		CreatedAt: createdAt.Format(time.RFC3339),
	}, nil
}

func (r *ReplyRepository) GetReplies(tweetID int64, limit int, cursorID *int64) (tweetdtos.GetRepliesRes, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// sort if have cursorID
	var rows pgx.Rows
	var err error

	if cursorID != nil {
		const q = `
			SELECT r.id, r.tweet_id, r.user_id, u.username, u.avatar_url, r.content, r.created_at
			FROM tweet_replies r
			JOIN users u ON u.id = r.user_id
			WHERE r.tweet_id = $1 AND r.id < $2
			ORDER BY r.id DESC
			LIMIT $3;
		`
		rows, err = r.db.Query(ctx, q, tweetID, *cursorID, limit)
	} else {
		const q = `
			SELECT r.id, r.tweet_id, r.user_id, u.username, u.avatar_url, r.content, r.created_at
			FROM tweet_replies r
			JOIN users u ON u.id = r.user_id
			WHERE r.tweet_id = $1
			ORDER BY r.id DESC
			LIMIT $2;
		`
		rows, err = r.db.Query(ctx, q, tweetID, limit)
	}

	if err != nil {
		return tweetdtos.GetRepliesRes{}, apperror.DB("failed to fetch replies", err)
	}
	defer rows.Close()

	res := tweetdtos.GetRepliesRes{
		Replies:    make([]tweetdtos.Reply, 0, limit),
		Count:      0,
		NextCursor: nil,
	}

	var lastID int64
	for rows.Next() {
		var (
			rp        tweetdtos.Reply
			createdAt time.Time
		)
		if err := rows.Scan(&rp.ID, &rp.TweetID, &rp.UserID, &rp.Username, &rp.AvatarURL, &rp.Content, &createdAt); err != nil {
			return tweetdtos.GetRepliesRes{}, apperror.DB("failed to scan reply", err)
		}
		rp.CreatedAt = createdAt.Format(time.RFC3339)
		res.Replies = append(res.Replies, rp)
		lastID = rp.ID
	}

	res.Count = len(res.Replies)
	if res.Count == limit {
		res.NextCursor = &lastID
	}

	return res, nil
}
