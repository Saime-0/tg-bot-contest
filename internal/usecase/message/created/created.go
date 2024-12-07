package created

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/Saime-0/tg-bot-contest/internal/autotx"
	"github.com/Saime-0/tg-bot-contest/internal/model"
	chatUpdate "github.com/Saime-0/tg-bot-contest/internal/usecase/chat/update"
	ticketCounting "github.com/Saime-0/tg-bot-contest/internal/usecase/ticket/counting"
	userUpdate "github.com/Saime-0/tg-bot-contest/internal/usecase/user/update"
)

type Params struct {
	DB *sqlx.DB

	Chat    model.Chat
	User    model.User
	Text    string
	TopicID int
}

type Out struct {
	CreatedTickets        []model.Ticket
	Comp                  model.Contest
	CalculationWasStarted bool
}

func (p Params) Run() (Out, error) {
	if p.Text == "" {
		return Out{}, nil
	}

	if err := chatUpdate.Run(p.DB, p.Chat); err != nil {
		return Out{}, err
	}
	if err := userUpdate.Run(p.DB, p.User); err != nil {
		return Out{}, err
	}

	tx, err := p.DB.BeginTxx(context.Background(), nil)
	defer func() { autotx.Commit(tx, err, recover()) }()

	var comp model.Contest
	err = tx.Get(&comp, `
		select * from contests 
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
		TX:   tx,
		Chat: p.Chat,
		User: p.User,
		Comp: comp,
	}).Run(); err != nil {
		return Out{}, err
	}

	return Out{
		CreatedTickets:        ticketCountingOut.CreatedTickets,
		Comp:                  comp,
		CalculationWasStarted: true,
	}, nil
}
