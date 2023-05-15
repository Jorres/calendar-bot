package handlers

import (
	"database/sql"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

func getUserNotes(logger *zap.Logger, db *sql.DB, userID int) ([]string, error) {
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
		reply := "Error fetching notes. Please try again."
		msg := tgbotapi.NewMessage(message.Chat.ID, reply)
		msg.ReplyToMessageID = message.MessageID
		_, err := bot.Send(msg)
		if err != nil {
			logger.Error("Error sending message", zap.Error(err))
		}
		return
	}

	if len(notes) == 0 {
		reply := "You have no notes."
		msg := tgbotapi.NewMessage(message.Chat.ID, reply)
		msg.ReplyToMessageID = message.MessageID
		_, err := bot.Send(msg)
		if err != nil {
			logger.Error("Error sending message", zap.Error(err))
		}
		return
	}

	reply := "Here are your notes:\n\n" + strings.Join(notes, "\n")
	msg := tgbotapi.NewMessage(message.Chat.ID, reply)
	msg.ReplyToMessageID = message.MessageID
	_, err = bot.Send(msg)
	if err != nil {
		logger.Error("Error sending message", zap.Error(err))
	}
}
