package tweetdtos

type BookmarkTweet struct {
	ID            int64    `json:"id"`
	Content       string   `json:"content"`
	Tags          []string `json:"tags"`
	LikeCount     int64    `json:"like_count"`
	DislikeCount  int64    `json:"dislike_count"`
	ReplyCount    int64    `json:"reply_count"`
	BookMarkCount int64    `json:"book_mark_count"`
	CreatedAt     string   `json:"created_at"`

	AuthorUsername string  `json:"author_username"`
	AuthorAvatar   *string `json:"author_avatar_url"`
}

type BookmarksListRes struct {
	Bookmarks  []BookmarkTweet `json:"bookmarks"`
	Limit      int             `json:"limit"`
	Count      int             `json:"count"`
	NextCursor *string         `json:"next_cursor"`
}
