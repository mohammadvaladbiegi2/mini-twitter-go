package userdtos

type SearchUsersByUsernameRes struct {
	Username  string  `json:"username"`
	Bio       *string `json:"bio"`
	AvatarURL *string `json:"avatar_url"`
}
