package queries

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// for debug
func showDB(db *sql.DB) {
	fmt.Println("======= users =======")
	rows, err := db.Query(`
	SELECT id, name, chat_id
	FROM users
	`)
	if err != nil {
		fmt.Println("Error querying events", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, chat_id int64
		var name string
		rows.Scan(&id, &name, &chat_id)
		if err != nil {
			fmt.Println("Error scanning event row", err)
			continue
		}
		fmt.Printf("id=%d, name='%s', chat_id=%d\n", id, name, chat_id)
	}

	fmt.Println("======= notes with id =======")
	rows, err = db.Query(`
	SELECT id, user_id, event_date, note, reminder_sent
	FROM notes
	`)
	if err != nil {
		fmt.Println("Error querying events", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id sql.NullInt64
		var user_id int64
		var event_date, note string
		var reminder_sent int
		rows.Scan(&id, &user_id, &event_date, &note, &reminder_sent)
		if err != nil {
			fmt.Println("Error scanning event row", err)
			continue
		}
		fmt.Printf("id=%d, user_id=%d, event_date='%s'\nnote='%s'\nreminder_sent=%d\n", id.Int64, user_id, event_date, note, reminder_sent)
	}
	fmt.Println("====================")
}

func NotifySender(logger *zap.Logger, bot *tgbotapi.BotAPI, db *sql.DB) {
	for {
		time.Sleep(time.Minute)

		logger.Info("Running NotifySender..")
		// fmt.Println("Running NotifySender..")
		// showDB(db)

		rows, err := db.Query(`
		SELECT n.id, n.note, u.chat_id, u.name, n.event_date, u.id
		FROM notes n
		INNER JOIN users u ON u.id = n.user_id
		WHERE event_date BETWEEN datetime('now') AND datetime('now', '+5 minutes')
		  AND reminder_sent = 0`)
		if err != nil {
			logger.Error("Error querying events", zap.Error(err))
			continue
		}
		defer rows.Close()

		var event_ids []int
		messages := []struct {
			chatID     int64
			note       string
			login      string
			event_date string
			userID     int64
		}{}
		for rows.Next() {
			var eventID int
			var note, login, event_date string
			var chatID, userID int64
			err := rows.Scan(&eventID, &note, &chatID, &login, &event_date, &userID)
			if err != nil {
				logger.Error("Error scanning event row", zap.Error(err))
				continue
			}

			messages = append(messages, struct {
				chatID     int64
				note       string
				login      string
				event_date string
				userID     int64
			}{chatID, note, login, event_date, userID})
			event_ids = append(event_ids, eventID)
		}

		args := make([]interface{}, len(event_ids))
		for i, id := range event_ids {
			args[i] = id
		}
		_, err = db.Exec("UPDATE notes SET reminder_sent = 1 WHERE id IN ("+placeholders(len(event_ids))+")", args...)
		if err != nil {
			logger.Error("Error updating events: ", zap.Ints("event_ids", event_ids), zap.Error(err))
		}

		for _, message := range messages {
			chats, err := GetAllGrantedUsersChats(logger, db, message.userID)
			message_string := " event is starting in 5 minutes!\n\n" + message.note + "\nPlanned at:" + message.event_date
			bot.Send(tgbotapi.NewMessage(message.chatID, "Your"+message_string))
			for _, chat := range chats {
				msg := tgbotapi.NewMessage(chat.chat, "@"+chat.login+message_string)
				_, err = bot.Send(msg)
				if err != nil {
					logger.Warn("Could not send a notify for "+message.login, zap.Error(err))
				}
			}
		}
	}
}

func placeholders(n int) string {
	placeholders := make([]string, n)
	for i := 0; i < n; i++ {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}
