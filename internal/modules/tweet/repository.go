package tweet

import (
	"context"
	"strings"
	"time"
	tweetdtos "twitter_clone/internal/modules/tweet/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateTweet(UserID int64, req tweetdtos.CreateTweetReq) (tweetdtos.CreateTweetRes, *apperror.AppError)
}

type TweetRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *TweetRepository {
	return &TweetRepository{db: db}
}

func (r *TweetRepository) CreateTweet(UserID int64, req tweetdtos.CreateTweetReq) (tweetdtos.CreateTweetRes, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return tweetdtos.CreateTweetRes{}, apperror.DB("failed to start transaction", err)
	}
	defer tx.Rollback(ctx)

	// 1. create tweet
	var tweetID int64
	insertTweetQuery := `
        INSERT INTO tweets (user_id, content)
        VALUES ($1, $2)
        RETURNING id;
    `
	if err := tx.QueryRow(ctx, insertTweetQuery, UserID, req.Content).Scan(&tweetID); err != nil {
		return tweetdtos.CreateTweetRes{}, apperror.DB("failed to insert tweet", err)
	}

	tagIDs := []int64{}
	uniqueTags := map[string]struct{}{}

	// 2. remove extra tag in request
	for _, t := range req.Tags {
		cleaned := strings.TrimSpace(strings.ToLower(t))
		if cleaned == "" {
			continue
		}
		uniqueTags[cleaned] = struct{}{}
	}

	if len(uniqueTags) > 0 {
		tags := make([]string, 0, len(uniqueTags))
		for t := range uniqueTags {
			tags = append(tags, t)
		}

		// 3. insert all tage just once time and ignoar extra tags
		insertTagsQuery := `
            INSERT INTO tags (name)
            SELECT unnest($1::text[])
            ON CONFLICT (name) DO NOTHING;
        `
		if _, err := tx.Exec(ctx, insertTagsQuery, tags); err != nil {
			return tweetdtos.CreateTweetRes{}, apperror.DB("failed to insert tags", err)
		}

		// 4. get all tags oldest and newest
		selectTagsQuery := `SELECT id FROM tags WHERE name = ANY($1)`
		rows, err := tx.Query(ctx, selectTagsQuery, tags)
		if err != nil {
			return tweetdtos.CreateTweetRes{}, apperror.DB("failed to fetch tag ids", err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int64
			if err := rows.Scan(&id); err != nil {
				return tweetdtos.CreateTweetRes{}, apperror.DB("failed to scan tag id", err)
			}
			tagIDs = append(tagIDs, id)
		}

		// 5. bulk insert tags to tweets
		insertTweetTagsQuery := `
            INSERT INTO tweet_tags (tweet_id, tag_id)
            SELECT $1, unnest($2::bigint[])
            ON CONFLICT DO NOTHING;
        `
		if _, err := tx.Exec(ctx, insertTweetTagsQuery, tweetID, tagIDs); err != nil {
			return tweetdtos.CreateTweetRes{}, apperror.DB("failed to insert tweet_tags", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return tweetdtos.CreateTweetRes{}, apperror.DB("failed to commit transaction", err)
	}

	// creat respons tags
	tagsOut := make([]string, 0, len(uniqueTags))
	for t := range uniqueTags {
		tagsOut = append(tagsOut, t)
	}

	return tweetdtos.CreateTweetRes{
		ID:       tweetID,
		UserID:   UserID,
		Content:  req.Content,
		ImageURL: nil,
		Tags:     tagsOut,
	}, nil
}
