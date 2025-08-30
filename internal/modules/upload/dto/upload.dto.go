package uploaddto

import (
	"mime/multipart"
	"time"
)

type UploadResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	FileName  string    `json:"file_name"`
	Hash      string    `json:"hash"`
	MimeType  string    `json:"mime_type"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	URL       string    `json:"url"`
}

type UploadTweetImageRes struct {
	UploadID int64  `json:"upload_id"`
	ImageURL string `json:"image_url"`
}

type UploadTweetImageReq struct {
	TweetID int64                 `form:"tweet_id" binding:"required"`
	File    *multipart.FileHeader `form:"file" binding:"required"`
}
