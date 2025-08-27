package upload

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	uploaddto "twitter_clone/internal/modules/upload/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/google/uuid"
)

type Service interface {
	UploadFile(userID int64, fileHeader *multipart.FileHeader) (*uploaddto.UploadResponse, *apperror.AppError)
}

type UploadService struct {
	Repo    Repository
	Storage Storage
}

func NewUploadService(repo Repository, storage Storage) *UploadService {
	return &UploadService{Repo: repo, Storage: storage}
}

func (s *UploadService) UploadFile(userID int64, fileHeader *multipart.FileHeader) (*uploaddto.UploadResponse, *apperror.AppError) {
	// open uploaded file
	f, err := fileHeader.Open()
	if err != nil {
		return nil, apperror.Validation("cannot open file", nil, err)
	}
	defer f.Close()

	// read bytes
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, apperror.Server("failed to read file", err)
	}

	// content type: prefer header, otherwise detect
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = httpDetectContentType(data)
	}

	// compute sha256 hash
	sum := sha256.Sum256(data)
	hashHex := hex.EncodeToString(sum[:])

	// check DB if already exists
	existingID, existingFile, existingMime, existingSize, appErr := s.Repo.FindByHash(hashHex)
	if appErr != nil {
		return nil, appErr
	}
	if existingID != 0 {
		// file already exists: construct public URL from storage and return
		url := s.Storage.PublicURL(existingFile)
		return &uploaddto.UploadResponse{
			ID:       existingID,
			URL:      url,
			Hash:     hashHex,
			MimeType: existingMime,
			Size:     existingSize,
			FileName: existingFile,
		}, nil
	}

	// generate safe object name
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		// fallback extension by mime
		ext = mimeExtFromContentType(contentType)
	}
	objectName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	// we can optionally prefix by user id & date: e.g. "user/123/2025/08/25/<object>"
	objectName = fmt.Sprintf("uploads/%d/%s", userID, objectName)

	// upload to storage
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	publicURL, appErr := s.Storage.Save(ctx, objectName, data, contentType)
	if appErr != nil {
		return nil, appErr
	}

	// save metadata in DB (SaveMetadata handles duplicate key race)
	id, appErr := s.Repo.SaveMetadata(userID, objectName, hashHex, contentType, int64(len(data)))
	if appErr != nil {
		// if metadata save failed but upload was successful, you may decide to delete object (optional).
		// For now we just return DB error (and object remains). Production: add cleanup or compensation.
		return nil, appErr
	}

	return &uploaddto.UploadResponse{
		ID:       id,
		URL:      publicURL,
		Hash:     hashHex,
		MimeType: contentType,
		Size:     int64(len(data)),
		FileName: objectName,
	}, nil
}

// helper - detect content-type like http.DetectContentType but independent here
func httpDetectContentType(data []byte) string {
	if len(data) < 512 {
		return http.DetectContentType(data)
	}
	return http.DetectContentType(data[:512])
}

// helper to map common mime to extension fallback
func mimeExtFromContentType(mime string) string {
	switch mime {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".bin"
	}
}
