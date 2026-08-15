package db

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
