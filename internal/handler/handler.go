package handler

import (
	"aauth/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
}

type signupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
} //should be in domain package

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (s *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/register", s.Register)
	mux.HandleFunc("POST /api/login", s.Login)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var creds signupRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &creds); err != nil {
		http.Error(w, "invalid json format", http.StatusBadRequest)
		return
	} //this part should filter the json

	if err := h.authService.Register(r.Context(), creds.Username, creds.Password); err != nil {
		log.Printf("REGISTER ERROR: %v", err) // ← ADD THIS

		if errors.Is(err, service.ErrUserExists) {
			http.Error(w, "username already exists", http.StatusConflict)
			return
		}
		fmt.Println("Oops triggered internal server error ")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Success register")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user created"})
}

//TODO: pw validation

func (s *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds signupRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &creds); err != nil {
		http.Error(w, "failed to unmarshal json body", http.StatusBadRequest)
		return
	}

	if err = s.authService.Login(r.Context(), creds.Username, creds.Password); err != nil {
		if errors.Is(err, service.ErrWrongPw) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		} else {
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login success"})
}
