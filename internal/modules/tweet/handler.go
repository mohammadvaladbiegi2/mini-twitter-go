package tweet

import (
	"net/http"
	"strconv"
	tweetaction "twitter_clone/internal/modules/tweet/action"
	tweetdtos "twitter_clone/internal/modules/tweet/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service       Service
	actionService tweetaction.Service
}

func NewTweetHandler(service Service, actionService tweetaction.Service) *Handler {
	return &Handler{
		service:       service,
		actionService: actionService,
	}
}

type SimpleResMessage struct {
	Message string `json:"message"`
}

// Tweet godoc
// @Summary      new Tweet
// @Description  Create new tweet, The `tags` field is optional
// @Tags         Tweet
// @Produce      json
// @Success      200 {object} tweetdtos.CreateTweetRes
// @Failure      400 {object} apperror.AppError
// @Router       /tweets/create-new-tweet [post]
func (h Handler) CreateNewTweet(c echo.Context) error {
	userID := c.Get("userID").(int64)

	var req tweetdtos.CreateTweetReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apperror.Validation("Invalid request body", nil, err))
	}

	tweet, appErr := h.service.CreateTweet(userID, req)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusOK, tweet)
}

// @Summary Like a tweet
// @Tags Tweet Action
// @Param tweet_id path int true "Tweet ID"
// @Success 200 {object} SimpleResMessage
// @Failure 400 {object} apperror.AppError
// @Router /tweets/{tweet_id}/like [post]
func (h *Handler) Like(c echo.Context) error {
	userID := c.Get("userID").(int64)
	tweetID, _ := strconv.ParseInt(c.Param("tweet_id"), 10, 64)

	if err := h.actionService.Like(userID, tweetID); err != nil {
		return c.JSON(err.StatusCode, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Like successfully"})
}

// @Summary Dislike a tweet
// @Tags Tweet Action
// @Param tweet_id path int true "Tweet ID"
// @Success 200 {object} SimpleResMessage
// @Failure 400 {object} apperror.AppError
// @Router /tweets/{tweet_id}/dislike [post]
func (h *Handler) Dislike(c echo.Context) error {
	userID := c.Get("userID").(int64)
	tweetID, _ := strconv.ParseInt(c.Param("tweet_id"), 10, 64)

	if err := h.actionService.Dislike(userID, tweetID); err != nil {
		return c.JSON(err.StatusCode, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Dislike successfully"})
}

// @Summary Remove reaction
// @Tags Tweet Action
// @Param tweet_id path int true "Tweet ID"
// @Success 200 {object} SimpleResMessage
// @Failure 400 {object} apperror.AppError
// @Router /tweets/{tweet_id}/reaction [delete]
func (h *Handler) RemoveReaction(c echo.Context) error {
	userID := c.Get("userID").(int64)
	tweetID, _ := strconv.ParseInt(c.Param("tweet_id"), 10, 64)

	if err := h.actionService.RemoveReaction(userID, tweetID); err != nil {
		return c.JSON(err.StatusCode, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Remove Reaction successfully"})
}
