package upload

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"twitter_clone/internal/pkg/apperror"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage interface {
	Save(ctx context.Context, objectName string, data []byte, contentType string) (string, *apperror.AppError)
	PublicURL(objectName string) string
	PresignedURL(objectName string, expiry time.Duration) (string, *apperror.AppError)
	Delete(objectName string) *apperror.AppError
}

type MinioStorage struct {
	client       *minio.Client
	bucketName   string
	endpoint     string
	secure       bool
	usePresigned bool
}

func NewMinioStorageFromEnv() (*MinioStorage, *apperror.AppError) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET")
	secureStr := os.Getenv("MINIO_SECURE")
	usePresignedStr := os.Getenv("MINIO_USE_PRESIGNED")

	if endpoint == "" || bucket == "" {
		return nil, apperror.Server("MINIO_ENDPOINT and MINIO_BUCKET env vars required", nil)
	}
	if accessKey == "" || secretKey == "" {
		accessKey = os.Getenv("MINIO_ROOT_USER")
		secretKey = os.Getenv("MINIO_ROOT_PASSWORD")
	}
	if accessKey == "" || secretKey == "" {
		return nil, apperror.Server("MINIO_credentials not set in env", nil)
	}

	secure := false
	if secureStr != "" {
		s, _ := strconv.ParseBool(secureStr)
		secure = s
	}

	usePresigned := false
	if usePresignedStr != "" {
		p, _ := strconv.ParseBool(usePresignedStr)
		usePresigned = p
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, apperror.Server("failed to init minio client", err)
	}

	// ensure bucket exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, apperror.Server("failed to check bucket", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, apperror.Server("failed to create bucket", err)
		}

		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": "*",
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, bucket)

		if err := client.SetBucketPolicy(ctx, bucket, policy); err != nil {
			fmt.Printf("Warning: failed to set bucket policy: %v\n", err)
		}
	}

	return &MinioStorage{
		client:       client,
		bucketName:   bucket,
		endpoint:     endpoint,
		secure:       secure,
		usePresigned: usePresigned,
	}, nil
}

func (s *MinioStorage) Save(ctx context.Context, objectName string, data []byte, contentType string) (string, *apperror.AppError) {
	reader := bytes.NewReader(data)
	size := int64(len(data))

	_, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", apperror.Server("failed to upload file to storage", err)
	}

	if s.usePresigned {
		fmt.Println("DEBUG: Generating presigned URL...")
		presignedURL, appErr := s.PresignedURL(objectName, 24*time.Hour)
		if appErr != nil {
			return "", appErr
		}
		return presignedURL, nil
	}

	publicURL := s.PublicURL(objectName)
	return publicURL, nil
}

func (s *MinioStorage) PublicURL(objectName string) string {
	scheme := "http"
	if s.secure {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s/%s/%s", scheme, s.endpoint, s.bucketName, url.PathEscape(objectName))
}

func (s *MinioStorage) PresignedURL(objectName string, expiry time.Duration) (string, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", apperror.Server("failed to generate presigned URL", err)
	}

	return presignedURL.String(), nil
}

func (s *MinioStorage) Delete(objectName string) *apperror.AppError {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return apperror.Server("failed to delete object", err)
	}
	return nil
}
