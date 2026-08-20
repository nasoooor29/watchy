package db

type Show struct {
	Id     int64  `db:"id"`
	Poster []byte `db:"poster"`
	Path   string `db:"path"`
	Title  string `db:"title"`
	MalId  int    `db:"mal_id"`
}

func (q *SQL) GetShow(id int64) (*Show, error) {
	show, err := Get[Show](q.db, "SELECT * FROM shows WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return show, nil
}

func (q *SQL) GetShowByPath(path string) (*Show, error) {
	return Get[Show](q.db, "SELECT * FROM shows WHERE path = ?", path)
}

func (q *SQL) GetAllShows() ([]Show, error) {
	return GetAll[Show](q.db, "SELECT * FROM shows")
}

func (q *SQL) CreateShow(poster []byte, path, title string, malId int) (*Show, error) {
	result, err := q.db.Exec(
		"INSERT INTO shows (poster, path, title, mal_id) VALUES (?, ?, ?, ?)",
		poster, path, title, malId,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return q.GetShow(id)
}
