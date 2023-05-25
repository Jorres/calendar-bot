package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	_ "github.com/mattn/go-sqlite3"

	"calendarbot/handlers"
	"calendarbot/utils"

	"go.uber.org/zap"
)

func main() {
	config := zap.NewProductionConfig()

	config.OutputPaths = []string{
		"/app/logs/calendarbot.log",
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	tokenFile := "token.txt"
	botToken, err := readTokenFromFile(tokenFile)
	if err != nil {
		logger.Panic("Error reading token from file", zap.String("file", tokenFile), zap.Error(err))
	}

	botIDStr, err := readTokenFromFile("bot_id.txt")
	if err != nil {
		logger.Panic("Could not read bot_id.txt", zap.Error(err))
		return
	}
	botID, err := strconv.ParseInt(botIDStr, 10, 64)
	if err != nil {
		logger.Panic("Could not parse bot_id", zap.Error(err))
		return
	}

	db, err := utils.InitDB("notes.db")
	if err != nil {
		logger.Panic("Error initializing database", zap.Error(err))
	}
	defer db.Close()

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logger.Panic("Error creating new Bot API", zap.Error(err))
	}

	bot.Debug = true
	logger.Info("Authorized on account", zap.String("account", bot.Self.UserName))

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	helpMessage := "Available commands:\n/note -- edit your notes.\n/show -- show your notes.\n/show <user> -- show the notes of users who have given you access.\n/permissions - setting access permissions for other users."

	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		logger.Info("Received message",
			zap.String("user", update.Message.From.UserName),
			zap.String("text", update.Message.Text),
		)

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start", "help":
				utils.ReplyMessage(logger, bot, update.Message, helpMessage)
			case "notes":
				handlers.HandleNotesCommand(logger, bot, db, update.Message, &updates)
			case "show":
				handlers.HandleShowNotesCommand(logger, bot, db, update.Message)
			case "permissions":
				handlers.HandlePermissionsCommand(logger, bot, db, update.Message, &updates, botID)
			default:
				utils.ReplyMessage(logger, bot, update.Message, "Unknown command.")
			}
		}
	}
}

func readTokenFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	token, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(token), nil
}
