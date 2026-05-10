package service

import (
	"aauth/internal/auth"
	"aauth/internal/db"
	"context"
	"errors"
	"fmt"
)

type AuthService struct {
	Queries *db.Queries
}

func NewAuthService(queries *db.Queries) *AuthService {
	return &AuthService{Queries: queries}
}

var ErrUserExists = errors.New("Username already exist") // sentinel error for the business logic errors

func (s *AuthService) Register(ctx context.Context, username string, password string) error {

	fmt.Println("Received username : %v", username)
	fmt.Println("Received password : %v", password)

	exists, err := s.Queries.CheckUserExists(ctx, username)
	if err != nil {
		return fmt.Errorf("Failed to check user exist: %w", err)

	}

	if exists {
		return ErrUserExists
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("Failed to hash password %w", err)

	}

	if err := s.Queries.Register(ctx, db.RegisterParams{
		Username: username,
		Password: hashedPassword,
	}); err != nil {
		return fmt.Errorf("Failed to register user into DB %w", err)
	}

	return nil

}
