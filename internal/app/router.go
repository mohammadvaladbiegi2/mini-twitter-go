package app

import (
	"twitter_clone/internal/modules/auth"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	// Create Dependency
	authRepo := auth.NewAuthRepository(db)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)

	// Routs
	e.POST("/signup", authHandler.SignUp)
	e.POST("/login", authHandler.Login)

	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
