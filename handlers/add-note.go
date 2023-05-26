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

const formatsHelp = "Date formats are:\n`/add 1917-11-07 09:20 ; my note`\n`/add 07 November 1917, 09:20 ; my note`\nIf you are not in UTC timezome, please specify your zone using format\n`/add 1917-11-07 09:20 +03:00 ; my note with time zone`\n`/add 07 November 1917, 09:20 -07:00 ; my note with zone`"

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

func parseDatetimeToUTC(datetimeStr string) (time.Time, error) {
	// Parse the datetime string in the provided formats
	formats := []string{
		"2006-01-02 15:04 -07:00",
		"2006-01-02 15:04",
		"02 January 2006, 15:04",
		"02 January 2006, 15:04 -07:00",
	}

	var parsedTime time.Time
	var err error

	// Attempt to parse the datetime string using each format
	for _, format := range formats {
		parsedTime, err = time.Parse(format, datetimeStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse datetime: %w", err)
	}

	// Convert the parsed time to UTC
	utcTime := parsedTime.UTC()

	return utcTime, nil
}

func HandleAddNoteCommand(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message) error {
	args := message.CommandArguments()
	datetime_string, note, err := parseAddNoteArguments(args)

	if err != nil {
		reply := "Please provide a date and a note in the format: /add <date> ; <note>"
		utils.ReplyMessage(logger, bot, message, utils.TransformMessage(reply)+formatsHelp)
		return errors.New(reply)
	}

	date, err := parseDatetimeToUTC(datetime_string)
	if err != nil {
		reply := "Invalid date format. Please use provided formats."
		utils.ReplyMessage(logger, bot, message, utils.TransformMessage(reply)+formatsHelp)
		return errors.New(reply)
	}

	err = queries.AddNote(logger, db, message, date.Format("2006-01-02 15:04:05"), note)
	if err != nil {
		reply := "Error adding note. Please try again."
		utils.ReplyMessage(logger, bot, message, utils.TransformMessage(reply))
		return errors.New(reply)
	}

	utils.ReplyMessage(logger, bot, message, fmt.Sprintf("Note successfully added:\n__*Date: %s*__\nNote: %s", utils.TransformMessage(datetime_string), utils.TransformMessage(note)))
	return nil
}

func handleEraseAllNotes(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message) {
	err := queries.EraseAllNotes(logger, db, message.From.ID)
	if err != nil {
		utils.ReplyMessage(logger, bot, message, utils.TransformMessage("Error erasing all notes. Please try again."))
	} else {
		utils.ReplyMessage(logger, bot, message, utils.TransformMessage("Succussfully deleted all notes."))
	}
}

func HandleNotesCommand(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message, updates *tgbotapi.UpdatesChannel) {
	utils.ReplyMessageWithOneTimeKeyboard(logger, bot, message, utils.TransformMessage("Would you like to add note or erase all of them?"), "Add", "Erase", "Go back")
	for update := range *updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "Add" {
			utils.ReplyMessageWithOneTimeKeyboard(logger, bot, update.Message, utils.TransformMessage("To add a note please use the format: /add <date> ; <note>")+"\n"+formatsHelp, "Add", "Erase", "Go back")
		} else if update.Message.IsCommand() {
			if update.Message.Command() == "add" {
				err := HandleAddNoteCommand(logger, bot, db, update.Message)
				if err == nil {
					logger.Error("User error on add logic", zap.Error(err))
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
