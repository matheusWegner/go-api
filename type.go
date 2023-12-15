package main

import "math/rand"

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

func newUser(userName, email string) *User {
	return &User{
		ID:       rand.Intn(10000),
		UserName: userName,
		Email:    email,
	}
}
