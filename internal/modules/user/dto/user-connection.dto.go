package userdtos

type UsersFollower struct {
	Username  string  `json:"username"`
	Bio       *string `json:"bio"`
	AvatarURL *string `json:"avatar_url"`
}

type UsersFollowing struct {
	Username  string  `json:"username"`
	Bio       *string `json:"bio"`
	AvatarURL *string `json:"avatar_url"`
}

type UsersFollowersRes struct {
	Followers []UsersFollower `json:"followers"`
	Limit     int             `json:"limit"`
	Offset    int             `json:"offset"`
	Count     int             `json:"count"`
	Total     int64           `json:"total"`
	HasNext   bool            `json:"has_next"`
}

type UsersFollowingsRes struct {
	Followings []UsersFollowing `json:"followings"`
	Limit      int              `json:"limit"`
	Offset     int              `json:"offset"`
	Count      int              `json:"count"`
	Total      int64            `json:"total"`
	HasNext    bool             `json:"has_next"`
}
