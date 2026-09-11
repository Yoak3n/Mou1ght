package request

type UserLoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	UserLoginRequest
	Email string `json:"email"`
}

type UpdateUserProfileRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
