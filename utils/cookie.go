package utils

import (
	"backend/db"
	"backend/models"
	"log/slog"
	"net/http"
)

const AUTH_COOKE_NAME = "auth_token"

func GetCookie(r *http.Request, name string) *http.Cookie {
	c, err := r.Cookie(name)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil
		}
		return nil
	}
	return c
}

func GetUserByCookie(r *http.Request, q *db.SQL) (*db.User, error) {
	c := GetCookie(r, AUTH_COOKE_NAME)
	if c == nil {
		return nil, models.ErrUnauthorized
	}
	data, err := q.ReadCookieByUUID(c.Value)
	if err != nil {
		return nil, models.ErrUnauthorized
	}

	u, err := q.GetUser(data.Userid)
	if err != nil {
		slog.Error("error happened", "err", err)
		return nil, err
	}
	return u, nil
}
