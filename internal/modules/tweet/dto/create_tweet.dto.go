package tweetdtos

type CreateTweetReq struct {
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type CreateTweetRes struct {
	ID      int64    `json:"id"`
	UserID  int64    `json:"user_id"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}
