package queries

import (
	"database/sql"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func AddNote(logger *zap.Logger, db *sql.DB, message *tgbotapi.Message, day, note string) error {
	user := message.From
	InsertUserWithChat(logger, db, user, message.Chat.ID)

	insertNoteQuery := `
	INSERT INTO notes (user_id, event_date, note, reminder_sent)
	VALUES (?, ?, ?, 0)
	`

	_, err := db.Exec(insertNoteQuery, user.ID, day, note)
	fmt.Printf("INSERTED user.ID=%d, day='%s', note='%s'\n", user.ID, day, note)
	if err != nil {
		logger.Error("Error inserting note into database", zap.Error(err))
	}
	return err
}

func EraseAllNotes(logger *zap.Logger, db *sql.DB, userID int64) error {
	kDeleteNotes := `
	DELETE FROM notes
	WHERE user_id = ?
	`

	_, err := db.Exec(kDeleteNotes, userID)
	if err != nil {
		logger.Error("Error deleting all notes from database", zap.Error(err))
	}
	return err
}

func GetUserNotes(logger *zap.Logger, db *sql.DB, userID int64) ([]string, error) {
	selectUserNotesQuery := `
	SELECT day, note
	FROM notes
	WHERE user_id = ?
	ORDER BY day
	`

	rows, err := db.Query(selectUserNotesQuery, userID)
	if err != nil {
		logger.Error("Error querying user notes", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var notes []string
	for rows.Next() {
		var day, note string
		err := rows.Scan(&day, &note)
		if err != nil {
			logger.Error("Error scanning row", zap.Error(err))
			return nil, err
		}
		notes = append(notes, fmt.Sprintf("%s: %s", day, note))
	}

	return notes, nil
}
