package upload

import (
	"net/http"
	"strconv"

	"twitter_clone/internal/pkg/apperror"

	"github.com/labstack/echo/v4"
)

// Handler handles upload endpoints
type Handler struct {
	Service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{Service: s}
}

// Upload godoc
// @Summary      Upload a file
// @Description  Upload a file (image/audio/etc.). If the same file (by SHA256 hash) was uploaded before,
// @Description  the existing file metadata & URL will be returned (no duplicate storage).
// @Tags         Upload
// @Accept       multipart/form-data
// @Produce      json
// @Param        file     formData file true  "File to upload"
// @Success      200 {object} uploaddto.UploadResponse "Upload succeeded or existing file returned"
// @Failure      400 {object} apperror.AppError "Validation error (missing file or invalid user_id)"
// @Router       /uploads [post]
func (h *Handler) Upload(c echo.Context) error {
	// get userID from context (JWT middleware should set it). Fallback: form value
	userID := c.Get("userID").(int64)

	// fallback to form param (only for manual/testing)
	if userID == 0 {
		userIDStr := c.FormValue("user_id")
		if userIDStr == "" {
			return c.JSON(http.StatusBadRequest, apperror.Validation("user_id is required (or set JWT token)", nil, nil))
		}
		id, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, apperror.Validation("invalid user_id", nil, err))
		}
		userID = id
	}

	// read uploaded file
	fh, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, apperror.Validation("file is required", nil, err))
	}

	resp, appErr := h.Service.UploadFile(userID, fh)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusOK, resp)
}
