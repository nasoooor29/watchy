package db

import (
	"backend/models"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"

	_ "modernc.org/sqlite"
)

type SQL struct {
	db *sql.DB
}

func GetDB() (*SQL, error) {
	env := models.GetEnv()
	rawDB, err := sql.Open("sqlite", env.DatabaseURL)
	if err != nil {
		return nil, err
	}
	s := &SQL{db: rawDB}

	err = s.migrate()
	if err != nil {
		slog.Error("", "err", err)
		return nil, err
	}

	return s, nil
}

func (s *SQL) Close() error {
	return s.db.Close()
}

func (s *SQL) migrate() error {
	env := models.GetEnv()
	// here we will loop over this dir and load all the migrations and execute them in order
	files, err := os.ReadDir(env.MigrationsDir)
	if err != nil {
		return err
	}
	slog.Info("migrating database", "dir", env.MigrationsDir, "count", len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// read the file
		content, err := os.ReadFile(filepath.Join(env.MigrationsDir, file.Name()))
		if err != nil {
			return err
		}
		slog.Info("executing migration", "file", file.Name())
		// execute the migration
		_, err = s.db.Exec(string(content))
		if err != nil {
			return err
		}
	}
	return nil
}

// if 2 issues or more came from any of the code below we will delete it and use plain sql queries
// current counter: 0

func Get[T any](db *sql.DB, query string, args ...any) (*T, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}

	result, err := scanRow[T](rows)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GetAll[T any](db *sql.DB, query string, args ...any) ([]T, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T

	for rows.Next() {
		result, err := scanRow[T](rows)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func scanRow[T any](rows *sql.Rows) (T, error) {
	var result T

	v := reflect.ValueOf(&result).Elem()
	t := v.Type()

	if v.Kind() != reflect.Struct {
		return result, fmt.Errorf("T must be a struct")
	}

	columns, err := rows.Columns()
	if err != nil {
		return result, err
	}

	// db column name -> struct field
	fields := make(map[string]int)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		name := field.Tag.Get("db")
		if name == "" {
			name = field.Name
		}

		fields[name] = i
	}

	dest := make([]any, len(columns))

	for i, column := range columns {
		fieldIndex, ok := fields[column]

		if !ok {
			// Ignore columns that aren't in the struct.
			var discard any
			dest[i] = &discard
			continue
		}

		dest[i] = v.Field(fieldIndex).Addr().Interface()
	}

	if err := rows.Scan(dest...); err != nil {
		return result, err
	}

	return result, nil
}
