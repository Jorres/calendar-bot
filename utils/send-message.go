package utils

import (
	"html"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func createMessage(message *tgbotapi.Message, reply string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, reply)
	msg.ReplyToMessageID = message.MessageID
	return msg
}

func escapeUsernames(text string) string {
	// Create a regular expression that matches usernames.
	pattern := regexp.MustCompile(`^@[a-zA-Z0-9_]+$`)

	// Replace all matches with the escaped version.
	return pattern.ReplaceAllString(text, html.EscapeString(text))
}

func transformMessage(message tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	escapes := []string{"<", ">", ".", "-", "+", "(", ")", "{", "}", "[", "]", "!", "?", "@"}
	for _, escape := range escapes {
		message.Text = strings.ReplaceAll(message.Text, escape, "\\"+escape)
	}
	message.Text = escapeUsernames(message.Text)

	message.ParseMode = tgbotapi.ModeMarkdownV2
	return message
}

func sendOrLogError(logger *zap.Logger, bot *tgbotapi.BotAPI, message tgbotapi.MessageConfig) {
	// TODO In tests bot might be nil, will make proper mocking later
	// This is to not send any messages while testing
	if bot != nil {
		_, err := bot.Send(message)
		if err != nil {
			logger.Error("Error sending message", zap.Error(err))
		}
	}
}

func ReplyMessage(logger *zap.Logger, bot *tgbotapi.BotAPI, message *tgbotapi.Message, reply string) {
	sendOrLogError(logger, bot, transformMessage(createMessage(message, reply)))
}

func ReplyMessageOriginal(logger *zap.Logger, bot *tgbotapi.BotAPI, message *tgbotapi.Message, reply string) {
	sendOrLogError(logger, bot, createMessage(message, reply))
}

func ReplyMessageWithOneTimeKeyboard(logger *zap.Logger, bot *tgbotapi.BotAPI, message *tgbotapi.Message, reply string, inputs ...string) {
	msg := createMessage(message, reply)

	var input_list []string
	input_list = append(input_list, inputs...)

	var buttons []tgbotapi.KeyboardButton
	for _, input := range input_list {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(input))
	}
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(buttons)

	sendOrLogError(logger, bot, transformMessage(msg))
}
