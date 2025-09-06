package timeline

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	timelinedto "twitter_clone/internal/modules/timeline/dto"
	"twitter_clone/internal/pkg/apperror"
	"twitter_clone/internal/pkg/redisclient"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	MyTimeline(userID int64, limit int, cursor *string) (*timelinedto.TimelineResponse, *apperror.AppError)
}

type TimeLineRepo struct {
	db *pgxpool.Pool
}

func NewTimeLineRepo(db *pgxpool.Pool) *TimeLineRepo {
	return &TimeLineRepo{db}
}

func (r *TimeLineRepo) MyTimeline(
	userID int64,
	limit int,
	cursor *string,
) (*timelinedto.TimelineResponse, *apperror.AppError) {

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// --- Step 1: Get following IDs ---
	followingIDs, appErr := r.getFollowingIDs(ctx, userID)
	if appErr != nil {
		return nil, appErr
	}

	if len(followingIDs) == 0 {
		return &timelinedto.TimelineResponse{
			Tweets:      []timelinedto.Tweet{},
			NextCursor:  nil,
			HasNextPage: false,
			Limit:       limit,
		}, nil
	}

	// --- Step 2: Parse cursor ---
	var afterCreatedAt time.Time
	var afterTweetID int64
	hasCursor := false

	if cursor != nil && *cursor != "" {
		if decoded, err := base64.URLEncoding.DecodeString(*cursor); err == nil {
			parts := strings.Split(string(decoded), "|")
			if len(parts) == 2 {
				if ts, err := time.Parse(time.RFC3339Nano, parts[0]); err == nil {
					afterCreatedAt = ts
					if id, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
						afterTweetID = id
						hasCursor = true
					}
				}
			}
		}
	}

	// --- Step 3: Execute query ---
	query, params := buildTimelineQuery(followingIDs, limit, hasCursor, afterCreatedAt, afterTweetID)

	rows, err := r.db.Query(ctx, query, params...)
	if err != nil {
		return nil, apperror.DB("timeline.tweets.query", err)
	}
	defer rows.Close()

	// --- Step 4: Parse rows ---
	tweets, lastCreatedAt, lastTweetID := parseTweetRows(rows, limit)

	// --- Step 5: Build next cursor and has_next_page ---
	var nextCursorPtr *string
	hasNextPage := false

	if lastCreatedAt != nil && lastTweetID != nil {
		// check for have more tweet
		hasMore, err := r.hasMoreTweets(ctx, followingIDs, *lastCreatedAt, *lastTweetID)
		if err != nil {
			return nil, apperror.DB("timeline.check_more_tweets", err)
		}

		hasNextPage = hasMore

		if hasNextPage {
			nextCursor := encodeCursor(*lastCreatedAt, *lastTweetID)
			nextCursorPtr = &nextCursor
		}
	}

	return &timelinedto.TimelineResponse{
		Tweets:      tweets,
		NextCursor:  nextCursorPtr,
		HasNextPage: hasNextPage,
		Limit:       limit,
	}, nil
}

// Helper function to check if there are more tweets
func (r *TimeLineRepo) hasMoreTweets(ctx context.Context, followingIDs []int64, lastCreatedAt time.Time, lastTweetID int64) (bool, error) {
	const query = `
SELECT EXISTS (
	SELECT 1 
	FROM tweets 
	WHERE user_id = ANY($1) 
	AND (created_at, id) < ($2, $3)
)`

	var exists bool
	err := r.db.QueryRow(ctx, query, followingIDs, lastCreatedAt, lastTweetID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Helper function for getting following IDs
func (r *TimeLineRepo) getFollowingIDs(ctx context.Context, userID int64) ([]int64, *apperror.AppError) {
	cacheKey := fmt.Sprintf("user_%d_followings", userID)

	// Try Redis first
	if b, err := redisclient.Rdb.Get(ctx, cacheKey).Bytes(); err == nil {
		var followingIDs []int64
		if json.Unmarshal(b, &followingIDs) == nil {
			return followingIDs, nil
		}
	}

	// Fallback to database
	const q = `SELECT following_id FROM follows WHERE follower_id = $1`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, apperror.DB("timeline.followings.query", err)
	}
	defer rows.Close()

	var followingIDs []int64
	for rows.Next() {
		var fid int64
		if err := rows.Scan(&fid); err != nil {
			return nil, apperror.DB("timeline.followings.scan", err)
		}
		followingIDs = append(followingIDs, fid)
	}

	// check error after Iteration
	if err := rows.Err(); err != nil {
		return nil, apperror.DB("timeline.followings.rows", err)
	}

	// Cache results
	if len(followingIDs) > 0 {
		if data, err := json.Marshal(followingIDs); err == nil {
			redisclient.Rdb.Set(ctx, cacheKey, data, 30*time.Minute)
		}
	}

	return followingIDs, nil
}

// Helper function for encoding cursor
func encodeCursor(createdAt time.Time, id int64) string {
	cursorValue := fmt.Sprintf("%s|%d", createdAt.Format(time.RFC3339Nano), id)
	return base64.URLEncoding.EncodeToString([]byte(cursorValue))
}

// Helper function for building query
func buildTimelineQuery(followingIDs []int64, limit int, hasCursor bool, afterCreatedAt time.Time, afterTweetID int64) (string, []interface{}) {
	baseQuery := `
SELECT
  t.id,
  t.content,
  t.like_count,
  t.dislike_count,
  t.reply_count,
  t.bookmark_count,
  t.image_url,
  t.created_at,
  u.username,
  u.avatar_url,
  COALESCE(ARRAY_AGG(DISTINCT tg.name) FILTER (WHERE tg.name IS NOT NULL), '{}') AS tags
FROM tweets t
JOIN users u ON u.id = t.user_id
LEFT JOIN tweet_tags tt ON tt.tweet_id = t.id
LEFT JOIN tags tg ON tg.id = tt.tag_id
WHERE t.user_id = ANY($1)
`

	params := []interface{}{followingIDs}
	paramCount := 1

	if hasCursor {
		baseQuery += ` AND (t.created_at, t.id) < ($2, $3)`
		params = append(params, afterCreatedAt, afterTweetID)
		paramCount += 2
	}

	baseQuery += `
GROUP BY t.id, u.username, u.avatar_url
ORDER BY t.created_at DESC, t.id DESC
LIMIT $` + strconv.Itoa(paramCount+1)

	params = append(params, limit)
	return baseQuery, params
}

// Helper function for parsing rows
func parseTweetRows(rows pgx.Rows, limit int) ([]timelinedto.Tweet, *time.Time, *int64) {
	tweets := make([]timelinedto.Tweet, 0, limit)
	var lastCreatedAt *time.Time
	var lastTweetID *int64

	for rows.Next() {
		var tweet timelinedto.Tweet
		var createdAt time.Time
		var tags []string

		if err := rows.Scan(
			&tweet.ID,
			&tweet.Content,
			&tweet.LikeCount,
			&tweet.DislikeCount,
			&tweet.ReplyCount,
			&tweet.BookmarkCount,
			&tweet.ImageURL,
			&createdAt,
			&tweet.Author.Username,
			&tweet.Author.AvatarURL,
			&tags,
		); err != nil {
			// Log error but continue
			continue
		}

		tweet.CreatedAt = createdAt.Format(time.RFC3339Nano)
		tweet.Tags = tags
		tweets = append(tweets, tweet)

		// Store last values for cursor
		lastCreatedAt = &createdAt
		lastID := tweet.ID
		lastTweetID = &lastID
	}

	return tweets, lastCreatedAt, lastTweetID
}
