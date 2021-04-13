package util

import (
	"encoding/json"
	"fmt"
	"github.com/DictumMortuum/servus/config"
	"github.com/jmoiron/sqlx"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"net/http"
	"time"
)

type telegramRecipient struct {
	Id        int64     `db:"id"`
	CrDate    time.Time `db:"cr_date"`
	Date      time.Time `db:"date"`
	Username  string    `db:"username" json:"username,omitempty"`
	UserId    int64     `db:"user_id" json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name,omitempty" default:""`
	Language  string    `db:"language_code" json:"language_code,omitempty" default:"en"`
	Intro     bool      `db:"intro"`
}

type telegramUpdates struct {
	OK bool `json:"ok"`
	Rs []struct {
		Id      int `json:"update_id"`
		Message struct {
			Id   int               `json:"message_id"`
			From telegramRecipient `json:"from"`
		} `json:"message"`
		Chat struct {
			Member struct {
				Status string `json:"status"`
			} `json:"new_chat_member"`
		} `json:"my_chat_member,omitempty"`
	} `json:"result"`
}

func (r telegramRecipient) Recipient() string {
	return fmt.Sprintf("%d", r.UserId)
}

func getUsers(db *sqlx.DB) ([]telegramRecipient, error) {
	rs := []telegramRecipient{}

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

func introduceUser(db *sqlx.DB, data telegramRecipient) error {
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

func createUser(db *sqlx.DB, data telegramRecipient) (int64, error) {
	sql := `
	insert into tboardgamepricesusers (
		cr_date,
		date,
		username,
		user_id,
		first_name,
		last_name,
		language_code,
		intro
	) values (
		NOW(),
		NOW(),
		:username,
		:user_id,
		:first_name,
		:last_name,
		:language_code,
		0
	) on duplicate key update
		date = NOW()
	`

	res, err := db.NamedExec(sql, &data)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func GetUpdates(db *sqlx.DB) (*telegramUpdates, error) {
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

	var rs telegramUpdates

	err = json.Unmarshal(body, &rs)
	if err != nil {
		return nil, err
	}

	for _, user := range rs.Rs {
		if user.Chat.Member.Status != "kicked" {
			_, err := createUser(db, user.Message.From)
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
			return err
		}
	}

	return nil
}
