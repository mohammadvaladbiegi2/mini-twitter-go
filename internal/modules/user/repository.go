package user

import (
	"context"
	"time"
	userdtos "twitter_clone/internal/modules/user/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	UpdateProfile(userID int64, updateProfileReq userdtos.UpdateProfileReq) (userdtos.UpdateProfileRes, *apperror.AppError)
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db}
}

func (r UserRepository) UpdateProfile(userID int64, updateProfileReq userdtos.UpdateProfileReq) (userdtos.UpdateProfileRes, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE users
		SET username = $1,
		    bio = $2,
		    updated_at = NOW()
		WHERE id = $3
		RETURNING username, bio
	`

	var res userdtos.UpdateProfileRes
	err := r.db.QueryRow(ctx, query, updateProfileReq.Username, updateProfileReq.Bio, userID).
		Scan(&res.Username, &res.Bio)

	if err != nil {
		return userdtos.UpdateProfileRes{}, apperror.DB("failed to update profile", err)
	}

	return res, nil
}
