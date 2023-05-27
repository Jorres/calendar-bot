package handlers

import (
	"calendarbot/queries"
	"calendarbot/utils"
	"database/sql"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func showUsersNotes(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message, userID int64, user_name string) {
	notes, err := queries.GetUserNotes(logger, db, userID)
	if err != nil {
		utils.ReplyMessage(logger, bot, message, "Error fetching notes\\. Please try again\\.")
		return
	}

	if len(notes) == 0 {
		utils.ReplyMessage(logger, bot, message, fmt.Sprintf("%s have no notes\\.", user_name))
		return
	}

	utils.ReplyMessage(logger, bot, message, fmt.Sprintf("Here are %s notes:\n\n", user_name)+strings.Join(notes, "\n"))
}

func showGrantedUsersNotes(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message, granted_user string) {
	granted_user = strings.TrimSpace(granted_user)
	if strings.HasPrefix(granted_user, "@") {
		granted_user = granted_user[1:]
	}
	if strings.Count(granted_user, " ") > 1 {
		utils.ReplyMessage(logger, bot, message, "The query contains more than one word\\.\nMust be only 1 login\\.\nFor example, `/show @three`")
	}

	granted, granted_user_id, err := queries.CheckUserGavePermission(logger, db, message.From.ID, granted_user)
	if err != nil {
		logger.Error("Got an error while checking granted user", zap.Error(err))
		utils.ReplyMessage(logger, bot, message, fmt.Sprintf("An error occurred while checking if @%s granted you an access\\. Please try later\\. error: %s", granted_user, err))
		return
	}
	if !granted {
		utils.ReplyMessage(logger, bot, message, utils.TransformMessage(fmt.Sprintf("@%s did not give you an access", granted_user)))
	} else {
		showUsersNotes(logger, bot, db, message, granted_user_id, utils.TransformMessage("@")+granted_user)
	}
}

func HandleShowNotesCommand(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message) {
	granted_user := message.CommandArguments()
	if granted_user == "" {
		showUsersNotes(logger, bot, db, message, message.From.ID, "your")
	} else {
		showGrantedUsersNotes(logger, bot, db, message, granted_user)
	}
}
