package app

import (
	internalMiddleware "twitter_clone/internal/middleware"
	"twitter_clone/internal/modules/auth"
	"twitter_clone/internal/modules/tweet"
	"twitter_clone/internal/modules/user"

	"github.com/labstack/echo/v4/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Create auth Dependency
	authRepo := auth.NewAuthRepository(db)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)

	// Create user Dependency
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	// Create tweet Dependency
	tweetRepo := tweet.NewRepository(db)
	tweetService := tweet.NewTweetService(tweetRepo)
	tweetHandler := tweet.NewTweetHandler(tweetService)

	// Routs
	e.POST("/signup", authHandler.SignUp)
	e.POST("/login", authHandler.Login)

	// route need protection and token
	authGroup := e.Group("")
	authGroup.Use(internalMiddleware.JWTauthentication)

	// user
	authGroup.PUT("users/update-profile", userHandler.UpdateProfile)
	authGroup.GET("users/get-profile", userHandler.GetProfile)
	authGroup.GET("users/search-by-user-name", userHandler.SearchUsersByUserName)

	// tweet
	authGroup.POST("tweets/create-new-tweet", tweetHandler.CreateNewTweet)

	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
