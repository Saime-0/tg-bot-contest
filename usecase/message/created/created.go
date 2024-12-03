package created

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"

	"tgBotCompetition/model"
	chatUpdate "tgBotCompetition/usecase/chat/update"
	ticketRecalc "tgBotCompetition/usecase/ticket/recalc"
	userUpdate "tgBotCompetition/usecase/user/update"
)

type Params struct {
	DB *sqlx.DB

	Chat    model.Chat
	User    model.User
	Text    string
	TopicID int
}

func (p Params) Run() error {
	if p.Text == "" {
		return nil
	}

	if err := (&chatUpdate.Params{DB: p.DB, Chat: p.Chat}).Run(); err != nil {
		return err
	}
	if err := (&userUpdate.Params{DB: p.DB, User: p.User}).Run(); err != nil {
		return err
	}

	var comp model.Competition
	err := p.DB.Get(&comp, `
		select * from competitions 
		where chat_id=? 
			and topic_id=?
			and ended_at is null
	`, p.Chat.ID, p.TopicID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if err != nil {
		return err
	}

	text := strings.ToLower(p.Text)
	kw := strings.ToLower(comp.Keyword)

	if strings.Contains(text, "\""+kw+"\"") ||
		!strings.Contains(text, kw) {
		return nil
	}

	if err = (&ticketRecalc.Params{
		DB:   p.DB,
		Chat: p.Chat,
		User: p.User,
		Comp: comp,
	}).Run(); err != nil {
		return err
	}

	return nil
}
