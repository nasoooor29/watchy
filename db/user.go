package db

import (
	"github.com/google/uuid"
)

type User struct {
	Id       string `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func (q *SQL) GetUser(id string) (*User, error) {
	user, err := Get[User](q.db, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (q *SQL) GetUserByEmail(email string) (*User, error) {
	return Get[User](q.db, "SELECT * FROM users WHERE email = ?", email)
}

func (q *SQL) CreateUser(name, email, password string) (*User, error) {
	uid := uuid.NewString()
	_, err := q.db.Exec(
		"INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)",
		uid, name, email, password,
	)
	if err != nil {
		return nil, err
	}
	return q.GetUser(uid)
}
