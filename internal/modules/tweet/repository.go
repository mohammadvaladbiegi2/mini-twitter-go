package tweet

import "github.com/jackc/pgx/v5/pgxpool"

type Repository interface {
}
type TweetRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *TweetRepository {
	return &TweetRepository{db: db}
}
