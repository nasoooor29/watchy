package db

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const SESSION_DURATION = 7 * 24 * time.Hour
const FOREVER_DURATION = 30 * 24 * time.Hour

type Cookie struct {
	UUID      string `db:"uuid"`
	Userid    string `db:"userid"`
	ExpiresAt string `db:"expires_at"`
}

// ReadCookieByUUID
func (q *SQL) ReadCookieByUUID(uuid string) (*Cookie, error) {
	cookie, err := Get[Cookie](q.db, "SELECT * FROM cookies WHERE uuid = ?", uuid)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

func (q *SQL) CreateCookie(userId string, forever bool) (*Cookie, error) {
	var duration time.Duration
	if forever {
		duration = FOREVER_DURATION
	} else {
		duration = SESSION_DURATION
	}
	id := uuid.NewString()
	_, err := q.db.Exec(
		"INSERT INTO cookies (uuid, userid, expires_at) VALUES (?, ?, ?)",
		id, userId, time.Now().Add(duration).UTC().Format(time.RFC3339),
	)
	cc, err := q.ReadCookieByUUID(id)
	if err != nil {
		slog.Error("", "err", err)
		return nil, err
	}
	return cc, err
}

func (q *SQL) DeleteCookie(uuid string) error {
	_, err := q.db.Exec("DELETE FROM cookies WHERE uuid = ?", uuid)
	return err
}

func (q *SQL) DeleteExpiredCookies() error {
	_, err := q.db.Exec("DELETE FROM cookies WHERE expires_at <= ?", time.Now().UTC().Format(time.RFC3339))
	return err
}
