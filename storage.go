package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	createUser(*User) error
	deleteUser(int) error
	updateUser(*User) error
	getUsers() ([]*User, error)
	getUserById(int) (*User, error)
	getAccountByUserName(string) (*User, error)
}

type PostgresStore struct {
	db *sql.DB
}

func newPostgressStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=907010 sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) init() error {
	return s.createTable()
}

func (s *PostgresStore) createTable() error {
	query := `create table  if not exists users (
         id serial primary key,
		 email varchar(50),
		 user_name  varchar(50),
		 created_at  timestamp
 	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createUser(u *User) error {
	query := `
			insert into users 
			(email,user_name,encrypted_password,created_at) 
			values 
			($1,$2,$3,$4)
		`
	_, err := s.db.Query(
		query,
		u.Email,
		u.UserName,
		u.EncryptedPassword,
		u.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) updateUser(*User) error {
	return nil
}

func (s *PostgresStore) deleteUser(id int) error {
	_, err := s.db.Query("delete from users where id = $1", id)
	return err
}

func (s *PostgresStore) getAccountByUserName(userName string) (*User, error) {
	rows, err := s.db.Query("select * from account where user_name = $1", userName)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("account with userName [%d] not found", userName)
}

func (s *PostgresStore) getUserById(id int) (*User, error) {
	rows, err := s.db.Query("select * from users where id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoUser(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) getUsers() ([]*User, error) {
	rows, err := s.db.Query("select * from users")
	if err != nil {
		return nil, err
	}
	users := []*User{}
	for rows.Next() {
		user, err := scanIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func scanIntoUser(rows *sql.Rows) (*User, error) {
	user := new(User)
	err := rows.Scan(&user.ID, &user.Email, &user.UserName, &user.CreatedAt)

	return user, err
}
