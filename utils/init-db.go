package utils

import "database/sql"

func InitDB(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	shouldCreateTable := false

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='notes'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		shouldCreateTable = true
	}

	if shouldCreateTable {
		createTableQuery := `
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			day TEXT,
			note TEXT
		);
		`

		_, err = db.Exec(createTableQuery)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
