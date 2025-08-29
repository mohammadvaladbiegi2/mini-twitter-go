package uploadavataruser

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetAvatarUploadID(userID int64) (int64, error)
	UpdateAvatar(userID, uploadID int64, uploadURL string) error
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetAvatarUploadID(userID int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var uploadID sql.NullInt64
	err := r.db.QueryRow(ctx, `SELECT avatar_upload_id FROM users WHERE id=$1`, userID).Scan(&uploadID)
	if err != nil {
		return 0, err
	}
	if uploadID.Valid {
		return uploadID.Int64, nil
	}
	return 0, nil
}

func (r *UserRepository) UpdateAvatar(userID, uploadID int64, uploadURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, `UPDATE users SET avatar_upload_id=$1, avatar_url=$2  WHERE id=$3`, uploadID, uploadURL, userID)
	return err
}
