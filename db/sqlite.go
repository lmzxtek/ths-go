package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite 驱动（匿名导入）
)

func Connect(dbfile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}
	return db, nil
}
func Close(db *sql.DB) error {
	return db.Close()
}
