package handler

import (
	"auth-service/internal/middleware"
	"auth-service/internal/usecase"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type AuthHandler struct {
	UseCase usecase.AuthUseCase
}

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{UseCase: uc}
}

func (h *AuthHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/signup", h.SignUp).Methods(http.MethodPost)
	r.HandleFunc("/confirm", h.ConfirmEmail).Methods(http.MethodGet)
	r.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	r.Handle("/change-password", middleware.AuthMiddleware(h.ChangePassword)).Methods(http.MethodPost)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	token, err := h.UseCase.SignUp(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered. Please confirm via token.",
		"token":   token,
	})
	if err != nil {
		return
	}
}

func (h *AuthHandler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if err := h.UseCase.ConfirmEmail(r.Context(), token); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := json.NewEncoder(w).Encode(map[string]string{"message": "Email confirmed"})
	if err != nil {
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	_ = json.NewDecoder(r.Body).Decode(&creds)

	token, user, err := h.UseCase.Login(r.Context(), creds.Email, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "JWT",
		Value: "Bearer " + token,
	})
	err = json.NewEncoder(w).Encode(map[string]string{
		"email": user.Email,
		"name":  user.Name,
	})
	if err != nil {
		return
	}
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)

	var req struct {
		OldPassword     string `json:"oldPassword"`
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	if req.NewPassword != req.ConfirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	err := h.UseCase.ChangePassword(r.Context(), email, req.OldPassword, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"})
	if err != nil {
		return
	}
}
