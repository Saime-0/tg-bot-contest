package created

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"

	"tgBotCompetition/model"
	chatUpdate "tgBotCompetition/usecase/chat/update"
	ticketCounting "tgBotCompetition/usecase/ticket/counting"
	userUpdate "tgBotCompetition/usecase/user/update"
)

type Params struct {
	DB *sqlx.DB

	Chat    model.Chat
	User    model.User
	Text    string
	TopicID int
}

type Out struct {
	CreatedTickets []model.Ticket
	Comp           model.Competition
}

func (p Params) Run() (Out, error) {
	if p.Text == "" {
		return Out{}, nil
	}

	if err := (&chatUpdate.Params{DB: p.DB, Chat: p.Chat}).Run(); err != nil {
		return Out{}, err
	}
	if err := (&userUpdate.Params{DB: p.DB, User: p.User}).Run(); err != nil {
		return Out{}, err
	}

	var comp model.Competition
	err := p.DB.Get(&comp, `
		select * from competitions 
		where chat_id=? 
			and topic_id=?
			and ended_at is null
	`, p.Chat.ID, p.TopicID)
	if errors.Is(err, sql.ErrNoRows) {
		return Out{}, nil
	} else if err != nil {
		return Out{}, err
	}

	text := strings.ToLower(p.Text)
	kw := strings.ToLower(comp.Keyword)

	if strings.Contains(text, "\""+kw+"\"") ||
		!strings.Contains(text, kw) {
		return Out{}, nil
	}
	var ticketCountingOut ticketCounting.Out
	if ticketCountingOut, err = (&ticketCounting.Params{
		DB:   p.DB,
		Chat: p.Chat,
		User: p.User,
		Comp: comp,
	}).Run(); err != nil {
		return Out{}, err
	}

	return Out{
		CreatedTickets: ticketCountingOut.CreatedTickets,
		Comp:           comp,
	}, nil
}
