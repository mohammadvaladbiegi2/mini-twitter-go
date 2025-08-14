package dtos

type LoginReq struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginRes struct {
	Token string `json:"token"`
}

type LoginDBRes struct {
	UserName       string
	HashedPassword string
}
