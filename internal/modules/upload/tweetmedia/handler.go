package tweetmedia

import (
	"net/http"
	"strconv"
	uploaddto "twitter_clone/internal/modules/upload/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

// UploadTweetImage godoc
// @Summary      Upload an image for a tweet
// @Description  Uploads an image and associates it with the given tweet
// @Tags         Upload
// @Accept       mpfd
// @Produce      json
// @Param        tweet_id path int true "ID of the tweet"
// @Param        file formData file true "Image file to upload"
// @Success      200 {object} uploaddto.UploadTweetImageRes
// @Failure      400 {object} apperror.AppError
// @Failure      500 {object} apperror.AppError
// @Router       /tweets/{tweet_id}/upload-image [post]
func (h *Handler) UploadTweetImage(c echo.Context) error {
	userID := c.Get("userID").(int64)
	tweetIDStr := c.Param("tweet_id")
	tweetID, err := strconv.ParseInt(tweetIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, apperror.Validation("invalid tweet_id", nil, err))
	}

	file, err := c.FormFile("file")

	if err != nil || file == nil {
		return c.JSON(http.StatusBadRequest, apperror.Validation("file is required", nil, err))
	}

	uploadRes, appErr := h.svc.UploadTweetImage(userID, tweetID, file)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusOK, uploaddto.UploadTweetImageRes{
		UploadID: uploadRes.ID,
		ImageURL: uploadRes.URL,
	})
}
