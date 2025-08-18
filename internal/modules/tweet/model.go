package tweet

import "time"

type Tweet struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []Tag     `json:"tags,omitempty"`
	Likes     []Like    `json:"likes,omitempty"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Like struct {
	TweetID   int64     `json:"tweet_id"`
	UserID    int64     `json:"user_id"`
	IsLike    bool      `json:"is_like"`
	CreatedAt time.Time `json:"created_at"`
}
