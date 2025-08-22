package userdtos

type UserGetProfileRes struct {
	Username       string            `json:"username"`
	Email          string            `json:"email"`
	Bio            *string           `json:"bio"`
	AvatarURL      *string           `json:"avatar_url"`
	FollowerCount  int64             `json:"follower_count"`
	FollowingCount int64             `json:"following_count"`
	Tweets         []TweetForProfile `json:"tweets"`
}

type TweetForProfile struct {
	ID           int64    `json:"id"`
	Content      string   `json:"content"`
	Tags         []string `json:"tags"`
	LikeCount    int64    `json:"like_count"`
	DislikeCount int64    `json:"dislike_count"`
	ReplyCount   int64    `json:"reply_count"`
	CreatedAt    string   `json:"created_at"`
}
