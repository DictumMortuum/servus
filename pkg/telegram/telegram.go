package telegram

import (
	"encoding/json"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/config"
	DB "github.com/DictumMortuum/servus/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"net/http"
)

func getUsers(db *sqlx.DB) ([]models.TelegramRecipient, error) {
	rs := []models.TelegramRecipient{}

	sql := `
	select
		*
	from
		tboardgamepricesusers`

	err := db.Select(&rs, sql)
	if err != nil {
		return rs, err
	}

	return rs, nil
}

func introduceUser(db *sqlx.DB, data models.TelegramRecipient) error {
	sql := `
	update
		tboardgamepricesusers
	set
		intro = 1
	where
		id = :id`

	_, err := db.NamedExec(sql, &data)
	if err != nil {
		return err
	}

	return nil
}

func GetUpdates(db *sqlx.DB) (*models.TelegramUpdates, error) {
	req, err := http.NewRequest("GET", "https://api.telegram.org/bot"+config.App.Telegram.Token+"/getUpdates", nil)
	if err != nil {
		return nil, err
	}

	conn := &http.Client{}

	resp, err := conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rs models.TelegramUpdates

	err = json.Unmarshal(body, &rs)
	if err != nil {
		return nil, err
	}

	for _, user := range rs.Rs {
		if user.Chat.Member.Status != "kicked" {
			_, err := DB.InsertIfNotExists(db, user.Message.From)
			if err != nil {
				return nil, err
			}
		}
	}

	SendIntros(db)

	return &rs, nil
}

func SendIntros(db *sqlx.DB) error {
	if !config.App.Telegram.Enabled {
		return nil
	}

	settings := tb.Settings{
		Token: config.App.Telegram.Token,
	}

	bot, err := tb.NewBot(settings)
	if err != nil {
		return err
	}

	users, err := getUsers(db)
	if err != nil {
		return err
	}

	for _, user := range users {
		if !user.Intro {
			fmt.Printf("Hello %s! Thanks for subscribing! Once new offers are available, you'll receive them here.\n", user.FirstName)
			msg := fmt.Sprintf("Hello %s! Thanks for subscribing! Once new offers are available, you'll receive them here.\n", user.FirstName)
			_, err := bot.Send(user, msg)
			if err != nil {
				return err
			}

			err = introduceUser(db, user)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func TelegramMessage(db *sqlx.DB, message string) error {
	if !config.App.Telegram.Enabled {
		fmt.Println(message)
		return nil
	}

	settings := tb.Settings{
		Token: config.App.Telegram.Token,
	}

	bot, err := tb.NewBot(settings)
	if err != nil {
		return err
	}

	users, err := getUsers(db)
	if err != nil {
		return err
	}

	for _, user := range users {
		_, err := bot.Send(user, message)
		if err != nil {
			log.Println(user, err)
		}
	}

	return nil
}
