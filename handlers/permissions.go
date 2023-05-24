package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"calendarbot/utils"
)

func waitForForwardedMessage(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message, updates *tgbotapi.UpdatesChannel, botID int64) {
	utils.ReplyMessage(logger, bot, message, "Please forward here a message from user you want to grant your notes.")

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

		if update.Message.ForwardFrom == nil {
			utils.ReplyMessage(logger, bot, update.Message, "It is not forwarded message. No new users will be added")
			return
		}

		if update.Message.ForwardFrom.ID == message.From.ID {
			utils.ReplyMessage(logger, bot, update.Message, "It is your message. No new users will be added")
		} else if update.Message.ForwardFrom.ID == botID {
			utils.ReplyMessage(logger, bot, update.Message, "Ooh! So glad you like me :) However, I can't afford to look at your wonderful notes, sorry...")
		} else {
			inserted, err := utils.ChaeckAndInsertNewGrantedUser(logger, db, message.From.ID, update.Message.From)
			if err != nil {
				logger.Error("Error in waiting forwarded message", zap.Error(err))
			} else if !inserted {
				utils.ReplyMessage(logger, bot, update.Message, "The user is already granted to your notes!")
			} else {
				utils.ReplyMessage(logger, bot, update.Message, "Successfully granted!")
			}
		}
		return
	}
}

func sendNotGrantedMessageAndWait(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message, updates *tgbotapi.UpdatesChannel, botID int64) {
	utils.ReplyMessageWithOneTimeKeyboard(logger, bot, message, "You have not granted your notes to anyone. Would you like to add someone?", "Yes", "No")

	for update := range *updates {
		if update.Message == nil {
			continue
		}

		answer := update.Message.Text
		if answer == "No" {
			utils.ReplyMessage(logger, bot, update.Message, "Okay")
			return
		} else if answer == "Yes" {
			waitForForwardedMessage(logger, bot, db, update.Message, updates, botID)
			return
		}
	}
}

func sendListMessageAndWait(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message, updates *tgbotapi.UpdatesChannel, granted_users []string, botID int64) {
	reply := "You have already granted your notes to:"
	for _, user := range granted_users {
		reply = fmt.Sprintf("%s @%s", reply, user)
		if len(reply) > 900 {
			reply += " (truncated)..."
			break
		}
	}
	utils.ReplyMessage(logger, bot, message, reply)
	utils.ReplyMessageWithOneTimeKeyboard(logger, bot, message, "Would you like add more or erase all of them?", "Add", "Erase")

	for update := range *updates {
		if update.Message == nil {
			continue
		}

		answer := update.Message.Text
		if answer == "Add" {
			waitForForwardedMessage(logger, bot, db, update.Message, updates, botID)
			return
		} else if answer == "Erase" {
			err := utils.DeleteAllGrantedUsers(logger, db, update.Message.From.ID)
			if err != nil {
				msg := "An error occurred while deleting all granted users"
				logger.Error(msg, zap.Error(err))
				utils.ReplyMessage(logger, bot, update.Message, msg+". Please try later.")
			} else {
				utils.ReplyMessage(logger, bot, update.Message, "Successfully deleted!")
			}
			return
		}
	}
}

func HandlePermissionsCommand(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message, updates *tgbotapi.UpdatesChannel, botID int64) {
	granted_users, err := utils.GetUserPermissions(logger, db, message.From.ID)
	if err != nil {
		utils.ReplyMessage(logger, bot, message, "Error while getting permissions. Please try again later.")
		return
	}

	if len(granted_users) == 0 {
		sendNotGrantedMessageAndWait(logger, bot, db, message, updates, botID)
	} else {
		sendListMessageAndWait(logger, bot, db, message, updates, granted_users, botID)
	}
}
