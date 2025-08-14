package dtos

type UserSignUpReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignUpRes struct {
	Token string `json:"token"`
}
