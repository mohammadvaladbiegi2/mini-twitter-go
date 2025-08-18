package tweet

type Handler struct {
	Service Service
}

func NewTweetHandler(service Service) *Handler {
	return &Handler{Service: service}
}
