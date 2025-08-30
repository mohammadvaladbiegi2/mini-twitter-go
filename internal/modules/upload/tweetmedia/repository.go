package tweetmedia

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TweetRepository struct {
	db *pgxpool.Pool
}

func NewTweetRepository(db *pgxpool.Pool) *TweetRepository {
	return &TweetRepository{db: db}
}

func (r *TweetRepository) UpdateTweetImage(tweetID, uploadID int64, imageURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx,
		`UPDATE tweets SET upload_id=$1, image_url=$2 WHERE id=$3`,
		uploadID, imageURL, tweetID,
	)
	return err
}

func (r *TweetRepository) IsOwner(tweetID, userID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var count int
	query := `SELECT COUNT(1) FROM tweets WHERE id = $1 AND user_id = $2`
	err := r.db.QueryRow(ctx, query, tweetID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
