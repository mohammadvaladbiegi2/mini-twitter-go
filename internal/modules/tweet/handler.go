package tweet

import (
	"net/http"
	"strconv"
	tweetaction "twitter_clone/internal/modules/tweet/action"
	tweetbookmark "twitter_clone/internal/modules/tweet/bookmark"
	tweetdtos "twitter_clone/internal/modules/tweet/dto"
	tweetreply "twitter_clone/internal/modules/tweet/reply"
	"twitter_clone/internal/pkg/apperror"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service         Service
	actionService   tweetaction.Service
	replyService    tweetreply.Service
	bookmarkService tweetbookmark.Service
}

func NewTweetHandler(service Service, actionService tweetaction.Service, replyService tweetreply.Service, bookmarkService tweetbookmark.Service) *Handler {
	return &Handler{
		service:         service,
		actionService:   actionService,
		replyService:    replyService,
		bookmarkService: bookmarkService,
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

// Tweet godoc
// @Summary      Reply to a tweet
// @Description  Creates a reply for the given tweet_id
// @Tags         Tweet Action
// @Accept       json
// @Produce      json
// @Param        tweet_id path int true "Tweet ID"
// @Param        data body tweetdtos.CreateReplyReq true "reply payload"
// @Success      201 {object} tweetdtos.CreateReplyRes
// @Failure      400 {object} apperror.AppError
// @Failure      404 {object} apperror.AppError
// @Failure      500 {object} apperror.AppError
// @Router       /tweets/{tweet_id}/reply [post]
func (h *Handler) CreateReply(c echo.Context) error {
	userID := c.Get("userID").(int64)

	tweetID, err := strconv.ParseInt(c.Param("tweet_id"), 10, 64)
	if err != nil || tweetID <= 0 {
		return c.JSON(http.StatusBadRequest, apperror.Validation("invalid tweet_id", nil, err))
	}

	var req tweetdtos.CreateReplyReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apperror.Validation("invalid request body", nil, err))
	}

	res, appErr := h.replyService.CreateReply(userID, tweetID, req)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}
	return c.JSON(http.StatusCreated, res)
}

// Tweet godoc
// @Summary      List replies of a tweet
// @Description  Returns replies (newest first). Keyset pagination via optional cursor_id.
// @Tags         Tweet Action
// @Produce      json
// @Param        tweet_id  path int   true  "Tweet ID"
// @Param        limit     query int  false "Max items (<=100, default 20)"
// @Param        cursor_id query int  false "Fetch items with id < cursor_id"
// @Success      200 {object} tweetdtos.GetRepliesRes
// @Failure      400 {object} apperror.AppError
// @Failure      500 {object} apperror.AppError
// @Router       /tweets/{tweet_id}/replies [get]
func (h *Handler) GetReplies(c echo.Context) error {
	tweetID, err := strconv.ParseInt(c.Param("tweet_id"), 10, 64)
	if err != nil || tweetID <= 0 {
		return c.JSON(http.StatusBadRequest, apperror.Validation("invalid tweet_id", nil, err))
	}

	limit := 20
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	var cursorID *int64
	if v := c.QueryParam("cursor_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			cursorID = &n
		}
	}

	res, appErr := h.replyService.GetReplies(tweetID, limit, cursorID)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}
	return c.JSON(http.StatusOK, res)
}

// POST /tweets/:tweet_id/bookmark
// @Summary  Bookmark tweet
// @Tags     Tweet Action
// @Produce  json
// @Router   /tweets/{tweet_id}/bookmark [post]
func (h *Handler) Bookmark(c echo.Context) error {
	userID := c.Get("userID").(int64)

	tid, err := strconv.ParseInt(c.Param("tweet_id"), 10, 64)
	if err != nil || tid <= 0 {
		return c.JSON(http.StatusBadRequest, apperror.Validation("invalid tweet_id", nil, err))
	}

	if appErr := h.bookmarkService.Add(userID, tid); appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "bookmarked"})
}

// DELETE /tweets/:tweet_id/bookmark
// @Summary  Unbookmark tweet
// @Tags     Tweet Action
// @Produce  json
// @Router   /tweets/{tweet_id}/bookmark [delete]
func (h *Handler) Unbookmark(c echo.Context) error {
	userID := c.Get("userID").(int64)

	tid, err := strconv.ParseInt(c.Param("tweet_id"), 10, 64)
	if err != nil || tid <= 0 {
		return c.JSON(http.StatusBadRequest, apperror.Validation("invalid tweet_id", nil, err))
	}

	if appErr := h.bookmarkService.Remove(userID, tid); appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "unbookmarked"})
}

// GET /tweets/bookmarks?limit=&cursor=...
// @Summary  List my bookmarks (keyset pagination)
// @Tags     User Infos
// @Produce  json
// @Param    limit  query int    false "max 50 (default 20)"
// @Param    cursor query string false "opaque cursor"
// @Success  200    {object} tweetdtos.BookmarksListRes
// @Router   /tweets/bookmarks [get]
func (h *Handler) ListBookmarks(c echo.Context) error {
	userID := c.Get("userID").(int64)

	limit := 20
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 50 {
			limit = n
		}
	}
	cursor := c.QueryParam("cursor")

	res, appErr := h.bookmarkService.List(userID, limit, cursor)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}
	return c.JSON(http.StatusOK, res)
}
