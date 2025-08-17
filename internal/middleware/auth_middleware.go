package middleware

import (
	"net/http"
	"strings"
	"twitter_clone/internal/pkg/apperror"
	"twitter_clone/internal/pkg/jwt"

	"github.com/labstack/echo/v4"
)

func JWTauthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, apperror.UnauthorizedErr("token is requier", nil))
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, apperror.UnauthorizedErr("token is requier", nil))
		}
		tokenString := parts[1]

		claims, appErr := jwt.VerifyToken(tokenString)
		if appErr != nil {
			return c.JSON(appErr.StatusCode, appErr)
		}

		c.Set("userID", claims.ID)
		c.Set("userName", claims.UserName)

		return next(c)
	}
}
