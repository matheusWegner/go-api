package main

import (
	"time"
 )
type CreateUserRequest struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func newUser(email,userName string) *User {
	return &User{
		UserName: userName,
		Email:    email,
		CreatedAt: time.Now().UTC(),
	}
}
