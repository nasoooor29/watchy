package utils

import (
	"backend/db"
	"backend/models"
	"log/slog"
	"net/http"
	"time"
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
	expiresAt, err := time.Parse(time.RFC3339, data.ExpiresAt)
	if err != nil || !expiresAt.After(time.Now()) {
		return nil, models.ErrUnauthorized
	}

	u, err := q.GetUser(data.Userid)
	if err != nil {
		slog.Error("error happened", "err", err)
		return nil, err
	}
	return u, nil
}

func GenerateCookie(w http.ResponseWriter, q *db.SQL, userId string, forever bool) error {
	cookie, err := q.CreateCookie(userId, forever)
	if err != nil {
		slog.Error("", "err", err)
		return err
	}
	// insert into handler
	http.SetCookie(w, &http.Cookie{
		Name:     AUTH_COOKE_NAME,
		Value:    cookie.UUID,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}
