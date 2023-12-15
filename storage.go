package main

import "database/sql"

type Storage interface {
	CreateUser(*User) error
	DeleteUser(int) error
	UpdateUser(*User) error
	GetUserById(int) (*User, error)
}

type PostgressStore struct {
	db *sql.DB
}

func newPostgressStore() (*PostgressStore, error) {

}
