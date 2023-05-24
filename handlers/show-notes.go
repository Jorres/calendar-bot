package handlers

import (
	"calendarbot/utils"
	"database/sql"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func getUserNotes(logger *zap.Logger, db *sql.DB, userID int64) ([]string, error) {
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

func HandleShowNotesCommand(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message) {
	notes, err := getUserNotes(logger, db, message.From.ID)
	if err != nil {
		utils.ReplyMessage(logger, bot, message, "Error fetching notes. Please try again.")
		return
	}

	if len(notes) == 0 {
		utils.ReplyMessage(logger, bot, message, "You have no notes.")
		return
	}

	utils.ReplyMessage(logger, bot, message, "Here are your notes:\n\n"+strings.Join(notes, "\n"))
}
