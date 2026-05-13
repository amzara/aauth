package service

import (
	"aauth/internal/auth"
	"aauth/internal/db"
	"aauth/internal/session"
	"context"
	"errors"
	"fmt"
)

type AuthService struct {
	Queries *db.Queries
	Store   *session.Store
}

func NewAuthService(queries *db.Queries, s *session.Store) *AuthService {
	return &AuthService{Queries: queries, Store: s}
}

var ErrUserExists = errors.New("Username already exist") // sentinel error for the business logic errors
var ErrWrongPw = errors.New("Invalid password")

func (s *AuthService) Register(ctx context.Context, username string, password string) error {

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

func (s *AuthService) Login(ctx context.Context, username string, password string) (string, error) {

	var cred db.Cred
	cred, err := s.Queries.GetUserByUsername(ctx, username) //this is retarded, just get pw
	if err != nil {

		return ("Failed getting username and password from DB"), err
	}

	compare := auth.CheckPassword(password, cred.Password)

	if compare == false {
		return "invalid pw", ErrWrongPw

	}
	sessionToken, err := s.Store.Create(ctx, username, nil)

	if err != nil {
		return "failed to create session token", err
	}

	return sessionToken, nil

}
