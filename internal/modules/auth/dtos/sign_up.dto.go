package dtos

type SignUpReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRes struct {
	Token string `json:"token"`
}
