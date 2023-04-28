package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	_ "github.com/mattn/go-sqlite3"

	"calendarbot/handlers"
	"calendarbot/utils"
)

func main() {
	tokenFile := "token.txt"
	botToken, err := readTokenFromFile(tokenFile)
	if err != nil {
		log.Panicf("Error reading token from file %s: %v", tokenFile, err)
	}

	db, err := utils.InitDB("notes.db")
	if err != nil {
		log.Panicf("Error initializing database: %v", err)
	}
	defer db.Close()

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "add":
				handlers.HandleAddNoteCommand(bot, db, update.Message)
			case "show":
				handlers.HandleShowNotesCommand(bot, db, update.Message)
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
