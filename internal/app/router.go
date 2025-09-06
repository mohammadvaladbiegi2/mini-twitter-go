package app

import (
	"log"
	internalMiddleware "twitter_clone/internal/middleware"
	"twitter_clone/internal/modules/auth"
	"twitter_clone/internal/modules/timeline"
	"twitter_clone/internal/modules/tweet"
	tweetaction "twitter_clone/internal/modules/tweet/action"
	tweetbookmark "twitter_clone/internal/modules/tweet/bookmark"
	tweetreply "twitter_clone/internal/modules/tweet/reply"
	"twitter_clone/internal/modules/upload"
	"twitter_clone/internal/modules/upload/avatar"
	"twitter_clone/internal/modules/upload/tweetmedia"
	uploadavataruser "twitter_clone/internal/modules/upload/user"
	"twitter_clone/internal/modules/user"
	useraction "twitter_clone/internal/modules/user/action"
	userconnection "twitter_clone/internal/modules/user/connection"
	userprofile "twitter_clone/internal/modules/user/profile"
	usersearch "twitter_clone/internal/modules/user/search"
	"twitter_clone/internal/pkg/apperror"
	"twitter_clone/internal/pkg/redisclient"

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

	userProfileRepo := userprofile.NewProfileRepository(db)
	userProfileService := userprofile.NewProfileService(userProfileRepo)

	userActionRepo := useraction.NewUserActionRepository(db)
	userActionService := useraction.NewUserActionService(userActionRepo)

	userSearchRepo := usersearch.NewSearchRepository(db)
	userSearchService := usersearch.NewSearchService(userSearchRepo)

	userConnectionRepo := userconnection.NewUserConnectionRepository(db)
	userConnectionService := userconnection.NewUserConnectionService(userConnectionRepo)

	userHandler := user.NewUserHandler(userProfileService, userSearchService, userActionService, userConnectionService)

	// Create tweet Dependency
	tweetRepo := tweet.NewRepository(db)
	tweetService := tweet.NewTweetService(tweetRepo)
	tweetActionRepo := tweetaction.NewUserActionRepository(db)
	tweetActionService := tweetaction.NewTweetActionService(tweetActionRepo)
	tweetReplyRepo := tweetreply.NewReplyRepository(db)
	tweetReplyService := tweetreply.NewReplyService(tweetReplyRepo)
	tweetBookMarkRepo := tweetbookmark.NewBookMarkRepository(db)
	tweetBookMarkService := tweetbookmark.NewBookMarkService(tweetBookMarkRepo)

	tweetHandler := tweet.NewTweetHandler(tweetService, tweetActionService, tweetReplyService, tweetBookMarkService)

	// init storage
	minioStorage, err := upload.NewMinioStorageFromEnv()
	if err != nil {
		log.Fatal(apperror.Server("cant init the minio storage", err))
	}

	// init redis
	if err := redisclient.Init(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	// repos
	uploadRepo := upload.NewUploadRepository(db)
	userRepo := uploadavataruser.NewUserRepository(db)
	tweetImageRepo := tweetmedia.NewTweetRepository(db)

	// services
	uploadSvc := upload.NewUploadService(uploadRepo, minioStorage)
	avatarSvc := avatar.NewService(uploadSvc, userRepo, uploadRepo, minioStorage)
	tweetImageService := tweetmedia.NewService(uploadSvc, tweetImageRepo, uploadRepo, minioStorage)

	// handlers
	avatarHandler := avatar.NewHandler(avatarSvc)
	tweetImageHandler := tweetmedia.NewHandler(tweetImageService)

	//  TimeLine
	timelineRepo := timeline.NewTimeLineRepo(db)
	timelineService := timeline.NewTimeLineService(timelineRepo)
	timelineHandler := timeline.NewTimeLineHandler(timelineService)

	// register route (protected)
	authGroup := e.Group("")
	authGroup.Use(internalMiddleware.JWTauthentication)

	// Routs Auth
	e.POST("/signup", authHandler.SignUp)
	e.POST("/login", authHandler.Login)

	// upload
	authGroup.POST("/users/me/avatar", avatarHandler.UploadAvatar)
	authGroup.POST("/tweets/:tweet_id/upload-image", tweetImageHandler.UploadTweetImage)

	// user
	authGroup.PUT("users/update-profile", userHandler.UpdateProfile)
	authGroup.GET("users/get-profile", userHandler.GetProfile)
	authGroup.GET("users/search-by-user-name", userHandler.SearchUsersByUserName)
	authGroup.GET("users/get-by-user-name", userHandler.GetUserByUsername)
	authGroup.POST("users/follow", userHandler.Follow)
	authGroup.POST("users/unfollow", userHandler.Unfollow)
	authGroup.GET("users/followers", userHandler.GetFollowers)
	authGroup.GET("users/followings", userHandler.GetFollowings)

	// tweet
	authGroup.POST("tweets/create-new-tweet", tweetHandler.CreateNewTweet)
	authGroup.POST("tweets/:tweet_id/like", tweetHandler.Like)
	authGroup.POST("tweets/:tweet_id/dislike", tweetHandler.Dislike)
	authGroup.DELETE("tweets/:tweet_id/reaction", tweetHandler.RemoveReaction)
	authGroup.POST("tweets/:tweet_id/reply", tweetHandler.CreateReply)
	authGroup.GET("tweets/:tweet_id/replies", tweetHandler.GetReplies)
	authGroup.POST("tweets/:tweet_id/bookmark", tweetHandler.Bookmark)
	authGroup.DELETE("tweets/:tweet_id/bookmark", tweetHandler.Unbookmark)
	authGroup.GET("tweets/bookmarks", tweetHandler.ListBookmarks)

	// timeline

	authGroup.GET("timeline/my_time_line", timelineHandler.MyTimeLine)

	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
