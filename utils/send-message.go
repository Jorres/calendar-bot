package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func createMessage(message *tgbotapi.Message, reply string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, reply)
	msg.ReplyToMessageID = message.MessageID
	return msg
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

	sendOrLogError(logger, bot, msg)
}
