package db

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (q *SQL) UpsertUser(id, name, email, password string) (*User, error) {
	if id == "" {
		id = uuid.NewString()
	}

	passwordHash := ""
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		passwordHash = string(hash)
	}

	return Get[User](
		q.db,
		`
		INSERT INTO users (id, name, email, password)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			email = excluded.email,
			password = CASE WHEN excluded.password != '' THEN excluded.password ELSE users.password END
		RETURNING *
		`,
		id, name, email, passwordHash,
	)
}
