package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"calendarbot/utils"
)

func waitForForwardedMessage(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message, updates *tgbotapi.UpdatesChannel) {
	reply := "Please forward here a message from user you want grant to your notes."
	bot.Send(tgbotapi.NewMessage(message.Chat.ID, reply))

	for update := range *updates {
		if update.Message == nil {
			continue
		}

		json_msg, err := json.Marshal(update.Message)
		if err != nil {
			logger.Error("Cannot convert message into json", zap.Error(err))
			return
		}
		logger.Info("Forwarded message",
			zap.String("message_json", string(json_msg)),
		)

		if update.Message.From.ID == message.From.ID {
			reply := "It is your message. No new users will be added"
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, reply))
			return
		} else {
			inserted, err := utils.ChaeckAndInsertNewGrantedUser(logger, db, message.From.ID, update.Message.From)
			if err != nil {
				logger.Error("Error in waiting forwarded message", zap.Error(err))
			} else if !inserted {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "The user is already granted to your notes!"))
			} else {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Successfully granted!"))
			}
			return
		}
	}
}

func sendNotGrantedMessageAndWait(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message, updates *tgbotapi.UpdatesChannel) {
	first_message := "You have not granted your notes to anyone. Would you like to add someone?"
	msg := tgbotapi.NewMessage(message.Chat.ID, first_message)

	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Yes"),
			tgbotapi.NewKeyboardButton("No"),
		),
	)
	bot.Send(msg)

	for update := range *updates {
		if update.Message == nil {
			continue
		}

		answer := update.Message.Text
		if answer == "No" {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Okay"))
			return
		} else if answer == "Yes" {
			waitForForwardedMessage(logger, bot, db, update.Message, updates)
			return
		}
	}
}

func sendListMessageAndWait(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message, updates *tgbotapi.UpdatesChannel, granted_users []string) {
	reply := "You have already granted your notes to:"
	for _, user := range granted_users {
		reply = fmt.Sprintf("%s @%s", reply, user)
		if len(reply) > 900 {
			reply += " (truncated)..."
			break
		}
	}
	bot.Send(tgbotapi.NewMessage(message.Chat.ID, reply))

	question_reply := "Would you like add more or erase all of them?"
	msg := tgbotapi.NewMessage(message.Chat.ID, question_reply)

	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Add"),
			tgbotapi.NewKeyboardButton("Erase"),
		),
	)
	bot.Send(msg)

	for update := range *updates {
		if update.Message == nil {
			continue
		}

		answer := update.Message.Text
		if answer == "Add" {
			waitForForwardedMessage(logger, bot, db, update.Message, updates)
			return
		} else if answer == "Erase" {
			err := utils.DeleteAllGrantedUsers(logger, db, update.Message.From.ID)
			if err != nil {
				msg := "An error occurred while deleting all granted users"
				logger.Error(msg, zap.Error(err))
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg+". Please try later."))
			} else {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Successfully deleted!"))
			}
			return
		}
	}
}

func HandlePermissionsCommand(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message, updates *tgbotapi.UpdatesChannel) {
	granted_users, err := utils.GetUserPermissions(logger, db, message.From.ID)
	if err != nil {
		reply := "Error while getting permissions. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, reply)
		msg.ReplyToMessageID = message.MessageID
		_, err := bot.Send(msg)
		if err != nil {
			logger.Error("Error sending message", zap.Error(err))
		}
		return
	}

	if len(granted_users) == 0 {
		sendNotGrantedMessageAndWait(logger, bot, db, message, updates)
	} else {
		sendListMessageAndWait(logger, bot, db, message, updates, granted_users)
	}
}
