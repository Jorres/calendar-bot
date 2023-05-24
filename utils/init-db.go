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
			user_id BIGINT,
			day TEXT,
			note TEXT
		);

		CREATE TABLE IF NOT EXISTS permissions (
			user_id BIGINT,
			granted_user_id BIGINT,
			granted_user_login TEXT,

			CONSTRAINT unique__user_id__granted_user_id UNIQUE (user_id, granted_user_id)
		);
		`

		_, err = db.Exec(createTableQuery)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
