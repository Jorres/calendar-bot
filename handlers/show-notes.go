package handlers

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func getUserNotes(db *sql.DB, userID int) ([]string, error) {
	selectUserNotesQuery := `
	SELECT day, note
	FROM notes
	WHERE user_id = ?
	ORDER BY day
	`

	rows, err := db.Query(selectUserNotesQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []string
	for rows.Next() {
		var day, note string
		err := rows.Scan(&day, &note)
		if err != nil {
			return nil, err
		}
		notes = append(notes, fmt.Sprintf("%s: %s", day, note))
	}

	return notes, nil
}

func HandleShowNotesCommand(bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message) {
	notes, err := getUserNotes(db, message.From.ID)
	if err != nil {
		reply := "Error fetching notes. Please try again."
		msg := tgbotapi.NewMessage(message.Chat.ID, reply)
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	if len(notes) == 0 {
		reply := "You have no notes."
		msg := tgbotapi.NewMessage(message.Chat.ID, reply)
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	reply := "Here are your notes:\n\n" + strings.Join(notes, "\n")
	msg := tgbotapi.NewMessage(message.Chat.ID, reply)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}
