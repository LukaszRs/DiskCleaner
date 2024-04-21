package repositories

import (
	"database/sql"
	"duplicates/internal/models"
)

func Create(db *sql.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS files`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS files (
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

func InsertFile(db *sql.DB, f models.FileDef) {
	_, err := db.Exec(`INSERT INTO files (path, filename, size, created_at, modified_at, checksum, hash) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		f.Path, f.Filename, f.Size, f.CreatedAt, f.ModifiedAt, f.Checksum, f.Hash)
	if err != nil {
		panic(err)
	}
}

func GetPotentialDuplicates(db *sql.DB) []models.FileDef {
	rows, err := db.Query(`SELECT f1.id, f1.path, f1.filename, f1.size, f1.created_at, f1.modified_at, f1.checksum, f1.hash
		FROM files f1
		INNER JOIN files f2 ON
			(f1.id != f2.id AND f1.hash = "" AND f2.hash = "") AND
			(f1.checksum = f2.checksum OR f1.size = f2.size OR f1.filename = f2.filename)
		LIMIT 100;`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var files []models.FileDef
	for rows.Next() {
		var f models.FileDef
		err := rows.Scan(&f.ID, &f.Path, &f.Filename, &f.Size, &f.CreatedAt, &f.ModifiedAt, &f.Checksum, &f.Hash)
		if err != nil {
			panic(err)
		}
		files = append(files, f)
	}
	return files
}

func UpdateHash(db *sql.DB, id string, hash string) {
	_, err := db.Exec(`UPDATE files SET hash = ? WHERE id = ?`, hash, id)
	if err != nil {
		panic(err)
	}
}

func GetDuplicates(db *sql.DB) []models.FileDef {
	rows, err := db.Query(`SELECT id, path, filename, size, created_at, modified_at, checksum, hash
		FROM files
		WHERE hash != "";`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var files []models.FileDef
	for rows.Next() {
		var f models.FileDef
		err := rows.Scan(&f.ID, &f.Path, &f.Filename, &f.Size, &f.CreatedAt, &f.ModifiedAt, &f.Checksum, &f.Hash)
		if err != nil {
			panic(err)
		}
		files = append(files, f)
	}
	return files
}
