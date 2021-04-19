package models

import (
	"fmt"
	"time"
)

type TelegramRecipient struct {
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

func (obj TelegramRecipient) Recipient() string {
	return fmt.Sprintf("%d", obj.UserId)
}

func (obj TelegramRecipient) Insert() string {
	return `insert into tboardgamepricesusers (
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
	)`
}

func (obj TelegramRecipient) Exists() string {
	return `select id from tboardgamepricesusers where user_id = :user_id`
}

type TelegramUpdates struct {
	OK bool `json:"ok"`
	Rs []struct {
		Id      int `json:"update_id"`
		Message struct {
			Id   int               `json:"message_id"`
			From TelegramRecipient `json:"from"`
		} `json:"message"`
		Chat struct {
			Member struct {
				Status string `json:"status"`
			} `json:"new_chat_member"`
		} `json:"my_chat_member,omitempty"`
	} `json:"result"`
}
