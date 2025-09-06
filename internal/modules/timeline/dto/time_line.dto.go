package timelinedto

type Tweet struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	ImageURL  *string  `json:"image_url"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	Author    struct {
		Username  string  `json:"username"`
		AvatarURL *string `json:"avatar_url"`
	} `json:"author"`
	LikeCount     int64 `json:"like_count"`
	DislikeCount  int64 `json:"dislike_count"`
	ReplyCount    int64 `json:"reply_count"`
	BookmarkCount int64 `json:"bookmark_count"`
}

type TimelineResponse struct {
	Tweets      []Tweet `json:"tweets"`
	NextCursor  *string `json:"next_cursor"`
	HasNextPage bool    `json:"has_next_page"`
	Limit       int     `json:"limit"`
}

type Cursor struct {
	CreatedAt string `json:"created_at"`
	ID        int64  `json:"id"`
}
