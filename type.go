package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	UserName string `json:"userName"`
	Token    string `json:"token"`
}

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID                int       `json:"id"`
	UserName          string    `json:"userName"`
	Email             string    `json:"email"`
	EncryptedPassword string    `json:"-"`
	CreatedAt         time.Time `json:"createdAt"`
}

func (u *User) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(pw)) == nil
}

func newUser(email, userName, password string) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		UserName:          userName,
		Email:             email,
		EncryptedPassword: string(encpw),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
