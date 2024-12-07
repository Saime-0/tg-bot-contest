package create

import (
	"github.com/jmoiron/sqlx"

	"github.com/Saime-0/tg-bot-contest/internal/l10n"
	"github.com/Saime-0/tg-bot-contest/internal/model"
	"github.com/Saime-0/tg-bot-contest/internal/ue"
)

type Params struct {
	TX *sqlx.Tx

	Multiplicity int
	Keyword      string
	ChatID       int
	TopicID      int
	CreatorID    int
}

const (
	DefaultMultiplicity = 10
	DefaultKeyword      = l10n.DefaultKeyword
)

func (p *Params) Run() error {
	var exists bool
	if err := p.TX.Get(&exists, `
		select exists(
		    select 1 from contests 
			where chat_id=? 
			  and ended_at is null
		)
  	`, p.ChatID); err != nil {
		return err
	}
	if exists {
		return ue.New(l10n.ContestCreatePreviousNotOverYet)
	}

	if p.Multiplicity <= 0 {
		p.Multiplicity = DefaultMultiplicity
	}
	if p.Keyword == "" {
		p.Keyword = DefaultKeyword
	}

	_, err := p.TX.NamedExec(`
		insert into contests (creator_id,chat_id,topic_id,keyword,multiplicity)
		values (:creator_id,:chat_id,:topic_id,:keyword,:multiplicity)`,
		model.Contest{
			CreatorID:    p.CreatorID,
			ChatID:       p.ChatID,
			TopicID:      p.TopicID,
			Keyword:      p.Keyword,
			Multiplicity: p.Multiplicity,
		})

	return err
}
