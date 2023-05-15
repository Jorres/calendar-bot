package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	_ "github.com/mattn/go-sqlite3"

	"calendarbot/handlers"
	"calendarbot/utils"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	tokenFile := "token.txt"
	botToken, err := readTokenFromFile(tokenFile)
	if err != nil {
		logger.Panic("Error reading token from file", zap.String("file", tokenFile), zap.Error(err))
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

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		logger.Panic("Error getting updates channel", zap.Error(err))
	}

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
			case "add":
				handlers.HandleAddNoteCommand(logger, bot, db, update.Message)
			case "show":
				handlers.HandleShowNotesCommand(logger, bot, db, update.Message)
			default:
				fmt.Println(update.Message.Command())
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command.")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
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
