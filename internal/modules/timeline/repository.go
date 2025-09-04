package timeline

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"twitter_clone/internal/pkg/apperror"
	"twitter_clone/internal/pkg/redisclient"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	AddDataToRedis(userID int64) (R, *apperror.AppError)
}

type TimeLineRepo struct {
	db *pgxpool.Pool
}

func NewTimeLineRepo(db *pgxpool.Pool) *TimeLineRepo {
	return &TimeLineRepo{db}
}

func (r TimeLineRepo) AddDataToRedis(userID int64) (R, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cacheKey := fmt.Sprintf("user_%d", userID)
	fmt.Println(cacheKey)

	// 1. check exist in redise
	val, err := redisclient.Rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedUser R
		if err := json.Unmarshal([]byte(val), &cachedUser); err == nil {
			fmt.Println("get from redis")
			return cachedUser, nil
		}
	}

	// 2. if not exist in redis read from database
	var user R
	query := "SELECT id, username, email, password_hash, bio, avatar_url FROM users WHERE id = $1"
	if err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.PasswordHash, &user.Bio, &user.AvatarURL,
	); err != nil {
		return R{}, apperror.DB("failed to find user", err)
	}

	// 3. save user in redis for hour
	jsonData, err := json.Marshal(user)
	if err == nil {
		_ = redisclient.Rdb.Set(ctx, cacheKey, jsonData, time.Hour).Err()
	}
	fmt.Println("get from DB")
	return user, nil
}
