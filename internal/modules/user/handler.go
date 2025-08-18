package user

import (
	"net/http"
	"strconv"
	userdtos "twitter_clone/internal/modules/user/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewUserHandler(service Service) *Handler {
	return &Handler{service: service}
}

// User godoc
// @Summary      Update user profile
// @Description  Update user profile. At least one of `username` or `bio` must be provided. Both fields are optional.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        data body userdtos.UpdateProfileReq true "Update profile payload (username and bio are optional, at least one required)"
// @Success      200 {object} userdtos.UpdateProfileRes "Updated user profile"
// @Failure      400 {object} apperror.AppError "Validation failed, e.g. no fields provided or invalid values"
// @Failure      500 {object} apperror.AppError "Internal server error"
// @Router       /users/update-profile [put]
func (h *Handler) UpdateProfile(c echo.Context) error {
	userID := c.Get("userID").(int64)

	var req userdtos.UpdateProfileReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apperror.Validation("Invalid request body", nil, err))
	}

	userUpdated, appErr := h.service.UpdateProfile(userID, req)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusCreated, userUpdated)
}

// User godoc
// @Summary      user profile
// @Description  get  user profile
// @Tags         User
// @Produce      json
// @Success      200 {object} userdtos.UserGetProfileRes
// @Failure      400 {object} apperror.AppError
// @Router       /users/get-profile [get]
func (h *Handler) GetProfile(c echo.Context) error {
	userID := c.Get("userID").(int64)

	userProfile, appErr := h.service.GetProfile(userID)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusCreated, userProfile)
}

// User godoc
// @Summary      Search users by username
// @Description  Search users by their username. The `limit` parameter is optional and defaults to 10 if not provided.
// @Tags         User
// @Produce      json
// @Param        username query string true  "Username prefix to search for"
// @Param        limit    query int    false "Maximum number of users to return (default 10)"
// @Success      200 {array} userdtos.SearchUsersByUsernameRes
// @Failure      400 {object} apperror.AppError "Validation failed"
// @Failure      500 {object} apperror.AppError "Internal server error"
// @Router       /users/search-by-user-name [get]
func (h *Handler) SearchUsersByUserName(c echo.Context) error {
	var validationErrors []map[string]string
	username := c.QueryParam("username")
	if username == "" {
		validationErrors = append(validationErrors, map[string]string{
			"error": "username query parameter is required",
		})

		return c.JSON(http.StatusBadRequest, apperror.Validation("Validation failed", validationErrors, nil))
	}

	maxLimit := 50
	limit := 10
	if param := c.QueryParam("limit"); param != "" {
		if num, err := strconv.Atoi(param); err != nil || num <= 0 {
			validationErrors = append(validationErrors, map[string]string{
				"error": "limit must be a positive integer",
			})

			return c.JSON(http.StatusBadRequest, apperror.Validation("Validation failed", validationErrors, nil))
		} else {
			if num > maxLimit {
				limit = maxLimit
			} else {
				limit = num
			}
		}
	}

	users, Serror := h.service.SearchUsersByUsername(username, limit)
	if Serror != nil {
		return c.JSON(Serror.StatusCode, Serror)
	}

	return c.JSON(http.StatusOK, users)
}
