package userdtos

type GetUserByUsernameRes struct {
	Username       string   `json:"username"`
	Bio            *string  `json:"bio"`
	AvatarURL      *string  `json:"avatar_url"`
	FollowerCount  int64    `json:"follower_count"`
	FollowingCount int64    `json:"following_count"`
	Tweets         []Tweets `json:"tweets"`
}

type Tweets struct {
	ID            int64    `json:"id"`
	Content       string   `json:"content"`
	Tags          []string `json:"tags"`
	LikeCount     int64    `json:"like_count"`
	DislikeCount  int64    `json:"dislike_count"`
	ReplyCount    int64    `json:"reply_count"`
	BookMarkCount int64    `json:"bookmark_count"`
	CreatedAt     string   `json:"created_at"`
}
