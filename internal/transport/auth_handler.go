package transport

import (
	"net/http"

	"github.com/JuanTobonV/blog_app/internal/service"
)

type AuthHandler struct {
	authService service.IAuthService
}

func NewAuthHandler(authService service.IAuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User interface{} `json:"user,omitempty"` 
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *AuthHandler) Register (w http.ResponseWriter, r *http.Request) {
	
}