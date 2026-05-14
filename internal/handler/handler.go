package handler

import (
	"aauth/internal/service"
	"aauth/internal/session"
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

type sessionToken struct {
	Token string `json:"token"`
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (s *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/register", s.Register)
	mux.HandleFunc("POST /api/login", s.Login)
	mux.HandleFunc("POST /api/auth", s.SessionCheck)
	mux.HandleFunc("POST /api/logout", s.Logout)
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

	sessionToken, err := s.authService.Login(r.Context(), creds.Username, creds.Password)
	if err != nil {
		if errors.Is(err, service.ErrWrongPw) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		} else {
			http.Error(w, "unexpected error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Login success, session token is %s", sessionToken),
	})
}

func (s *AuthHandler) SessionCheck(w http.ResponseWriter, r *http.Request) {
	var userToken sessionToken

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cant read request body", http.StatusBadRequest)
	}
	if err = json.Unmarshal(body, &userToken); err != nil {
		http.Error(w, "invalid json object", http.StatusBadRequest)
	}

	m := make(map[string]string)

	m, err = s.authService.Store.Get(r.Context(), userToken.Token)

	if err != nil {
		if errors.Is(err, session.ErrSessionNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Session does not exist",
			})
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(m)

}

func (s *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cant read request body", http.StatusBadRequest)
	}

	var token sessionToken

	if err := json.Unmarshal(body, &token); err != nil {
		http.Error(w, "cant read request body", http.StatusBadRequest)
	}

	err = s.authService.Logout(r.Context(), token.Token)
	if err != nil {
		if err == session.ErrSessionNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "session not found",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "internal server error",
		})
		return

	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Succesfully logged out",
	})

}
