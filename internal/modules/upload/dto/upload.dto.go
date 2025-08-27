package uploaddto

type UploadResponse struct {
	ID       int64  `json:"id"`
	URL      string `json:"url"`
	Hash     string `json:"hash"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
	FileName string `json:"file_name"`
}
