package tweetdtos

type CreateReplyReq struct {
	Content string `json:"content"`
	// for futuer: ParentReplyID *int64 `json:"parent_reply_id"`
}

type Reply struct {
	ID        int64   `json:"id"`
	TweetID   int64   `json:"tweet_id"`
	UserID    int64   `json:"user_id"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url"`
	Content   string  `json:"content"`
	CreatedAt string  `json:"created_at"`
}

type CreateReplyRes struct {
	Reply Reply `json:"reply"`
}

type GetRepliesRes struct {
	Replies    []Reply `json:"replies"`
	Count      int     `json:"count"`
	NextCursor *int64  `json:"next_cursor"`
}
