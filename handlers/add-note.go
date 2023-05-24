package handlers

import (
	"calendarbot/utils"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func addNote(logger *zap.Logger, db *sql.DB, userID int64, day, note string) error {
	insertNoteQuery := `
	INSERT INTO notes (user_id, day, note)
	VALUES (?, ?, ?)
	`

	_, err := db.Exec(insertNoteQuery, userID, day, note)
	if err != nil {
		logger.Error("Error inserting note into database", zap.Error(err))
	}
	return err
}

func parseAddNoteArguments(args string) (string, string, error) {
	parts := strings.Split(args, ";")
	if len(parts) < 2 {
		return "", "", errors.New("not enough arguments")
	}

	day := strings.TrimSpace(parts[0])
	note := strings.TrimSpace(strings.Join(parts[1:], ";"))

	if day == "" || note == "" {
		return "", "", errors.New("invalid arguments")
	}

	return day, note, nil
}

func HandleAddNoteCommand(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message) error {
	args := message.CommandArguments()
	day, note, err := parseAddNoteArguments(args)

	if err != nil {
		reply := "Please provide a date and a note in the format: /add <date> ; <note>"
		utils.ReplyMessage(logger, bot, message, reply)
		return errors.New(reply)
	}

	date, err := time.Parse("02 January 2006", day)
	if err != nil {
		reply := "Invalid date format. Please use the format: \"dd MMMM yyyy\", e.g., \"27 April 2023\""
		utils.ReplyMessage(logger, bot, message, reply)
		return errors.New(reply)
	}

	err = addNote(logger, db, message.From.ID, date.Format("2006-01-02"), note)
	if err != nil {
		reply := "Error adding note. Please try again."
		utils.ReplyMessage(logger, bot, message, reply)
		return errors.New(reply)
	}

	utils.ReplyMessage(logger, bot, message, fmt.Sprintf("Note successfully added:\nDate: %s\nNote: %s", day, note))
	return nil
}
