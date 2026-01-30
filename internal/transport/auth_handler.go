package transport

import (
    "encoding/json"
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
    Token string      `json:"token"`
    User  interface{} `json:"user,omitempty"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    // Only accept POST requests
    if r.Method != http.MethodPost {
        respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
        return
    }

    // Parse request body
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    // Call service to register user
    user, err := h.authService.Register(req.Username, req.Password)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    // Generate token for the new user
    token, err := h.authService.Login(req.Username, req.Password)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
        return
    }

    // Send success response
    respondWithJSON(w, http.StatusCreated, AuthResponse{
        Token: token,
        User:  user,
    })
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
        return
    }

    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    // Call service to login
    token, err := h.authService.Login(req.Username, req.Password)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, AuthResponse{
        Token: token,
    })
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
    respondWithJSON(w, status, ErrorResponse{Error: message})
}