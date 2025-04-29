package http

import (
	"api-gateway/internal/config"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func SignUp(req SignUpRequest) (map[string]interface{}, error) {
	return postJSON(config.AuthServiceURL+"/signup", req)
}

func ConfirmEmail(token string) (map[string]interface{}, error) {
	return getJSON(config.AuthServiceURL + "/confirm?token=" + token)
}

func Login(req LoginRequest) (map[string]interface{}, error) {
	return postJSON(config.AuthServiceURL+"/login", req)
}

func ChangePassword(req ChangePasswordRequest) (map[string]interface{}, error) {
	return postJSON(config.AuthServiceURL+"/change-password", req)
}
