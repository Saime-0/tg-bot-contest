package recalc

import (
	"github.com/jmoiron/sqlx"

	"tgBotCompetition/model"
)

type Params struct {
	DB *sqlx.DB

	Chat model.Chat
	User model.User
	Comp model.Competition
}

func (p Params) Run() error {
	return nil
}
