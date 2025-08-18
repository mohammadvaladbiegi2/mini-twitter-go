package tweet

type Service interface{}

type TweetService struct {
	Repo *Repository
}

func NewTweetService(repo *Repository) *TweetService {
	return &TweetService{Repo: repo}
}
