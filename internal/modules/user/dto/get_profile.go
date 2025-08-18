package userdtos

type UserGetProfileRes struct {
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	Bio            *string `json:"bio"`
	AvatarURL      *string `json:"avatar_url"`
	FollowerCount  int64   `json:"follower_count"`
	FollowingCount int64   `json:"following_count"`
}
