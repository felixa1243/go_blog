package dto

type LoginResponse struct {
	Message string            `json:"message"`
	Token   string            `json:"token"`
	User    UserLoginResponse `json:"user"`
}
