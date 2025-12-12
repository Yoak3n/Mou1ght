package request

type UserLoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	UserLoginRequest
	Email string `json:"email"`
}
