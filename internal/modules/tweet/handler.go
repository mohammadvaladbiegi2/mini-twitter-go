package tweet

import (
	"net/http"
	tweetdtos "twitter_clone/internal/modules/tweet/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service Service
}

func NewTweetHandler(service Service) *Handler {
	return &Handler{Service: service}
}

func (h Handler) CreateNewTweet(c echo.Context) error {
	userID := c.Get("userID").(int64)

	var req tweetdtos.CreateTweetReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apperror.Validation("Invalid request body", nil, err))
	}

	tweet, appErr := h.Service.CreateTweet(userID, req)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusOK, tweet)
}
