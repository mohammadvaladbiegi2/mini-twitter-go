package user

import (
	"net/http"
	"strconv"
	useraction "twitter_clone/internal/modules/user/action"
	userconnection "twitter_clone/internal/modules/user/connection"
	userdtos "twitter_clone/internal/modules/user/dto"
	userprofile "twitter_clone/internal/modules/user/profile"
	usersearch "twitter_clone/internal/modules/user/search"
	"twitter_clone/internal/pkg/apperror"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	profileService    userprofile.Service
	searchService     usersearch.Service
	actionService     useraction.Service
	connectionService userconnection.Service
}

func NewUserHandler(profileService userprofile.Service, searchService usersearch.Service, actionService useraction.Service, connectionService userconnection.Service) *Handler {
	return &Handler{
		profileService:    profileService,
		searchService:     searchService,
		actionService:     actionService,
		connectionService: connectionService,
	}
}

type SimpleResMessage struct {
	Message string `json:"message"`
}

// User godoc
// @Summary      Update user profile
// @Description  Update user profile. At least one of `username` or `bio` must be provided. Both fields are optional.
// @Tags         User Profile
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

	userUpdated, appErr := h.profileService.UpdateProfile(userID, req)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusCreated, userUpdated)
}

// User godoc
// @Summary      user profile
// @Description  get  user profile
// @Tags         User Profile
// @Produce      json
// @Success      200 {object} userdtos.UserGetProfileRes
// @Failure      400 {object} apperror.AppError
// @Router       /users/get-profile [get]
func (h *Handler) GetProfile(c echo.Context) error {
	userID := c.Get("userID").(int64)

	userProfile, appErr := h.profileService.GetProfile(userID)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusOK, userProfile)
}

// User godoc
// @Summary      Search users by username
// @Description  Search users by their username. The `limit` parameter is optional and defaults to 10 if not provided.
// @Tags         User Search
// @Produce      json
// @Param        username query string true  "User name prefix to search for"
// @Param        limit    query int    false "Maximum number of users to return (default 10)"
// @Success      200 {array} userdtos.SearchUsersByUsernameRes
// @Failure      400 {object} apperror.AppError "Validation failed"
// @Failure      500 {object} apperror.AppError "Internal server error"
// @Router       /users/search-by-user-name/{username}  [get]
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

	users, Serror := h.searchService.SearchUsersByUsername(username, limit)
	if Serror != nil {
		return c.JSON(Serror.StatusCode, Serror)
	}

	return c.JSON(http.StatusOK, users)
}

// User godoc
// @Summary      get users by username
// @Description  Search users by their `username`
// @Tags         User Search
// @Produce      json
// @Param        username query string true  "User name prefix to search for"
// @Success      200 {array} userdtos.GetUserByUsernameRes
// @Failure      400 {object} apperror.AppError "Validation failed"
// @Failure      500 {object} apperror.AppError "Internal server error"
// @Router       /users/get-by-user-name/{username}  [get]
func (h *Handler) GetUserByUsername(c echo.Context) error {
	var validationErrors []map[string]string
	username := c.QueryParam("username")
	if username == "" {
		validationErrors = append(validationErrors, map[string]string{
			"error": "username query parameter is required",
		})

		return c.JSON(http.StatusBadRequest, apperror.Validation("Validation failed", validationErrors, nil))
	}

	userID := c.Get("userID").(int64)

	users, Serror := h.searchService.GetUserByUsername(username, userID)
	if Serror != nil {
		return c.JSON(Serror.StatusCode, Serror)
	}

	return c.JSON(http.StatusOK, users)
}

// @Summary      Follow a user
// @Tags         User Actions
// @Produce      json
// @Param        target_id query int true "Target user ID"
// @Success      200 {object} SimpleResMessage
// @Failure      400 {object} apperror.AppError
// @Router       /users/follow  [post]
func (h *Handler) Follow(c echo.Context) error {
	userID := c.Get("userID").(int64)
	targetID, err := strconv.ParseInt(c.QueryParam("target_id"), 10, 64)
	if err != nil || targetID <= 0 {
		return c.JSON(http.StatusBadRequest, apperror.Validation("invalid target_id", nil, err))
	}

	appErr := h.actionService.Follow(userID, targetID)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Followed successfully"})
}

// @Summary      Unfollow a user
// @Tags         User Actions
// @Produce      json
// @Param        target_id query int true "Target user ID"
// @Success      200 {object} SimpleResMessage
// @Failure      400 {object} apperror.AppError
// @Router       /users/unfollow [post]
func (h *Handler) Unfollow(c echo.Context) error {
	userID := c.Get("userID").(int64)
	targetID, err := strconv.ParseInt(c.QueryParam("target_id"), 10, 64)
	if err != nil || targetID <= 0 {
		return c.JSON(http.StatusBadRequest, apperror.Validation("invalid target_id", nil, err))
	}

	appErr := h.actionService.Unfollow(userID, targetID)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Unfollowed successfully"})
}

// GetFollowers godoc
// @Summary      Get followers of a user
// @Description  Returns a paginated list of followers for the authenticated user.
// @Tags         User Connection
// @Produce      json
// @Param        page query int false "Page number (default 1)"
// @Param        page_size query int false "Number of items per page (max 50, default 10)"
// @Success      200 {object} userdtos.UsersFollowersRes
// @Failure      400 {object} apperror.AppError "Validation error"
// @Failure      500 {object} apperror.AppError "Database or server error"
// @Router       /users/followers [get]
func (h *Handler) GetFollowers(c echo.Context) error {
	var validationErrors []map[string]string

	userID := c.Get("userID").(int64)

	page := 1
	if param := c.QueryParam("page"); param != "" {
		if num, err := strconv.Atoi(param); err != nil || num <= 0 {
			validationErrors = append(validationErrors, map[string]string{
				"error": "page must be a positive integer",
			})
			return c.JSON(http.StatusBadRequest, apperror.Validation("Validation failed", validationErrors, nil))
		} else {
			page = num
		}
	}

	pageSize := 10
	maxPageSize := 50
	if param := c.QueryParam("page_size"); param != "" {
		if num, err := strconv.Atoi(param); err != nil || num <= 0 {
			validationErrors = append(validationErrors, map[string]string{
				"error": "page_size must be a positive integer",
			})
			return c.JSON(http.StatusBadRequest, apperror.Validation("Validation failed", validationErrors, nil))
		} else if num > maxPageSize {
			pageSize = maxPageSize
		} else {
			pageSize = num
		}
	}

	result, appErr := h.connectionService.GetFollowers(userID, page, pageSize)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	response := userdtos.UsersFollowersRes{
		Followers: result.Data,
		Limit:     result.PageSize,
		Offset:    (result.CurrentPage - 1) * result.PageSize,
		Count:     len(result.Data),
		Total:     result.TotalCount,
		HasNext:   result.HasNext,
	}

	return c.JSON(http.StatusOK, response)
}

// GetFollowings godoc
// @Summary      Get followings of a user
// @Description  Returns a paginated list of followings for the authenticated user.
// @Tags         User Connection
// @Produce      json
// @Param        page query int false "Page number (default 1)"
// @Param        page_size query int false "Number of items per page (max 50, default 10)"
// @Success      200 {object} userdtos.UsersFollowingsRes
// @Failure      400 {object} apperror.AppError "Validation error"
// @Failure      500 {object} apperror.AppError "Database or server error"
// @Router       /users/followings [get]
func (h *Handler) GetFollowings(c echo.Context) error {
	var validationErrors []map[string]string

	userID := c.Get("userID").(int64)

	page := 1
	if param := c.QueryParam("page"); param != "" {
		if num, err := strconv.Atoi(param); err != nil || num <= 0 {
			validationErrors = append(validationErrors, map[string]string{
				"error": "page must be a positive integer",
			})
			return c.JSON(http.StatusBadRequest, apperror.Validation("Validation failed", validationErrors, nil))
		} else {
			page = num
		}
	}

	pageSize := 10
	maxPageSize := 50
	if param := c.QueryParam("page_size"); param != "" {
		if num, err := strconv.Atoi(param); err != nil || num <= 0 {
			validationErrors = append(validationErrors, map[string]string{
				"error": "page_size must be a positive integer",
			})
			return c.JSON(http.StatusBadRequest, apperror.Validation("Validation failed", validationErrors, nil))
		} else if num > maxPageSize {
			pageSize = maxPageSize
		} else {
			pageSize = num
		}
	}

	result, appErr := h.connectionService.GetFollowings(userID, page, pageSize)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	response := userdtos.UsersFollowingsRes{
		Followings: result.Data,
		Limit:      result.PageSize,
		Offset:     (result.CurrentPage - 1) * result.PageSize,
		Count:      len(result.Data),
		Total:      result.TotalCount,
		HasNext:    result.HasNext,
	}

	return c.JSON(http.StatusOK, response)
}
