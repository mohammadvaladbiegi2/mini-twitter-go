package upload

import (
	"context"
	"time"

	uploaddto "twitter_clone/internal/modules/upload/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindByHash(hash string) (id int64, fileName, mime string, size int64, appErr *apperror.AppError)
	SaveMetadata(userID int64, fileName, hash, mime string, size int64) (int64, *apperror.AppError)
	FindByID(id int64) (*uploaddto.UploadResponse, error)
	DeleteByID(id int64) error
}

type UploadRepository struct {
	db *pgxpool.Pool
}

func NewUploadRepository(db *pgxpool.Pool) *UploadRepository {
	return &UploadRepository{db}
}

func (r *UploadRepository) FindByHash(hash string) (int64, string, string, int64, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int64
	var fileName, mime string
	var size int64
	q := `SELECT id, file_name, mime_type, size FROM uploads WHERE hash=$1 LIMIT 1;`
	err := r.db.QueryRow(ctx, q, hash).Scan(&id, &fileName, &mime, &size)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, "", "", 0, nil
		}
		return 0, "", "", 0, apperror.DB("failed to query upload by hash", err)
	}
	return id, fileName, mime, size, nil
}

func (r *UploadRepository) SaveMetadata(userID int64, fileName, hash, mime string, size int64) (int64, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int64
	q := `
		INSERT INTO uploads (user_id, file_name, hash, mime_type, size, created_at)
		VALUES ($1,$2,$3,$4,$5,NOW())
		ON CONFLICT (hash) DO UPDATE SET file_name = uploads.file_name
		RETURNING id;
	`
	err := r.db.QueryRow(ctx, q, userID, fileName, hash, mime, size).Scan(&id)
	if err != nil {
		return 0, apperror.DB("failed to insert upload metadata", err)
	}
	return id, nil
}

func (r *UploadRepository) FindByID(id int64) (*uploaddto.UploadResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.db.QueryRow(ctx, `
		SELECT id, user_id, file_name, hash, mime_type, size, created_at
		FROM uploads WHERE id=$1
	`, id)

	var res uploaddto.UploadResponse
	err := row.Scan(&res.ID, &res.UserID, &res.FileName, &res.Hash, &res.MimeType, &res.Size, &res.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *UploadRepository) DeleteByID(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, `DELETE FROM uploads WHERE id=$1`, id)
	return err
}
