package tweetbookmark

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	tweetdtos "twitter_clone/internal/modules/tweet/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	AddBookmark(userID, tweetID int64) *apperror.AppError
	RemoveBookmark(userID, tweetID int64) *apperror.AppError
	ListBookmarks(userID int64, limit int, afterCreatedAt *time.Time, afterTweetID *int64) ([]tweetdtos.BookmarkTweet, *time.Time, *int64, *apperror.AppError)
}

type BookmarkRepository struct {
	db *pgxpool.Pool
}

func NewBookMarkRepository(db *pgxpool.Pool) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

func (r *BookmarkRepository) AddBookmark(userID, tweetID int64) *apperror.AppError {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperror.DB("bookmark.begin", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// وجود توییت (برای خطای بهتر نسبت به FK)
	const checkTweet = `SELECT 1 FROM tweets WHERE id=$1`
	var d int
	if err := tx.QueryRow(ctx, checkTweet, tweetID).Scan(&d); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.NotFound("tweet not found", err)
		}
		return apperror.DB("bookmark.check_tweet", err)
	}

	// درج idempotent
	const ins = `
		INSERT INTO tweet_bookmarks (user_id, tweet_id)
		VALUES ($1,$2)
		ON CONFLICT (user_id, tweet_id) DO NOTHING;
	`
	tag, err := tx.Exec(ctx, ins, userID, tweetID)
	if err != nil {
		return apperror.DB("bookmark.insert", err)
	}

	// اگر درج شد، شمارنده را زیاد کن (اگه ستون را اضافه کرده‌ای)
	if tag.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE tweets SET bookmark_count=bookmark_count+1 WHERE id=$1`, tweetID); err != nil {
			return apperror.DB("bookmark.bump_count", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return apperror.DB("bookmark.commit", err)
	}
	return nil
}

func (r *BookmarkRepository) RemoveBookmark(userID, tweetID int64) *apperror.AppError {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperror.DB("unbookmark.begin", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	const del = `DELETE FROM tweet_bookmarks WHERE user_id=$1 AND tweet_id=$2`
	tag, err := tx.Exec(ctx, del, userID, tweetID)
	if err != nil {
		return apperror.DB("unbookmark.delete", err)
	}
	if tag.RowsAffected() == 0 {
		return apperror.NotFound("bookmark not found", nil)
	}

	// شمارنده را کم کن (اگه ستون را اضافه کرده‌ای)
	if _, err := tx.Exec(ctx, `UPDATE tweets SET bookmark_count=GREATEST(bookmark_count-1,0) WHERE id=$1`, tweetID); err != nil {
		return apperror.DB("unbookmark.dec_count", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return apperror.DB("unbookmark.commit", err)
	}
	return nil
}

func (r *BookmarkRepository) ListBookmarks(userID int64, limit int, afterCreatedAt *time.Time, afterTweetID *int64) ([]tweetdtos.BookmarkTweet, *time.Time, *int64, *apperror.AppError) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	// keyset: (b.created_at DESC, b.tweet_id DESC)
	const base = `
SELECT
  t.id,
  t.content,
  t.like_count,
  t.dislike_count,
  t.reply_count,
  t.created_at,
  u.username,
  u.avatar_url,
  COALESCE(ARRAY_AGG(DISTINCT tg.name) FILTER (WHERE tg.name IS NOT NULL), '{}') AS tags,
  b.created_at AS b_created_at -- برای cursor
FROM tweet_bookmarks b
JOIN tweets t ON t.id = b.tweet_id
JOIN users u  ON u.id = t.user_id
LEFT JOIN tweet_tags tt ON tt.tweet_id = t.id
LEFT JOIN tags tg       ON tg.id       = tt.tag_id
WHERE b.user_id = $1
`

	var rows pgx.Rows
	var err error

	if afterCreatedAt != nil && afterTweetID != nil {
		// صفحه‌های بعدی
		const q = base + `
  AND (b.created_at, b.tweet_id) < ($2, $3)
GROUP BY t.id, u.username, u.avatar_url, b_created_at
ORDER BY b_created_at DESC, t.id DESC
LIMIT $4;
`
		rows, err = r.db.Query(ctx, q, userID, *afterCreatedAt, *afterTweetID, limit)
	} else {
		// صفحه اول
		const q = base + `
GROUP BY t.id, u.username, u.avatar_url, b_created_at
ORDER BY b_created_at DESC, t.id DESC
LIMIT $2;
`
		rows, err = r.db.Query(ctx, q, userID, limit)
	}
	if err != nil {
		return nil, nil, nil, apperror.DB("bookmarks.list.query", err)
	}
	defer rows.Close()

	list := make([]tweetdtos.BookmarkTweet, 0, limit)
	var lastCreatedAt *time.Time
	var lastTweetID *int64

	for rows.Next() {
		var it tweetdtos.BookmarkTweet
		var tCreatedAt time.Time
		var tags []string
		var bCreatedAt time.Time

		if err := rows.Scan(
			&it.ID,
			&it.Content,
			&it.LikeCount,
			&it.DislikeCount,
			&it.ReplyCount,
			&tCreatedAt,
			&it.AuthorUsername,
			&it.AuthorAvatar,
			&tags,
			&bCreatedAt,
		); err != nil {
			return nil, nil, nil, apperror.DB("bookmarks.list.scan", err)
		}

		it.CreatedAt = tCreatedAt.Format(time.RFC3339)
		it.Tags = tags
		list = append(list, it)

		// برای next-cursor
		lastCreatedAt = &bCreatedAt
		lastID := it.ID
		lastTweetID = &lastID
	}

	return list, lastCreatedAt, lastTweetID, nil
}

/* --- Helpers برای cursor (در سرویس استفاده می‌کنیم) --- */
func EncodeCursor(ts time.Time, id int64) string {
	raw := fmt.Sprintf("%d|%d", ts.UnixNano(), id)
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func DecodeCursor(cur string) (*time.Time, *int64, error) {
	if cur == "" {
		return nil, nil, nil
	}
	b, err := base64.StdEncoding.DecodeString(cur)
	if err != nil {
		return nil, nil, err
	}
	parts := strings.Split(string(b), "|")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid cursor")
	}
	ns, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, nil, err
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, nil, err
	}
	t := time.Unix(0, ns).UTC()
	return &t, &id, nil
}
