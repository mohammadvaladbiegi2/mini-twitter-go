package avatar

import (
	"net/http"
	"twitter_clone/internal/pkg/apperror"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

// Upload godoc
// @Summary      Upload a user avatar
// @Description  Upload a image . If the same file (by SHA256 hash) was uploaded before,
// @Description  the existing file metadata & URL will be returned (no duplicate storage).
// @Tags         Upload
// @Accept       multipart/form-data
// @Produce      json
// @Param        file     formData file true  "File to upload"
// @Success      200 {object} uploaddto.UploadResponse "Upload succeeded or existing file returned"
// @Failure      400 {object} apperror.AppError "Validation error (missing file or invalid user_id)"
// @Router       /users/me/avatar [post]
func (h *Handler) UploadAvatar(c echo.Context) error {
	userID := c.Get("userID").(int64)
	fh, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, apperror.Validation("file is required", nil, err))
	}

	res, appErr := h.svc.UploadAvatar(userID, fh)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusOK, res)
}
