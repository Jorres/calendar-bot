package handlers

import (
	"calendarbot/queries"
	"calendarbot/utils"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

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

	err = queries.AddNote(logger, db, message.From, date.Format("2006-01-02"), note)
	if err != nil {
		reply := "Error adding note. Please try again."
		utils.ReplyMessage(logger, bot, message, reply)
		return errors.New(reply)
	}

	utils.ReplyMessage(logger, bot, message, fmt.Sprintf("Note successfully added:\nDate: %s\nNote: %s", day, note))
	return nil
}

func handleEraseAllNotes(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message) {
	err := queries.EraseAllNotes(logger, db, message.From.ID)
	if err != nil {
		utils.ReplyMessage(logger, bot, message, "Error erasing all notes. Please try again.")
	} else {
		utils.ReplyMessage(logger, bot, message, "Succussfully deleted all notes.")
	}
}

func HandleNotesCommand(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message, updates *tgbotapi.UpdatesChannel) {
	utils.ReplyMessageWithOneTimeKeyboard(logger, bot, message, "Would you like to add note or erase all of them?", "Add", "Erase", "Go back")
	for update := range *updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "Add" {
			utils.ReplyMessageWithOneTimeKeyboard(logger, bot, update.Message, "To add a note please use the format: /add \"dd MMMM yyyy\" ; <note>\ne.g., /add 07 November 1917 ; my note", "Add", "Erase", "Go back")
		} else if update.Message.IsCommand() {
			if update.Message.Command() == "add" {
				err := HandleAddNoteCommand(logger, bot, db, update.Message)
				if err == nil {
					return
				}
			}
		} else if update.Message.Text == "Erase" {
			handleEraseAllNotes(logger, bot, db, update.Message)
			return
		} else if update.Message.Text == "Go back" {
			utils.ReplyMessage(logger, bot, update.Message, "Okay, let's go back")
			return
		}
	}
}
