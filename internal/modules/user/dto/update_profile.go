package userdtos

type UpdateProfileReq struct {
	Username string  `json:"username,omitempty"`
	Bio      *string `json:"bio,omitempty"`
}

type UpdateProfileRes struct {
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
}
