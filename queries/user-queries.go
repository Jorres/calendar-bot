package queries

import (
	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func InsertUser(logger *zap.Logger, db *sql.DB, user *tgbotapi.User) {
	kInsertUser := `
	INSERT INTO users (id, name)
	VALUES (?, ?)
	`

	_, err := db.Exec(kInsertUser, user.ID, user.UserName)
	if err != nil {
		logger.Warn("Error inserting user into database. Maybe there is already a user", zap.Error(err))
	}
}
