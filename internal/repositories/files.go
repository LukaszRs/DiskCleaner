package repositories

import (
	"database/sql"
	"duplicates/internal/services"
)

func Create(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS files (
        id INTEGER PRIMARY KEY,
        path TEXT NOT NULL,
		filename TEXT NOT NULL,
		size INTEGER NOT NULL,
		created_at DATETIME,
		modified_at DATETIME,
		checksum TEXT NOT NULL,
		hash TEXT NOT NULL
    )`)
	if err != nil {
		panic(err)
	}
}

func InsertFile(db *sql.DB, f services.FileDef) {
	_, err := db.Exec(`INSERT INTO files (path, filename, size, created_at, modified_at, checksum, hash) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		f.Path, f.Filename, f.Size, f.CreatedAt, f.ModifiedAt, f.Checksum, f.Hash)
	if err != nil {
		panic(err)
	}
}
