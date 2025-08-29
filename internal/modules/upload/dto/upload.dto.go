package uploaddto

import "time"

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
