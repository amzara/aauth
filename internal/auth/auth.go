package auth

import (
	"aauth/internal/db"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Queries *db.Queries
}

type Credentials struct {
	Username string `json:"username", db: "username"`
	Password string `json:"username", db: "password"`
}

func HashPassword(pw string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 14)

	return string(bytes), err

}

func CheckPassword(pw, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	return err == nil //bool return can do this

}
