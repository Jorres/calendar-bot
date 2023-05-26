package tests

import (
	"calendarbot/handlers"
	"calendarbot/utils"

	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

var testHandleAddNoteCommandSuccessfulParams = []struct {
	date string
	note string
}{
	{"22 January 2021, 00:00", "Test note"},
	{"19 February 0001, 00:00", "Тест на кириллице"},
	{"01 March 9999, 00:00", "اختبار على اللغة العربية"},
	{"30 April 0721, 00:00", "中文測試"},
	{"07 May 1919, 00:00", "Иван Иванов Москва Москва улица Ленина метро Ⓜ проспект Вернадского 😁"},
	{"15 June 7070, 00:00", "ХВОСТ УДАЧИ – БЛАГОТВОРИТЕЛЬНЫЙ ПРОЕКТ, КОТОРЫЙ СОЗДАН ДЛЯ ПОМОЩИ ДОМАШНИМ ПИТОМЦАМ И ИХ ХОЗЯЕВАМ. Искать по фото      Объявления     В добрые руки     СМИ о нас     О проекте      Создать объявление     Войти         Объявления     В добрые руки     СМИ о нас     О проекте  Искать по фото Создать объявление Умный поиск домашних животных Умный поиск домашних животных  Найди питомца по фотографии Мой питомец потерялся Мой питомец потерялся  Создать объявление о пропаже питомца, чтобы все увидели Создать объявление Найден чужой питомец Найден чужой питомец  Вы нашли чужого питомца и хотите разместить объявление Создать объявление О нас говорят Перейти в раздел О нас О нас  Наша умная доска объявлений поможет найти пропавших собак и кошек. Мы используем искусственный интеллект для быстрого и удобного поиска питомцев с помощью простой фотографии. Если в базе есть объявления о пропаже и находке с фото одного и того же животного, публикации будут сопоставлены автоматически.  Мы хотим, чтобы люди более ответственно относились к своим питомцам, а общество стало более гуманным. И готовы в этом помочь. Заботиться – легко: не проходите мимо животных, которым нужна помощь. Присоединяйтесь к проекту. Вместе мы сделаем мир лучше! Подпишитесь на новости  И будьте в курсе всех событий из мира технологий  Отправляя форму, вы соглашаетесь с политикой конфиденциальности subscribe Как это работает 011Поиск по фотографии Поиск по фотографии  Используйте поиск по фото! Возможно, питомец которого вы ищете, уже есть в базе объявлений. 022Cоздайте объявление Cоздайте объявление  Создайте объявление о пропавшем или найденном питомце. Это просто! 033ИИ подберет питомца ИИ подберет питомца  Искусственный интеллект подберет объявления с похожими питомцами. 04Питомец вернулся! Питомец вернулся!  Питомец вернулся домой. Мы работаем для того, чтобы сделать мир лучше! leftright Попробуйте прямо сейчас Искать по фото Создать объявление Отдать или взять в добрые руки  Наш сервис позволяет найти или отдать питомца в добрые руки. Отдать питомца в добрые руки Отдать питомца в добрые руки  Вы решили отдать своего питомца другому хозяину? Создать объявление Взять питомца в добрые руки Взять питомца в добрые руки  Вы наконец решили взять питомца к себе домой? Взять питомца Вопрос-ответ Что делать, если пропала собака? Что делать, если пропала кошка? Как пользоваться базой объявлений? Какой должна быть фотография, чтобы вы смогли найти собаку или кошку по фото? Что делать, если поиск по фотографии ничего не дал? Чем проект «Хвост удачи» отличается от других досок объявлений по поиску домашних животных? Каких домашних питомцев я могу поискать по фото? Что еще я могу сделать, чтобы помочь животным? Подпишитесь на новости  И будьте в курсе всех событий из мира технологий  Отправляя форму, вы соглашаетесь с политикой конфиденциальности subscribe  Все права защищены и охраняются действующим законодательством РФ.  Администрация сайта не несет ответственности за содержание размещенных объявлений ГОЗНАК РКФ РКФ  Информационный партнер и консультант проекта      Объявления     СМИ о нас     В добрые руки     О проекте     Создать объявление     Обратная связь      https://vk.com/club206650352  © 2023 АО «ГОЗНАК». Политика конфиденциальности  На сайте используются файлы cookie  Оставаясь на сайте, вы выражаете свое согласие на обработку персональных данных в соответствии с политикой АО «ГОЗНАК» и соглашаетесь с политикой обработки файлов cookie"},
	{"07 July 0090, 00:00", ";;;;А как тебе;такое;Илон;Маск?"},
	{"22 August 2021, 00:00", "ꯆꯥꯏꯅꯤꯖꯗꯥ ꯇꯦꯁ꯭ꯠ ꯇꯧꯕꯥ꯫"},
	{"22 September 1010, 00:00", "ხაჭაპური იყიდე იაფად ძალიან გემრიელი"},
	{"22 October 2021, 00:00", "ꯆꯥꯏꯅꯤꯖꯗꯥ ꯇꯦꯁ꯭ꯠ ꯇꯧꯕꯥ꯫"},
	{"22 November 2021, 00:00", "အင်း ဒါဆို ဘယ်လိုလဲ။"},
	{"15 December 0021, 00:00", "من زندگی می کنم - بنابراین وجود دارم"},
	{"9999-12-29 23:59", "من زندگی می کنم - بنابراین وجود دارم"},
}

func testHandleAddNoteCommandSuccessful(date string, note string, t *testing.T) {
	db, err := utils.InitDB("test.db")
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	message := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
		Text: fmt.Sprintf("/add %s ; %s", date, note),
		Entities: []tgbotapi.MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: 4,
			},
		},
		From: &tgbotapi.User{
			ID: 12345,
		},
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	err = handlers.HandleAddNoteCommand(logger, nil, db, message)
	if err != nil {
		t.Fatalf("HandleAddNoteCommand (/add) returned an error: %v", err)
	}

	// Check that the note was added to the mock database
	rows, err := db.Query("SELECT note FROM notes WHERE user_id = ?", message.From.ID)
	if err != nil {
		t.Fatalf("Failed to query mock database: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Errorf("HandleAddNoteCommand did not add note to database")
	}

	for rows.Next() {
		var res_note string
		err = rows.Scan(&res_note)
		if err != nil {
			t.Errorf("Error while reading SELECT content: %v", err)
		}
		if res_note != note {
			t.Errorf("Expected %s as note, but got %s", note, res_note)
		}
	}

	_, err = db.Exec("DROP TABLE users")
	if err != nil {
		t.Fatalf("Failed to drop table from mock database: %v", err)
	}
	_, err = db.Exec("DROP TABLE notes")
	if err != nil {
		t.Fatalf("Failed to drop table from mock database: %v", err)
	}
	_, err = db.Exec("DROP TABLE permissions")
	if err != nil {
		t.Fatalf("Failed to drop table from mock database: %v", err)
	}
}

func TestHandleAddNoteCommandSuccessful(t *testing.T) {
	for _, tt := range testHandleAddNoteCommandSuccessfulParams {
		testHandleAddNoteCommandSuccessful(tt.date, tt.note, t)
	}
}

var testHandleAddNoteCommandDateParseFailParams = []struct {
	date string
}{
	{"02/17/2009"},
	{"17/02/2009"},
	{"2009/02/17"},
	{"February 17, 2009"},
	{"2/17/2009"},
	{"17/2/2009"},
	{"2009/2/17"},
	{" 2/17/2009"},
	{"17/ 2/2009"},
	{"2009/ 2/17"},
	{"02172009"},
	{"Feb172009"},
	{"17 February, 2009"},
	{"17 Feb 2009"},
	{"99 April 9999"},
	{"-1 April 9999"},
	{"00 April 9999"},
	{"01 April 0"},
	{"01 Апреля 2024"},
	{"二〇二二年一〇月二二日"},
}

func testHandleAddNoteCommandDateParseFail(date string, t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	note := "Test note"

	// Create a mock message
	message := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
		Text: fmt.Sprintf("/add %s ; %s", date, note),
		Entities: []tgbotapi.MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: 4,
			},
		},
		From: &tgbotapi.User{
			ID: 12345,
		},
	}

	err = handlers.HandleAddNoteCommand(logger, nil, nil, message)
	if err == nil {
		t.Fatalf("HandleAddNoteCommand (/add) does not return an error: %v", err)
	}
}

func TestHandleAddNoteCommandDateParseFail(t *testing.T) {
	for _, tt := range testHandleAddNoteCommandDateParseFailParams {
		testHandleAddNoteCommandDateParseFail(tt.date, t)
	}
}

// Uncomment for test locally
// func TestMain(t *testing.T) {
// 	// Create a context with a 1 second timeout.
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

// 	defer run.Run()

// 	// Wait for the goroutine to finish or timeout.
// 	select {
// 	case <-ctx.Done():
// 		// The goroutine has not finished, so it must have been stopped by the timeout.
// 		t.Errorf("goroutine did not finish within 1 second")
// 	default:
// 		// The goroutine finished before the timeout expired.
// 	}

// 	// Cancel the context to prevent any goroutines that are still running from being leaked.
// 	cancel()
// }
