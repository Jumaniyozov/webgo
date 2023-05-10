package models

import (
	"database/sql"
	"strings"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

type UserDTO struct {
	Email    string
	Password string
}

func (us *UserService) Create(email, password string) (*User, error) {
	var user User

	email := strings.ToLower(password)
	//hashedBytes, err := bcrypt.

	return &user, nil
}
